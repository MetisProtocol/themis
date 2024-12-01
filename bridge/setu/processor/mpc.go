package processor

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/metis-seq/themis/bridge/setu/rpc"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/helper"
	"github.com/metis-seq/themis/helper/tss"
	mpcTypes "github.com/metis-seq/themis/mpc/types"
	"github.com/metis-seq/themis/types"
)

const (
	localMpcCreateKey = "local-mpc-create"
)

// MpcProcessor - process mpc related events
type MpcProcessor struct {
	BaseProcessor

	// header listener subscription
	cancelMpcService context.CancelFunc
}

// Start starts new block subscription
func (mp *MpcProcessor) Start() error {
	mp.Logger.Info("Starting")

	// create cancellable context
	mpcCtx, cancelMpcService := context.WithCancel(context.Background())
	mp.cancelMpcService = cancelMpcService

	// start polling for mpc-create
	config := helper.GetConfig()
	mp.Logger.Info("Start polling for mpc", "pollInterval", config.MpcPollInterval, "blobUpgradeHeight", config.BlobUpgradeHeight)
	go mp.startPolling(mpcCtx, config.MpcPollInterval, config.BlobUpgradeHeight)

	return nil
}

// RegisterTasks - nil
func (mp *MpcProcessor) RegisterTasks() {
}

// startPolling - polls themis and checks if new mpc needs to be proposed
func (mp *MpcProcessor) startPolling(ctx context.Context, interval time.Duration, blobUpgradeHeight uint64) {
	ticker := time.NewTicker(interval)
	dbTicker := time.NewTicker(10 * time.Millisecond)
	// stop ticker when everything done
	defer ticker.Stop()
	defer dbTicker.Stop()

	for {
		select {
		case <-dbTicker.C:
			mp.KeySign()
		case <-ticker.C:
			if !util.MpcCommonGenerated {
				mp.checkAndPropose(types.CommonMpcType)
			}

			if !util.MpcStateCommitGenerated {
				mp.checkAndPropose(types.StateSubmitMpcType)
			}

			if !util.MpcRewardCommitGenerated {
				mp.checkAndPropose(types.RewardSubmitMpcType)
			}

			if !util.MpcBlobCommitGenerated {
				if blobUpgradeHeight == 0 {
					mp.Logger.Info("Blob upgrade height not set, ignore blob commitment mpc creation")
					continue
				}

				// start creating blob mpc account only when we reached the blob upgrade height
				lastBlockBytes, err := mp.storageClient.Get([]byte(lastMetisBlockKey), nil)
				if err != nil {
					mp.Logger.Error("Unable to fetch last block height from local db, probably not synced, ignore blob commitment mpc creation", "error", err)
					continue
				}
				dbBlockHeight, err := strconv.ParseUint(string(lastBlockBytes), 10, 64)
				if err != nil {
					mp.Logger.Error("Unable to parse last block height from local db, probably not synced, ignore blob commitment mpc creation", "error", err)
					continue
				}
				if dbBlockHeight < blobUpgradeHeight {
					mp.Logger.Info("Blob upgrade height not reached, ignore blob commitment mpc creation", "dbBlockHeight", dbBlockHeight, "blobUpgradeHeight", blobUpgradeHeight)
					continue
				}

				mp.checkAndPropose(types.BlobSubmitMpcType)
			}
		case <-ctx.Done():
			mp.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
}

// checkAndPropose - will check if current user is mpc proposer and proposes the mpc
func (mp *MpcProcessor) checkAndPropose(mpcType types.MpcType) {
	var lastMpc *types.Mpc
	helper.Retry(5, 1*time.Second, func() error {
		var getMpcErr error
		lastMpc, getMpcErr = mp.getLastMpc(mpcType)
		if getMpcErr != nil {
			return getMpcErr
		}
		if lastMpc == nil || lastMpc.ID == "" {
			return errors.New("mpc doesn't really exist")
		}
		return nil
	})

	if lastMpc == nil || lastMpc.ID == "" {
		// mpc not found, create new one
		mp.Logger.Info("checkAndPropose mpc not found")

		if mp.isMpcProposer() {
			lastSpan, _ := mp.getLastSpan()

			// Make sure that each span will only generate mpc once
			hasLocalMpcCreate, _ := mp.storageClient.Has([]byte(localMpcCreateKey+fmt.Sprintf("%v_%v", lastSpan.ID, mpcType)), nil)
			if !hasLocalMpcCreate {
				mp.Logger.Info("checkAndPropose mpc not found, and local mpc create not exist")

				// get current mpc set
				mpcSet, err := mp.getMpcSet()
				if err != nil {
					mp.Logger.Error("Unable to fetch mpcSet", "error", err)
					return
				}
				mp.propose(lastSpan, mpcSet, mpcType)
			}
		}
	} else {
		if mpcType == types.CommonMpcType {
			util.MpcCommonGenerated = true
		} else if mpcType == types.StateSubmitMpcType {
			util.MpcStateCommitGenerated = true
		} else if mpcType == types.RewardSubmitMpcType {
			util.MpcRewardCommitGenerated = true
		} else if mpcType == types.BlobSubmitMpcType {
			util.MpcBlobCommitGenerated = true
		}
		mp.Logger.Info("checkAndPropose mpc found", "mpcId", lastMpc.ID, "mpcAddress", lastMpc.MpcAddress, "checkMpcType", mpcType, "lastMpcType", lastMpc.MpcType)
	}
}

// propose new mpc if needed
func (mp *MpcProcessor) propose(lastSpan *types.Span, mpcSet *types.MpcSet, mpcType types.MpcType) {
	var mpcID string
	var mpcPub []byte
	var err error
	var threshold uint64

	// check mpc recover
	mpcRecover := os.Getenv("MPC_RECOVER")
	if mpcRecover != "" && mpcRecover == "TRUE" {
		switch mpcType {
		case types.CommonMpcType:
			mpcID = os.Getenv("MPC_COMMON_RECOVER_ID")
		case types.StateSubmitMpcType:
			mpcID = os.Getenv("MPC_STATE_RECOVER_ID")
		case types.RewardSubmitMpcType:
			mpcID = os.Getenv("MPC_REWARD_RECOVER_ID")
		case types.BlobSubmitMpcType:
			mpcID = os.Getenv("MPC_BLOB_RECOVER_ID")
		}

		mpcPub, threshold, err = helper.GetMpcKey(mpcID)
		if err != nil {
			mp.Logger.Error("GetMpcKey failed", "err", err)
			return
		}
	} else {
		// mpc id generate
		mpcID = uuid.New().String()
		mp.Logger.Info("✅ Proposing new mpc create", "mpcId", mpcID)

		if len(mpcSet.Result) < 3 {
			threshold = uint64(len(mpcSet.Result) - 1)
		} else {
			threshold = uint64(len(mpcSet.Result) - 2)
		}

		var allPartyId []*tss.PartyID
		for _, set := range mpcSet.Result {
			allPartyId = append(allPartyId,
				&tss.PartyID{
					Id:      set.ID,
					Moniker: set.Moniker,
					Key:     set.Key,
				},
			)
		}

		// call mpc node send keygen request
		mpcPub, err = helper.MpcCreate(&tss.KeyGenRequest{
			KeyId:      mpcID,
			Threshold:  int32(threshold),
			AllPartyId: allPartyId,
		})
		if err != nil {
			mp.Logger.Error("MpcCreate failed", "err", err)
			return
		}
		mp.Logger.Info("MpcCreate success", "pubkey", fmt.Sprintf("0x%x", mpcPub))
	}

	// clac mpc address
	publicKey, err := btcec.ParsePubKey(mpcPub, btcec.S256())
	if err != nil {
		mp.Logger.Error("MpcCreate ParsePubKey failed", "err", err)
		return
	}

	mpcAddress := crypto.PubkeyToAddress(*publicKey.ToECDSA())
	mp.Logger.Info("MpcCreate success", "mpcAddress", mpcAddress.String())

	// log new mpc
	mp.Logger.Info("✅ Proposing new mpc to themis", "mpcId", mpcID)

	// broadcast new mpc info to themis
	msg := mpcTypes.MsgProposeMpcCreate{
		ID:           mpcID,
		Threshold:    threshold,
		Participants: mpcSet.Result,
		Proposer:     types.BytesToThemisAddress(helper.GetAddress()),
		MpcAddress:   types.BytesToThemisAddress(mpcAddress.Bytes()),
		MpcPubkey:    publicKey.SerializeCompressed(),
		MpcType:      mpcType,
	}

	// return broadcast to themis
	if err := mp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
		mp.Logger.Error("Error while broadcasting mpc to themis", "mpcId", mpcID, "error", err)
		return
	}

	// set local key
	mp.setLocalMpcCreate(lastSpan.ID, mpcID, mpcType)
}

func (mp *MpcProcessor) setLocalMpcCreate(spanID uint64, mpcID string, mpcType types.MpcType) {
	if err := mp.storageClient.Put([]byte(localMpcCreateKey+fmt.Sprintf("%v_%v", spanID, mpcType)), []byte(mpcID), nil); err != nil {
		mp.Logger.Error("rl.storageClient.Put", "Error", err)
	}
	mp.Logger.Info("setLocalMpcCreate", "spanID", spanID, "mpcId", mpcID)
}

// checks mpc status
func (mp *MpcProcessor) getLastMpc(mpcType types.MpcType) (*types.Mpc, error) {
	// fetch latest start block from themis via rest query
	result, err := helper.FetchFromAPI(mp.cliCtx, helper.GetThemisServerEndpoint(fmt.Sprintf(util.LatestMpcURL, mpcType)))
	if err != nil {
		mp.Logger.Error("Error while fetching latest mpc", "err", err)
		return nil, err
	}

	var lastMpc types.Mpc
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastMpc); err != nil {
		mp.Logger.Error("Error unmarshalling mpc", "error", err)
		return nil, err
	}
	return &lastMpc, nil
}

func (mp *MpcProcessor) getMpcSet() (*types.MpcSet, error) {
	result, err := helper.FetchFromAPI(mp.cliCtx, helper.GetThemisServerEndpoint(util.MpcSetURL))
	if err != nil {
		mp.Logger.Error("Error while fetching mpc set")
		return nil, err
	}

	var mpcSet types.MpcSet
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &mpcSet.Result); err != nil {
		mp.Logger.Error("Error unmarshalling mpc", "error", err)
		return nil, err
	}
	return &mpcSet, nil
}

// checks span status
func (mp *MpcProcessor) getLastSpan() (*types.Span, error) {
	// fetch latest start block from themis via rest query
	result, err := helper.FetchFromAPI(mp.cliCtx, helper.GetThemisServerEndpoint(util.LatestSpanURL))
	if err != nil {
		mp.Logger.Error("Error while fetching latest span")
		return nil, err
	}

	var lastSpan types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastSpan); err != nil {
		mp.Logger.Error("Error unmarshalling span", "error", err)
		return nil, err
	}
	return &lastSpan, nil
}

// isMpcProposer checks if current user is mpc proposer
func (mp *MpcProcessor) isMpcProposer() bool {
	// get current mpc set
	mpcSet, err := mp.getMpcSet()
	if err != nil {
		mp.Logger.Error("Unable to fetch mpcSet", "error", err)
		return false
	}

	// anyone among next mpc producers can become next mpc proposer
	for _, val := range mpcSet.Result {
		mp.Logger.Info("isMpcProposer", "mpcSetMoniker", val.Moniker, "localAddress", helper.GetAddressStr())
		if strings.EqualFold(val.Moniker, helper.GetAddressStr()) {
			return true
		}
	}
	return false
}

// selectMpcProposer select a mpc proposer
func (mp *MpcProcessor) selectMpcProposer(signID string) string {
	// get current mpc set
	mpcSet, err := mp.getMpcSet()
	if err != nil {
		mp.Logger.Error("Unable to fetch mpcSet", "error", err)
		return ""
	}

	signBigInt, _ := big.NewInt(0).SetString(strings.ReplaceAll(signID, "-", ""), 16)
	rand.Seed(signBigInt.Int64())

	selectIndex := rand.Intn(len(mpcSet.Result))
	mpcProposer := mpcSet.Result[selectIndex].Moniker
	mp.Logger.Debug("selectMpcProposer", "signBigInt", signBigInt, "mpcProposer", mpcProposer)
	return mpcProposer
}

// Stop stops all necessary go routines
func (mp *MpcProcessor) Stop() {
	// cancel mpc polling
	mp.cancelMpcService()
}

func (mp *MpcProcessor) mpcKeySign(eventBytes string) error {
	mp.Logger.Info("Received mpcKeySign request", "eventBytes", eventBytes)

	var event sdk.StringEvent
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(eventBytes), &event); err != nil {
		mp.Logger.Error("Error unmarshalling event from themis", "error", err)
		return err
	}
	mp.Logger.Info("processing mpcKeySign confirmation event", "eventtype", event.Type)

	var (
		signID     string
		mpcID      string
		signMsgStr string
		signMsg    []byte
	)

	for _, attr := range event.Attributes {
		if attr.Key == mpcTypes.AttributeKeyMpcSignID {
			signID = attr.Value
		}

		if attr.Key == mpcTypes.AttributeKeyMpcSignMsg {
			signMsgStr = attr.Value
			signMsg, _ = hex.DecodeString(attr.Value)
		}

		if attr.Key == mpcTypes.AttributeKeyMpcID {
			mpcID = attr.Value
		}
	}

	// filter empty event
	if signID == "" || mpcID == "" || signMsgStr == "" {
		mp.Logger.Error("processing mpcKeySign invalid sign info,ignore it")
		return nil
	}

	selectedMpcProposer := mp.selectMpcProposer(signID)
	if strings.EqualFold(selectedMpcProposer, helper.GetAddressStr()) {
		mp.Logger.Info(
			"✅ Proposing new mpc sign",
			"event", event.Type,
			"signID", signID,
			"mpcID", mpcID,
			"signMsg", signMsgStr,
		)

		// set health value
		rpc.HealthValue.Mpc.IsMpcProposer = 1
		rpc.HealthValue.Mpc.SignID = signID
		rpc.HealthValue.Mpc.Timestamp = time.Now().Unix()

		// call mpc make sign
		signature, err := helper.MpcSign(signID, mpcID, signMsg)
		if err != nil {
			if err == errors.New("KeySignValidate signature already generated") { // signature had been generated
				themisSignInfo, err := getSignInfo(signID)
				if err != nil {
					mp.Logger.Error("mpc sign get exist sign result", err)
					return err
				}
				signature = themisSignInfo.Result.Signature
				rpc.HealthValue.Mpc.SignSuccess = 1
			} else {
				mp.Logger.Error("mpc sign failed", err)
				return err
			}
		}
		rpc.HealthValue.Mpc.SignSuccess = 1

		mp.Logger.Info(
			"mpc sign success",
			"signID", signID,
			"signature", hex.EncodeToString(signature))

		// broadcast to themis
		msg := mpcTypes.MsgMpcSign{
			ID:        signID,
			Signature: hex.EncodeToString(signature),
			Proposer:  types.BytesToThemisAddress(helper.GetAddress()),
		}

		// return broadcast to themis
		if err := mp.txBroadcaster.BroadcastToThemis(msg, event); err != nil {
			mp.Logger.Error("Error while broadcasting mpcSign to themis", "signID", signID, "error", err)
			return err
		}
	}

	return nil
}

func (mp *MpcProcessor) KeySign() {
	allEvents, _ := mp.sqlClient.BridgeSqliteThemisEvent.GetAllWaitPushThemisEventsByType(mpcTypes.EventTypeProposeMpcSign, 100, 0)
	for _, event := range allEvents {
		err := mp.mpcKeySign(event.EventLog)
		if err == nil {
			mp.sqlClient.BridgeSqliteThemisEvent.Delete(event.ID)
		}
	}
}

type ThemisSign struct {
	SignID    string `json:"sign_id,omitempty"`
	MpcID     string `json:"mpc_id,omitempty"`
	SignType  int    `json:"sign_type,omitempty"`
	SignData  []byte `json:"sign_data,omitempty"`
	SignMsg   []byte `json:"sign_msg,omitempty"`
	Proposer  string `json:"proposer,omitempty"`
	Signature []byte `json:"signature,omitempty"`
}

type ThemisResultSignInfo struct {
	Height string     `json:"height,omitempty"`
	Result ThemisSign `json:"result,omitempty"`
}

func getSignInfo(signId string) (*ThemisResultSignInfo, error) {
	url := helper.GetThemisServerEndpoint(fmt.Sprintf(util.MpcSignByIdURL, signId))
	client := &http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var output ThemisResultSignInfo
	err = json.Unmarshal(body, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
