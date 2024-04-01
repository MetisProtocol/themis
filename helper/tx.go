package helper

import (
	"context"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/metis-seq/themis/contracts/erc20"
	"github.com/metis-seq/themis/contracts/stakemanager"
)

func GenerateAuthObj(client *ethclient.Client, address common.Address, data []byte) (auth *bind.TransactOpts, err error) {
	// generate call msg
	callMsg := ethereum.CallMsg{
		To:   &address,
		Data: data,
	}

	// get priv key
	pkObject := GetPrivKey()

	// create ecdsa private key
	ecdsaPrivateKey, err := crypto.ToECDSA(pkObject[:])
	if err != nil {
		return
	}

	// from address
	fromAddress := common.BytesToAddress(pkObject.PubKey().Address().Bytes())
	// fetch gas price
	gasprice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return
	}

	mainChainMaxGasPrice := GetConfig().MainchainMaxGasPrice
	// Check if configured or not, Use default in case of invalid value
	if mainChainMaxGasPrice <= 0 {
		mainChainMaxGasPrice = DefaultMainchainMaxGasPrice
	}

	if gasprice.Cmp(big.NewInt(mainChainMaxGasPrice)) == 1 {
		Logger.Error("Gas price is more than max gas price", "gasprice", gasprice)
		err = fmt.Errorf("gas price is more than max_gas_price, gasprice = %v, maxGasPrice = %d", gasprice, mainChainMaxGasPrice)

		return
	}

	// fetch nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return
	}

	// fetch gas limit
	callMsg.From = fromAddress
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		Logger.Error("Unable to fetch ChainID", "error", err)
		return
	}

	// create auth
	auth, err = bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainId)
	if err != nil {
		Logger.Error("Unable to create auth object", "error", err)
		return
	}

	auth.GasPrice = gasprice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = gasLimit

	return
}

// StakeFor stakes for a validator
func (c *ContractCaller) StakeFor(val common.Address, stakeAmount *big.Int, feeAmount *big.Int, acceptDelegation bool, stakeManagerAddress common.Address, stakeManagerInstance *stakemanager.Stakemanager) error {
	signerPubkey := GetPubKey()
	signerPubkeyBytes := signerPubkey[1:] // remove 04 prefix

	// pack data based on method definition
	data, err := c.StakeManagerABI.Pack("stakeFor", val, stakeAmount, feeAmount, acceptDelegation, signerPubkeyBytes)
	if err != nil {
		Logger.Error("Unable to pack tx for stakeFor", "error", err)
		return err
	}

	auth, err := GenerateAuthObj(GetMainClient(), stakeManagerAddress, data)
	if err != nil {
		Logger.Error("Unable to create auth object", "error", err)
		return err
	}

	// stake for stake manager
	tx, err := stakeManagerInstance.LockFor(
		auth,
		val,
		stakeAmount,
		// acceptDelegation,
		signerPubkeyBytes,
	)

	if err != nil {
		Logger.Error("Error while submitting stake", "error", err)
		return err
	}

	Logger.Info("Submitted stake successfully", "txHash", tx.Hash().String())

	return nil
}

// ApproveTokens approves metis token for stake
func (c *ContractCaller) ApproveTokens(amount *big.Int, stakeManager common.Address, tokenAddress common.Address, metisTokenInstance *erc20.Erc20) error {
	data, err := c.MetisTokenABI.Pack("approve", stakeManager, amount)
	if err != nil {
		Logger.Error("Unable to pack tx for approve", "error", err)
		return err
	}

	auth, err := GenerateAuthObj(GetMainClient(), tokenAddress, data)
	if err != nil {
		Logger.Error("Unable to create auth object", "error", err)
		return err
	}

	tx, err := metisTokenInstance.Approve(auth, stakeManager, amount)
	if err != nil {
		Logger.Error("Error while approving approve", "error", err)
		return err
	}

	Logger.Info("Sent approve tx successfully", "txHash", tx.Hash().String())

	return nil
}
