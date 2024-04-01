package listener

import (
	"context"
	"math/big"
	"strings"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/helper"
	"github.com/stretchr/testify/assert"
)

func TestFilterLog(t *testing.T) {
	helper.InitThemisConfig("../../../data")

	ethURL := helper.GetConfig().EthRPCUrl
	t.Logf("eth url:%v", ethURL)

	ethURL = "https://rpc.ankr.com/eth_sepolia"

	client, err := ethclient.Dial(ethURL)
	if err != nil {
		panic(err)
	}

	logs, err := client.FilterLogs(context.TODO(), ethereum.FilterQuery{
		FromBlock: big.NewInt(5603203),
		ToBlock:   big.NewInt(5603204),
		Addresses: []common.Address{
			common.HexToAddress("0xaB260022803e6735a81256604d73115D663c6b82"),
		},
	})
	if err != nil {
		panic(err)
	}

	abiObj, _ := abi.JSON(strings.NewReader(stakinginfo.StakinginfoABI))

	t.Logf("logs len:%v", len(logs))
	for _, vLog := range logs {
		t.Logf("log hash:%v", vLog.TxHash.Hex())
		t.Logf("log index:%v", vLog.Index)

		topic := vLog.Topics[0].Bytes()
		selectedEvent := helper.EventByID(&abiObj, topic)
		if selectedEvent == nil {
			continue
		}
		t.Logf("event:%v", selectedEvent)
		assert.Equal(t, "Locked", selectedEvent.Name)
	}
}
