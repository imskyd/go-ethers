package mev

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req/v3"
	"github.com/metachris/flashbotsrpc"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"strings"
)

type Builders []Builder
type Builder struct {
	Name          string   `json:"name"`
	Rpc           string   `json:"rpc"`
	SupportedApis []string `json:"supported-apis"`
	Client        *flashbotsrpc.FlashbotsRPC
	ETHClient     *ethclient.Client
	ChainId       *big.Int
}

const GithubFlashBotsBuilders = "https://raw.githubusercontent.com/flashbots/dowg/refs/heads/main/builder-registrations.json"

func initBuildersFromRemote() Builders {
	var builders Builders
	resp, err := req.Get(GithubFlashBotsBuilders)
	if err != nil {
		log.Fatalf("Error getting builders: %v", err)
	}
	if err := resp.UnmarshalJson(&builders); err != nil {
		log.Fatalf("Error unmarshalling builders: %v", err)
	}

	for i, b := range builders {
		if ok := strings.Contains(b.Rpc, "http"); !ok {
			builders[i].Rpc = "https://" + b.Rpc
		}

	}
	return builders
}

func (b *Builder) init() {
	rpc := flashbotsrpc.New(b.Rpc)
	b.Client = rpc
}

func (b *Builder) setETHClient(client *ethclient.Client) {
	b.ETHClient = client
	chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain id: %v", err)
	}
	b.ChainId = chainId
}

func (b *Builder) sendBundle(privateKey *ecdsa.PrivateKey, txs types.Transactions, validBlock uint64) error {
	logrus.Infof("send bundle to builder %s", b.Name)

	var hexTxs []string

	for _, tx := range txs {
		txBytes, _ := tx.MarshalBinary()
		hexTxs = append(hexTxs, "0x"+common.Bytes2Hex(txBytes))
	}

	sendBundleArgs := flashbotsrpc.FlashbotsSendBundleRequest{
		Txs:         hexTxs,
		BlockNumber: fmt.Sprintf("0x%x", validBlock),
	}

	result, err := b.Client.FlashbotsSendBundle(privateKey, sendBundleArgs)
	if err != nil {
		logrus.Errorf("builder %s err: %s", b.Name, err.Error())
		return err
	}

	if result.BundleHash != "" {
		return nil
	}
	return nil
}
