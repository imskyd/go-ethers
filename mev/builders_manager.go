package mev

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"math/big"
	"sync"
	"time"
)

type BuildersManager struct {
	client   *ethclient.Client
	builders Builders
	chainId  *big.Int
}

func NewSender(rpc string) *BuildersManager {
	builders := initBuildersFromRemote()
	client, _ := ethclient.Dial(rpc)
	chainId, _ := client.ChainID(context.Background())

	for i, _ := range builders {
		builders[i].init()
		builders[i].setETHClient(client)
	}

	return &BuildersManager{
		client:   client,
		builders: builders,
		chainId:  chainId,
	}
}

func (s *BuildersManager) SendBundle(privateKey *ecdsa.PrivateKey, txs types.Transactions, validBlock uint64) error {
	signTxs, err := s.signTxs(privateKey, txs)
	if err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for _, builder := range s.builders {
		wg.Add(1)
		go func(builder Builder) {
			defer wg.Done()
			if err := builder.sendBundle(privateKey, signTxs, validBlock); err != nil {
				logrus.Errorf("failed to send bundle to builder: %v", err)
			}
		}(builder)
	}
	wg.Wait()

	for _, tx := range signTxs {
		go s.checkTxStatus(tx)
	}

	return nil
}

func (s *BuildersManager) signTxs(privateKey *ecdsa.PrivateKey, txs types.Transactions) (types.Transactions, error) {
	for i, tx := range txs {
		signTx, err := types.SignTx(tx, types.NewCancunSigner(s.chainId), privateKey)
		if err != nil {
			return types.Transactions{}, err
		}
		txs[i] = signTx
	}
	return txs, nil
}

func (s *BuildersManager) checkTxStatus(tx *types.Transaction) {
	maxTryTime := 60
	tried := 0
	for true {
		if tried > maxTryTime {
			logrus.Errorf("tried too many tries to check tx status")
		}
		receipt, err := s.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			logrus.Errorf("failed to get transaction receipt: %v", err)
			tried++
			time.Sleep(time.Second * 1)
			continue
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			logrus.Errorf("tx status failed, check https://etherscan.io/tx/%s", tx.Hash())
			break
		}
		if receipt.Status == types.ReceiptStatusSuccessful {
			logrus.Infof("tx success: https://etherscan.io/tx/%s", tx.Hash())
			break
		}
	}
}
