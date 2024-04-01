package processor

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	authTypes "github.com/metis-seq/themis/auth/types"
	bankTypes "github.com/metis-seq/themis/bank/types"
	"github.com/metis-seq/themis/bridge/setu/util"
	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/helper"
	stakingTypes "github.com/metis-seq/themis/staking/types"
	hmTypes "github.com/metis-seq/themis/types"
)

// StakingProcessor - process staking related events
type StakingProcessor struct {
	BaseProcessor
	stakingInfoAbi *abi.ABI
	cancelService  context.CancelFunc
}

// NewStakingProcessor - add  abi to staking processor
func NewStakingProcessor(stakingInfoAbi *abi.ABI) *StakingProcessor {
	return &StakingProcessor{
		stakingInfoAbi: stakingInfoAbi,
	}
}

// Start starts new block subscription
func (sp *StakingProcessor) Start() error {
	sp.Logger.Info("Starting")

	spCtx, cancelSpanService := context.WithCancel(context.Background())
	sp.cancelService = cancelSpanService

	go sp.StartPolling(spCtx)
	return nil
}

func (sp *StakingProcessor) StartPolling(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// need to broadcast tx to metis
			allEvents, _ := sp.sqlClient.BridgeSqliteEthereumEvent.GetAllWaitPushEthereumEvents(100, 0)
			var err error
			for _, event := range allEvents {
				switch event.EventName {
				case "Locked":
					err = sp.sendValidatorJoinToThemis(event.EventName, event.EventLog)
				case "UnlockInit":
					err = sp.sendUnstakeInitToThemis(event.EventName, event.EventLog)
				case "LockUpdate":
					err = sp.sendStakeUpdateToThemis(event.EventName, event.EventLog)
				case "SignerChange":
					err = sp.sendSignerChangeToThemis(event.EventName, event.EventLog)
				case "BatchSubmitReward":
					err = sp.sendBatchSubmitRewardToThemis(event.EventName, event.EventLog)
				case "UpdateEpochLength":
					err = sp.sendUpdateEpochLength(event.EventName, event.EventLog)
				}
				if err == nil {
					// delete event
					sp.sqlClient.BridgeSqliteEthereumEvent.Delete(uint64(event.ID))
				}
			}
		case <-ctx.Done():
			sp.Logger.Info("Polling stopped")
			ticker.Stop()
			return
		}
	}
}

// RegisterTasks - Registers staking tasks with machinery
func (sp *StakingProcessor) RegisterTasks() {
}

// Stop stops all necessary go routines
func (sp *StakingProcessor) Stop() {
	sp.cancelService()
}

func (sp *StakingProcessor) sendValidatorJoinToThemis(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoLocked)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		// check auto gas send
		isGasSender := os.Getenv("IS_GAS_SENDER")
		if isGasSender == "true" {
			sp.autoSendGas(event.Signer)
		}

		// check signer
		pubkey := helper.GetPubKey()
		if !bytes.Equal(event.SignerPubkey, pubkey[1:]) {
			sp.Logger.Info("not self lock event,ignore it.")
			return nil
		}

		signerPubKey := event.SignerPubkey
		if len(signerPubKey) == 64 {
			signerPubKey = util.AppendPrefix(signerPubKey)
		}
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send validatorjoin to themis as already processed",
				"event", eventName,
				"validatorID", event.SequencerId,
				"activationBatch", event.ActivationBatch,
				"nonce", event.Nonce,
				"amount", event.Amount,
				"totalAmount", event.Total,
				"SignerPubkey", hmTypes.NewPubKey(signerPubKey).String(),
				"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		// if account doesn't exists Retry with delay for topup to process first.
		if _, err := util.GetAccount(sp.cliCtx, hmTypes.ThemisAddress(event.Signer)); err != nil {
			sp.Logger.Info(
				"Themis Account doesn't exist. Retrying validator-join after 10 seconds",
				"event", eventName,
				"signer", event.Signer,
			)
			return errors.New("account doesn't exist")
		}

		sp.Logger.Info(
			"Received task to send validatorjoin to themis",
			"event", eventName,
			"validatorID", event.SequencerId,
			"activationBatch", event.ActivationBatch,
			"nonce", event.Nonce,
			"amount", event.Amount,
			"totalAmount", event.Total,
			"SignerPubkey", hmTypes.NewPubKey(signerPubKey).String(),
			"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		msg := stakingTypes.NewMsgValidatorJoin(
			hmTypes.BytesToThemisAddress(helper.GetAddress()),
			event.SequencerId.Uint64(),
			event.ActivationBatch.Uint64(),
			sdk.NewIntFromBigInt(event.Amount),
			hmTypes.NewPubKey(signerPubKey),
			hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to themis
		if err := sp.txBroadcaster.BroadcastToThemis(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting unstakeInit to themis", "validatorId", event.SequencerId.Uint64(), "error", err)
			return err
		}
	}

	return nil
}

func (sp *StakingProcessor) sendUnstakeInitToThemis(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoUnlockInit)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to themis as already processed",
				"event", eventName,
				"validator", event.User,
				"validatorID", event.SequencerId,
				"nonce", event.Nonce,
				"deactivatonBatch", event.DeactivationBatch,
				"amount", event.Amount,
				"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		sp.Logger.Info(
			"Received task to send unstake-init to themis",
			"event", eventName,
			"validator", event.User,
			"validatorID", event.SequencerId,
			"nonce", event.Nonce,
			"deactivatonBatch", event.DeactivationBatch,
			"amount", event.Amount,
			"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		msg := stakingTypes.NewMsgValidatorExit(
			hmTypes.BytesToThemisAddress(helper.GetAddress()),
			event.SequencerId.Uint64(),
			event.DeactivationBatch.Uint64(),
			hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to themis
		if err := sp.txBroadcaster.BroadcastToThemis(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting unstakeInit to themis", "validatorId", event.SequencerId.Uint64(), "error", err)
			return err
		}
	}

	return nil
}

func (sp *StakingProcessor) sendStakeUpdateToThemis(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoLockUpdate)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to themis as already processed",
				"event", eventName,
				"validatorID", event.SequencerId,
				"nonce", event.Nonce,
				"newAmount", event.NewAmount,
				"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		sp.Logger.Info(
			"Received task to send stake-update to themis",
			"event", eventName,
			"validatorID", event.SequencerId,
			"nonce", event.Nonce,
			"newAmount", event.NewAmount,
			"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		msg := stakingTypes.NewMsgStakeUpdate(
			hmTypes.BytesToThemisAddress(helper.GetAddress()),
			event.SequencerId.Uint64(),
			sdk.NewIntFromBigInt(event.NewAmount),
			hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to themis
		if err := sp.txBroadcaster.BroadcastToThemis(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting stakeupdate to themis", "validatorId", event.SequencerId.Uint64(), "error", err)
			return err
		}
	}

	return nil
}

func (sp *StakingProcessor) sendSignerChangeToThemis(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoSignerChange)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		newSignerPubKey := event.SignerPubkey
		if len(newSignerPubKey) == 64 {
			newSignerPubKey = util.AppendPrefix(newSignerPubKey)
		}

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to themis as already processed",
				"event", eventName,
				"validatorID", event.SequencerId,
				"nonce", event.Nonce,
				"NewSignerPubkey", hmTypes.NewPubKey(newSignerPubKey).String(),
				"oldSigner", event.OldSigner.Hex(),
				"newSigner", event.NewSigner.Hex(),
				"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		sp.Logger.Info(
			"Received task to send signer-change to themis",
			"event", eventName,
			"validatorID", event.SequencerId,
			"nonce", event.Nonce,
			"NewSignerPubkey", hmTypes.NewPubKey(newSignerPubKey).String(),
			"oldSigner", event.OldSigner.Hex(),
			"newSigner", event.NewSigner.Hex(),
			"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// signer change
		msg := stakingTypes.NewMsgSignerUpdate(
			hmTypes.BytesToThemisAddress(helper.GetAddress()),
			event.SequencerId.Uint64(),
			hmTypes.NewPubKey(newSignerPubKey),
			hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to themis
		if err := sp.txBroadcaster.BroadcastToThemis(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting signerChange to themis", "msg", msg, "validatorId", event.SequencerId.Uint64(), "error", err)
			return err
		}
	}

	return nil
}

func (sp *StakingProcessor) sendBatchSubmitRewardToThemis(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoBatchSubmitReward)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		currentBatch, err := util.GetCurrentBatch(sp.cliCtx)
		if err != nil {
			sp.Logger.Error("Error GetCurrentBatch", "error", err)
			return nil
		}
		if currentBatch >= event.NewBatchId.Uint64() {
			sp.Logger.Info("Ignoring task to send batchSubmitReward to themis", "chainBatch", currentBatch, "eventBatch", event.NewBatchId.Uint64())
			return nil
		}

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send batchSubmitReward to themis as already processed",
				"event", eventName,
				"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		sp.Logger.Info(
			"Received task to send batch-submit-reward to themis",
			"event", eventName,
			"txHash", hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		validatorId := sp.queryValidatorID()
		if validatorId <= 0 {
			sp.Logger.Info("Not a validator,skip it", "validatorId", validatorId)
			return nil
		}
		currentNonce, _, err := util.GetValidatorNonce(sp.cliCtx, validatorId)
		if err != nil {
			sp.Logger.Error("Failed to fetch validator nonce and height data from API", "validatorId", validatorId)
			return err
		}

		// batch submit reward
		msg := stakingTypes.NewMsgBatchSubmitReward(
			hmTypes.BytesToThemisAddress(helper.GetAddress()),
			validatorId,
			event.NewBatchId.Uint64(),
			hmTypes.BytesToThemisHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			currentNonce+1,
		)

		// return broadcast to themis
		if err := sp.txBroadcaster.BroadcastToThemis(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting submitReward to themis", "msg", msg, "error", err)
			return err
		}
	}

	return nil
}

func (sp *StakingProcessor) sendUpdateEpochLength(eventName string, logBytes string) error {
	return nil
}

func (sp *StakingProcessor) queryValidatorID() uint64 {
	valSet, err := util.GetValidatorSet(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error GetValidatorSet", "error", err)
		return 0
	}

	for _, val := range valSet.Validators {
		if bytes.Equal(val.Signer.Bytes(), helper.GetAddress()) {
			return uint64(val.ID)
		}
	}

	return 0
}

func (sp *StakingProcessor) autoSendGas(receipent common.Address) error {
	// get account getter
	accGetter := authTypes.NewAccountRetriever(sp.cliCtx)

	// get from account
	from := helper.GetFromAddress(sp.cliCtx)

	// to key
	to := hmTypes.HexToThemisAddress(receipent.Hex())
	if to.Empty() {
		return fmt.Errorf("Invalid to address")
	}

	if err := accGetter.EnsureExists(from); err != nil {
		return err
	}

	account, err := accGetter.GetAccount(from)
	if err != nil {
		return err
	}

	sendCoins := "1000000000000000000000themis"
	// parse coins trying to be sent
	coins, err := sdk.ParseCoins(sendCoins)
	if err != nil {
		return err
	}

	// ensure account has enough coins
	if !account.GetCoins().IsAllGTE(coins) {
		return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction,Required coin= %s Available Coins=%s", from, coins, account.GetCoins())
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := bankTypes.NewMsgSend(from, to, coins)
	// return broadcast to themis
	if err := sp.txBroadcaster.BroadcastToThemis(msg, nil); err != nil {
		sp.Logger.Error("Error while broadcasting signerChange to themis", "msg", msg, "error", err)
		return err
	}
	return nil
}
