package listener

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/metis-seq/themis/contracts/stakinginfo"
	"github.com/metis-seq/themis/helper"
)

// handleLog handles the given log
func (rl *RootChainListener) handleLog(vLog types.Log, selectedEvent *abi.Event) error {
	rl.Logger.Info("ReceivedEvent", "eventname", selectedEvent.Name)

	switch selectedEvent.Name {
	case "Locked":
		return rl.handleStakedLog(vLog, selectedEvent)
	case "LockUpdate":
		return rl.handleStakeUpdateLog(vLog, selectedEvent)
	case "SignerChange":
		return rl.handleSignerChangeLog(vLog, selectedEvent)
	case "UnlockInit":
		return rl.handleUnstakeInitLog(vLog, selectedEvent)
	case "BatchSubmitReward":
		return rl.handleBatchSubmitReward(vLog, selectedEvent)
	default:
		rl.Logger.Info("Unhandled event", "eventname", selectedEvent.Name)
	}
	return nil
}

func (rl *RootChainListener) handleStakedLog(vLog types.Log, selectedEvent *abi.Event) error {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
		return err
	}
	event := new(stakinginfo.StakinginfoLocked)
	if err = helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
		rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		return err
	}
	return rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleStakeUpdateLog(vLog types.Log, selectedEvent *abi.Event) error {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
		return err
	}
	return rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleSignerChangeLog(vLog types.Log, selectedEvent *abi.Event) error {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
		return err
	}
	return rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleUnstakeInitLog(vLog types.Log, selectedEvent *abi.Event) error {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
		return err
	}
	return rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) handleBatchSubmitReward(vLog types.Log, selectedEvent *abi.Event) error {
	logBytes, err := jsoniter.ConfigFastest.Marshal(vLog)
	if err != nil {
		rl.Logger.Error("Failed to marshal log", "Error", err)
		return err
	}
	return rl.insertToDb(selectedEvent.Name, string(logBytes))
}

func (rl *RootChainListener) insertToDb(eventName, eventLog string) error {
	_, err := rl.sqlClient.BridgeSqliteEthereumEvent.Insert(eventName, eventLog)
	return err
}
