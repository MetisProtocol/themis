package mpc

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/metis-seq/themis/mpc/types"
)

func TestQueryParams(t *testing.T) {
	data := []byte{123, 34, 83, 105, 103, 110, 73, 68, 34, 58, 34, 102, 48, 102, 55, 55, 50, 56, 98, 45, 98, 52, 97, 49, 45, 52, 98, 52, 100, 45, 57, 101, 98, 50, 45, 49, 98, 48, 48, 49, 57, 48, 98, 49, 98, 98, 51, 34, 125}
	cdc := MakeCodec()

	var params types.QueryMpcSignParams
	if err := cdc.UnmarshalJSON(data, &params); err != nil {
		t.Fatal(err)
	}

	t.Logf("params:%v", params)
}

func MakeCodec() *codec.Codec {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	cdc.Seal()
	return cdc
}

func TestSignatureVerify(t *testing.T) {
	pubkey := "02b8351bea9a451b74193b25f371c54ca363fa37c3231e194234ee4726de551ef3"
	pubkeyBytes, _ := hex.DecodeString(pubkey)
	pubKey, err := btcec.ParsePubKey(pubkeyBytes, btcec.S256())
	if err != nil {
		t.Fatal(err)
	}

	publicKeyBytes := crypto.FromECDSAPub(pubKey.ToECDSA())

	signature := "bde0bf9d75751c452bae4692b5f365a771d2449b7373dd29ff832ddf369df59227faba87f178a6f076d80714036cb0151dc8b7850272f0e015dd2657c9a2ac831c"
	sig, _ := hex.DecodeString(signature)

	signMsg := []byte("hello")
	hash := crypto.Keccak256Hash(signMsg)
	t.Logf("sig hash: %x", hash.Bytes())

	t.Logf("signature len: %v", len(sig))

	t.Logf("signature v: %v", sig[64])
	if sig[64] > 1 {
		sig[64] -= 27
	}
	t.Logf("signature v: %v", sig[64])
	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), sig)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(sigPublicKey, publicKeyBytes) {
		t.Fatalf("verify failed")
	}

	t.Logf("signature verify success")
}
