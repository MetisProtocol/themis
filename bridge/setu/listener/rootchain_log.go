package listener

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/helper"
)

// handleLog handles the given log
func (rl *RootChainListener) handleLog(vLog types.Log, selectedEvent *abi.Event) {
	rl.Logger.Debug("ReceivedEvent", "eventname", selectedEvent.Name)

	switch selectedEvent.Name {
	case "Locked":
		rl.handleStakedLog(vLog, selectedEvent)
	case "LockUpdate":
		rl.handleStakeUpdateLog(vLog, selectedEvent)
	case "SignerChange":
		rl.handleSignerChangeLog(vLog, selectedEvent)
	case "UnlockInit":
		rl.handleUnstakeInitLog(vLog, selectedEvent)
	case "BatchSubmitReward":
		rl.handleBatchSubmitReward(vLog, selectedEvent)
	}
}

func (rl *RootChainListener) handleStakedLog(vLog types.Log, selectedEvent *abi.Event) {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
	}

	event := new(stakinginfo.StakinginfoLocked)
	if err = helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
		rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
	}

	rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleStakeUpdateLog(vLog types.Log, selectedEvent *abi.Event) {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
	}
	rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleSignerChangeLog(vLog types.Log, selectedEvent *abi.Event) {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
	}
	rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleUnstakeInitLog(vLog types.Log, selectedEvent *abi.Event) {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
	}
	rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleBatchSubmitReward(vLog types.Log, selectedEvent *abi.Event) {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
	}

	rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) insertToDb(eventName, eventLog string) error {
	_, err := rl.sqlClient.BridgeSqliteEthereumEvent.Insert(eventName, eventLog)
	return err
}
