package argus

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Service struct {
	rpc        string
	privateKey *ecdsa.PrivateKey
	contract   *Argus
	chainId    *big.Int
}

func NewArgusService(rpc, argusAddress, botPk string) (*Service, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}
	_privateKey, err := crypto.HexToECDSA(botPk)
	if err != nil {
		return nil, err
	}
	argusContract, err := NewArgus(common.HexToAddress(argusAddress), client)
	if err != nil {
		return nil, err
	}
	_chainId, _ := client.ChainID(context.Background())
	argus := &Service{
		rpc:        rpc,
		chainId:    _chainId,
		privateKey: _privateKey,
		contract:   argusContract,
	}
	return argus, nil
}

func NewCallData(to common.Address, value *big.Int, data []byte) CallData {
	cd := CallData{
		Flag:  big.NewInt(0),
		To:    to,
		Value: value,
		Data:  data,
		Hint:  common.Hex2Bytes("0x"),
		Extra: common.Hex2Bytes("0x"),
	}
	return cd
}

func (a *Service) GetTransactOpts() (*bind.TransactOpts, error) {
	o, err := bind.NewKeyedTransactorWithChainID(a.privateKey, a.chainId)
	return o, err
}

func (a *Service) ExecTransaction(opts *bind.TransactOpts, callData CallData) (*types.Transaction, error) {
	if opts == nil {
		o, _ := bind.NewKeyedTransactorWithChainID(a.privateKey, a.chainId)
		opts = o
	}
	tx, err := a.contract.ExecTransaction(opts, callData)
	return tx, err
}

func (a *Service) ExecTransactions(opts *bind.TransactOpts, callDataList []CallData) (*types.Transaction, error) {
	if opts == nil {
		o, _ := bind.NewKeyedTransactorWithChainID(a.privateKey, a.chainId)
		opts = o
	}
	tx, err := a.contract.ExecTransactions(opts, callDataList)
	return tx, err
}
