package v2

import (
	"fmt"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	v2 "github.com/imskyd/go-ethers/mpc/v2"
	"testing"
)

var mpc *v2.EvmMpcV2

func init() {
	mpc = v2.NewEvmMpcV2(
		coboWaas2.ProdEnv,
		"",
		"",
		"",
	)
	mpc.SetDebug(true)
}

func TestSignMessage(t *testing.T) {
	tx1, err := mpc.GetTransactionByTransactionId("78e70793-68fb-47ec-b2ee-2bbb9933b7fd")
	fmt.Println(err, tx1.Result.TransactionSignatureResult.Signature, tx1.Result.TransactionSignatureResult.ResultType)
}

func TestListWalletAddresses(t *testing.T) {
	res, err := mpc.ListWalletAddresses("62d01110-7599-4221-8010-a6e86f326181")
	fmt.Println(err, res.Data)
}

func TestErc20Transfer(t *testing.T) {
	fee := v2.NewFeeLegacy("ARBITRUM_ETH", "1", "180000")
	txResp, err := mpc.TokenTransfer("ARBITRUM_ETH",
		"0x30724ac9213b6d7df24a0afcfe4d02b57878a18b",
		"0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
		"0x30724ac9213b6d7df24a0afcfe4d02b57878a18b",
		"10000",
		fee,
	)
	fmt.Println(txResp, err)
	txTokenResp, err := mpc.NativeTransfer("ARBITRUM_ETH", "0x30724ac9213b6d7df24a0afcfe4d02b57878a18b", "0x30724ac9213b6d7df24a0afcfe4d02b57878a18b", "0.001", fee)
	fmt.Println(txTokenResp, err)

	resp, err := mpc.CancelTransaction("01fbe88c-3b32-49f3-aa9d-210b7f02b1fc")
	fmt.Println(resp, err)
}
