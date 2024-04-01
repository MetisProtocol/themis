package types

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestPubKey(t *testing.T) {
	pubkeyValue := "BPdH9NyFrdWJSaCP6v9jGk2B9/7EAqGpseBidYTHZZjC5SAnKFy9Z9KoaGaTIRWq7Ux4eWZxK4TERZ6U/dIA8ZA="
	pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkeyValue)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("pubkey hex: %x", pubkeyBytes)
	// 04f747f4dc85add58949a08feaff631a4d81f7fec402a1a9b1e0627584c76598c2e52027285cbd67d2a86866932115aaed4c787966712b84c4459e94fdd200f190

	pubkey := NewPubKey(pubkeyBytes)
	address := pubkey.Address()
	t.Logf("pubkey address:%v", address)
}

func TestSignedTx(t *testing.T) {
	signedTxValue := "+QESAoUG/COsAIMEk+CUklfZ1Hj7cbmMwtGGaxqMUEqLZMeAuKqozaN7AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAlcAAAAAjwAABgAABQAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAY0evSAAAdmvhAAABAAAAAGNHtJ0AAHZsOwAAAgAAAABjR751AAB2bOAAAAIAAAAAY0e/KgAAdmzsABg9EDaSsAACaKq7LrY+sI5LrhedV+E1iI3S5hmnaBtfWF2d9uRBC5rvFIL1EKDIoI2RMaMI0FOzf1fBlsyMTW1S0Wq4SDKyeggYBXJNAaBLxlQ+chp9byiJa3Yl18TuDi02srQtohOmbiQMncmBpg=="
	signedTxBytes, err := base64.StdEncoding.DecodeString(signedTxValue)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("signed hex: 0x%x", signedTxBytes)
}

func TestPri2Pub(t *testing.T) {
	pri := "0x1f9a552c0aad1f104401316375f0737bba5fba0b34a83b0069f2a02c57514a0c"
	priBytes, _ := hex.DecodeString(strings.TrimPrefix(pri, "0x"))
	t.Logf("priBytes len:%v", len(priBytes))

	privKeyBytes := [32]byte{}
	for i := 0; i < 32; i++ {
		privKeyBytes[i] = priBytes[i]
	}

	priObj := secp256k1.PrivKeySecp256k1(privKeyBytes)
	t.Logf("prikey: 0x%x", priObj)

	pubkeyHex := hex.EncodeToString(priObj.PubKey().Bytes())
	t.Logf("pubkey: 0x%v", strings.TrimPrefix(pubkeyHex, "eb5ae9874104"))

	address := priObj.PubKey().Address()
	t.Logf("pubkey address: 0x%v", address)
}

func TestNewKey(t *testing.T) {
	priObj := secp256k1.GenPrivKey()
	pubkeyHex := hex.EncodeToString(priObj.PubKey().Bytes())
	t.Logf("pubBytes :%v, len:%v", pubkeyHex, len(pubkeyHex))

	t.Logf("prikey: 0x%x", priObj)
	t.Logf("pubkey: 0x%v", strings.TrimPrefix(pubkeyHex, "eb5ae9874104"))

	address := priObj.PubKey().Address()
	t.Logf("pubkey address: 0x%v", strings.ToLower(address.String()))
}

func TestStrToLower(t *testing.T) {
	str := "6A1d21F068f2456d620A2FCbB018C990b9D1b1bd"
	t.Logf("lower str: %v", strings.ToLower(str))

	str = "c5F6997BA8D5784698e0228985B41ae8b00963A0"
	t.Logf("lower str: %v", strings.ToLower(str))

	str = "d75b2b9B92e4d136B0518A5fE93ab51c80347bAD"
	t.Logf("lower str: %v", strings.ToLower(str))

	str = "D17c31878C3feAE00e675CbEe51D4023EDC90503"
	t.Logf("lower str: %v", strings.ToLower(str))
}
