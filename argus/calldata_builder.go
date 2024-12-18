package argus

import (
	"ethers/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

func NewCallDataErc20Transfer(token, to string, amount *big.Int) CallData {
	txData := abi.PackErc20Transfer(to, amount)
	callData := NewCallData(common.HexToAddress(token), big.NewInt(0), txData)
	return callData
}

func NewCallDataErc20Approve(token, spender string, amount *big.Int) CallData {
	txData := abi.PackErc20Approve(spender, amount)
	callData := NewCallData(common.HexToAddress(token), big.NewInt(0), txData)
	return callData
}