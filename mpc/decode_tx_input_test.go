package mpc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/metis-seq/themis/contracts/stakemanager"
)

func DecodeTransactionInputData1(contractABI *abi.ABI, data []byte) (map[string]interface{}, error) {
	methodSigData := data[:4]
	inputsSigData := data[4:]
	fmt.Printf("mentod id: %x\n", methodSigData)
	method, err := contractABI.MethodById(methodSigData)
	if err != nil {
		fmt.Printf("MethodById err: %v\n", err)
		return nil, err
	}

	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		fmt.Printf("UnpackIntoMap err: %v\n", err)
		return nil, err
	}

	return inputsMap, nil
}

func TestTxDataDecode(t *testing.T) {
	txDataStr := "02f9024a0581df58708302113594bf10ba6759fb74e251c0acb89686382b0daaa50e80b901e4a2512d790000000000000000000000000000000000000000000000000000000000000003000000000000000000000000b4ebe166513c578e33a8373f04339508bc7e8cfb0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000000010000000000000000000000001267397fb5bf6f6dcc3d18d673616d512dbcd8f0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000041b858bb278dd7606237e376801099195c3e4f558727c5be48058c2ff2296c77697b3cb114b6f61c836acafae104db3d5dca994f9afa8e6322756a15bcb55996ab1c00000000000000000000000000000000000000000000000000000000000000c080a0e4159c0a8041e1239548c4337e606dfef2d969926430eb1a80bc36ebc7a25e08a070fe173bffaeb8e7140f84f1d39f21b014d0842103d851a1baae2680ee298690"
	txData, _ := hex.DecodeString(txDataStr)

	tx := new(types.Transaction)
	err := tx.UnmarshalBinary(txData)
	if err != nil {
		t.Fatal(err)
	}

	stakingAbi, _ := abi.JSON(strings.NewReader(stakemanager.StakemanagerABI))
	result, err := DecodeTransactionInputData1(&stakingAbi, tx.Data())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("decode tx data result:%v", result)
}

func TestBytesEqual(t *testing.T) {
	str := "b926b53ae92ba7ffc2d7b1937c94d2be5eacff2180f7133da8cc08a3591ea4b6"
	data, _ := hex.DecodeString(str)
	result := bytes.Equal(data, data)
	t.Log(result)
}

func TestUnmarshalTx(t *testing.T) {
	signData := "0xf8ec808504a817c800830f42409499d02ce7160109bdda8bf81629eea80661a2150980b8844fb71bdd000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000009600000000000000000000000000000000000000000000000000000000000000a270000000000000000000000001267397fb5bf6f6dcc3d18d673616d512dbcd8f08204eda086a59fe1b95e21854208d9902d99b41e27dd35743c59ed9f67544ba53c499147a00695f186625f972f0ea97f8313fb68c9489096e2ecffd200f6d4e743f0a7b204"

	data, err := hex.DecodeString(strings.TrimPrefix(signData, "0x"))
	if err != nil {
		t.Fatal(err)
	}

	tx := new(types.Transaction)
	err = tx.UnmarshalBinary(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("sign data is ok")

	jsonData, _ := tx.MarshalJSON()
	t.Logf("json data is:%v", string(jsonData))

	// broadcast tx
	// conn, _ := ethclient.Dial("http://localhost:8545")
	// err = conn.SendTransaction(context.Background(), tx)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log("send tx success", tx.Hash().Hex())

	// wait for the transaction be mined
	// _, err = bind.WaitMined(context.Background(), conn, tx)
	// if err != nil {
	// 	t.Fatalf("WaitMined tx err:%v", err)
	// 	return
	// }
}
