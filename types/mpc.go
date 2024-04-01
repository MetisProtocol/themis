package types

import (
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	p2pSecp256 "github.com/decred/dcrd/dcrec/secp256k1/v4"
	p2pCrypto "github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type MpcType int

const (
	CommonMpcType       MpcType = iota
	StateSubmitMpcType          // Special type for state submit
	RewardSubmitMpcType         // Special type for reward submit
)

type SignType int

const (
	BatchSubmit SignType = iota
	BatchReward
	CommitEpochToMetis
	ReCommitEpochToMetis
	L1UpdateMpcAddress
	L2UpdateMpcAddress
)

type PartyID struct {
	ID      string `json:"id,omitempty" yaml:"id"`           // mpc p2p peer id
	Moniker string `json:"moniker,omitempty" yaml:"moniker"` // moniker
	Key     []byte `json:"key,omitempty" yaml:"key"`         // mpc p2p peer public key
}

// Mpc stores details for a mpc EOA on Metis chain
type Mpc struct {
	ID           string        `json:"mpc_id" yaml:"mpc_id"`
	Threshold    uint64        `json:"threshold" yaml:"threshold"`
	Participants []PartyID     `json:"participants" yaml:"participants"`
	MpcAddress   ThemisAddress `json:"mpc_address" yaml:"mpc_address"`
	MpcPubkey    []byte        `json:"mpc_pubkey" yaml:"mpc_pubkey"`
	MpcType      MpcType       `json:"mpc_type" yaml:"mpc_type"`
}

type MpcSet struct {
	Result []PartyID `json:"result,omitempty"`
}

// NewMpc creates new mpc
func NewMpc(id string, threshold uint64, parties []PartyID, address ThemisAddress, pubkey []byte, mpcType MpcType) Mpc {
	return Mpc{
		ID:           id,
		Threshold:    threshold,
		Participants: parties,
		MpcAddress:   address,
		MpcPubkey:    pubkey,
		MpcType:      mpcType,
	}
}

// String returns the string representation of mpc
func (s *Mpc) String() string {
	return fmt.Sprintf(
		"Mpc %v %v %v %v %x %v}",
		s.ID,
		s.Threshold,
		s.Participants,
		s.MpcAddress.EthAddress().Hex(),
		s.MpcPubkey,
		s.MpcType,
	)
}

// SortMpcByID sorts mpcs by MpcID
func SortMpcByID(a []*Mpc) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID < a[j].ID
	})
}

// MpcSign creates msg mpc create
type MpcSign struct {
	ID        string        `json:"sign_id" yaml:"sign_id"`
	MpcID     string        `json:"mpc_id" yaml:"mpc_id"`
	SignType  SignType      `json:"sign_type" yaml:"sign_type"`
	SignData  []byte        `json:"sign_data" yaml:"sign_data"`
	SignMsg   []byte        `json:"sign_msg" yaml:"sign_mgs"`
	Proposer  ThemisAddress `json:"proposer" yaml:"proposer"`
	Signature []byte        `json:"signature" yaml:"signature"`
	SignedTx  []byte        `json:"signed_tx" yaml:"signed_tx"`
}

// NewMpcSign creates new mpc create message
func NewMpcSign(
	id string,
	mpcID string,
	signType SignType,
	signData []byte,
	signMsg []byte,
	proposer ThemisAddress,
) MpcSign {
	return MpcSign{
		ID:       id,
		MpcID:    mpcID,
		SignType: signType,
		SignData: signData,
		SignMsg:  signMsg,
		Proposer: proposer,
	}
}

func NewMpcPartyIDFromPrivatekey(cdc *codec.Codec, pris []crypto.PrivKey) ([]PartyID, error) {
	parties := make([]PartyID, len(pris))

	for i, pri := range pris {
		var privObject secp256k1.PrivKeySecp256k1
		cdc.MustUnmarshalBinaryBare(pri.Bytes(), &privObject)

		p2pPrivate := p2pSecp256.PrivKeyFromBytes(pri.Bytes()[4:])
		p2pSecp256Key := (*p2pCrypto.Secp256k1PrivateKey)(p2pPrivate)

		peerId, err := peer.IDFromPrivateKey(p2pSecp256Key)
		if err != nil {
			return nil, err
		}

		p2pPubKeyBytes, err := p2pSecp256Key.GetPublic().Raw()
		if err != nil {
			return nil, err
		}

		moniker := BytesToThemisAddress(privObject.PubKey().Address().Bytes()).String()
		parties[i] = PartyID{
			ID:      peerId.String(),
			Key:     p2pPubKeyBytes,
			Moniker: moniker,
		}
	}
	return parties, nil
}
