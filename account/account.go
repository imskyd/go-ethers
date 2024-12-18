package account

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Account struct {
	pk         string
	PrivateKey *ecdsa.PrivateKey
	address    common.Address
}

func NewAccount(pk string) *Account {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		panic(err)
	}

	return &Account{
		pk:         pk,
		PrivateKey: privateKey,
		address:    crypto.PubkeyToAddress(privateKey.PublicKey),
	}
}

func (a *Account) Address() common.Address {
	return a.address
}

func (a *Account) SignTx(chainId int64, tx *types.Transaction) (*types.Transaction, error) {
	return types.SignTx(tx, types.NewCancunSigner(big.NewInt(chainId)), a.PrivateKey)
}
