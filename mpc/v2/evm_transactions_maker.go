package v2

import (
	"encoding/json"
	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imskyd/go-ethers/abi"
	"github.com/shopspring/decimal"
)

// TokenTransfer basic token transfer
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) ApiTokenTransfer(from, to, tokenId, amount string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	//source
	mpcSource := coboWaas2.NewMpcTransferSource(coboWaas2.WALLETSUBTYPE_ORG_CONTROLLED, m.walletId)
	mpcSource.SetAddress(from)
	source := coboWaas2.MpcTransferSourceAsTransferSource(mpcSource)
	//des
	des := coboWaas2.AddressTransferDestinationAsTransferDestination(coboWaas2.NewAddressTransferDestination(coboWaas2.TRANSFERDESTINATIONTYPE_ADDRESS))
	output := coboWaas2.NewAddressTransferDestinationAccountOutput(to, amount)
	des.AddressTransferDestination.AccountOutput = output
	//param
	transferParams := *coboWaas2.NewTransferParams(
		m.createRequestId(),
		source,
		tokenId,
		des,
	)
	if fee != nil {
		transferParams.Fee = fee
	}

	if m.debug {
		debugJson, _ := json.Marshal(transferParams)
		m.logger.Infof("token transfer parameters: %s", debugJson)
	}

	resp, r, err := m.client.TransactionsAPI.CreateTransferTransaction(m.getCtx()).TransferParams(transferParams).Execute()
	return m.formatResponse(resp, r, err)
}

// TokenTransfer basic contract call
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) TokenTransfer(chainId, from, token, receiver, transferAmount string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	amount, err := decimal.NewFromString(transferAmount)
	if err != nil {
		return nil, err
	}
	bytes := abi.PackErc20Transfer(receiver, amount.BigInt())
	resp, err := m.ContractCall(chainId, from, token, "0x"+common.Bytes2Hex(bytes), "", fee)
	return resp, err
}

// TokenApprove basic contract call
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) TokenApprove(chainId, from, token, spender, approveAmount string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	amount, err := decimal.NewFromString(approveAmount)
	if err != nil {
		return nil, err
	}
	bytes := abi.PackErc20Approve(spender, amount.BigInt())
	resp, err := m.ContractCall(chainId, from, token, "0x"+common.Bytes2Hex(bytes), "", fee)
	return resp, err
}

// NativeTransfer
// value is formatted, exam: if you want transfer 0.1 eth, value is "0.1"
func (m *EvmMpcV2) NativeTransfer(chainId, from, to, value string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	return m.ContractCall(chainId, from, to, "0x", value, fee)
}

// ContractCall basic contract call
// if fee is nil, use cobo calc fee param
func (m *EvmMpcV2) ContractCall(chainId, from, to, callData, value string, fee *coboWaas2.TransactionRequestFee) (*coboWaas2.CreateTransferTransaction201Response, error) {
	//source
	mpcSource := coboWaas2.NewMpcContractCallSource(
		coboWaas2.CONTRACTCALLSOURCETYPE_ORG_CONTROLLED,
		m.walletId,
		from,
	)
	source := coboWaas2.MpcContractCallSourceAsContractCallSource(mpcSource)
	//des
	evmContractCallDes := coboWaas2.NewEvmContractCallDestination(coboWaas2.CONTRACTCALLDESTINATIONTYPE_EVM_CONTRACT, to, callData)
	if value != "" && value != "0" {
		evmContractCallDes.SetValue(value)
	}
	des := coboWaas2.EvmContractCallDestinationAsContractCallDestination(evmContractCallDes)
	//param
	contractCallParams := *coboWaas2.NewContractCallParams(
		m.createRequestId(),
		chainId,
		source,
		des,
	)
	if fee != nil {
		contractCallParams.Fee = fee
	}

	if m.debug {
		debugJson, _ := json.Marshal(contractCallParams)
		m.logger.Infof("ContractCall parameters: %s", debugJson)
	}

	resp, r, err := m.client.TransactionsAPI.CreateContractCallTransaction(m.getCtx()).ContractCallParams(contractCallParams).Execute()
	return m.formatResponse(resp, r, err)
}

func (m *EvmMpcV2) Sign191Message(chainID, address, message string) (*coboWaas2.CreateTransferTransaction201Response, error) {
	mpcSource := coboWaas2.NewMpcMessageSignSource(coboWaas2.MESSAGESIGNSOURCETYPE_ORG_CONTROLLED, m.walletId, address)
	signSource := coboWaas2.MpcMessageSignSourceAsMessageSignSource(mpcSource)

	signDestination := coboWaas2.NewEvmEIP191MessageSignDestination(coboWaas2.MESSAGESIGNDESTINATIONTYPE__191_SIGNATURE, message)
	destination := coboWaas2.EvmEIP191MessageSignDestinationAsMessageSignDestination(signDestination)

	messageSignParams := *coboWaas2.NewMessageSignParams(
		m.createRequestId(),
		chainID,
		signSource,
		destination,
	)
	resp, r, err := m.client.TransactionsAPI.CreateMessageSignTransaction(m.getCtx()).MessageSignParams(messageSignParams).Execute()
	return m.formatResponse(resp, r, err)
}

func (m *EvmMpcV2) Sign712Message(chainID, address string, structuredData map[string]interface{}) (*coboWaas2.CreateTransferTransaction201Response, error) {
	mpcSource := coboWaas2.NewMpcMessageSignSource(coboWaas2.MESSAGESIGNSOURCETYPE_ORG_CONTROLLED, m.walletId, address)
	signSource := coboWaas2.MpcMessageSignSourceAsMessageSignSource(mpcSource)

	signDestination := coboWaas2.NewEvmEIP712MessageSignDestination(coboWaas2.MESSAGESIGNDESTINATIONTYPE__712_SIGNATURE, structuredData)
	destination := coboWaas2.EvmEIP712MessageSignDestinationAsMessageSignDestination(signDestination)

	messageSignParams := *coboWaas2.NewMessageSignParams(
		m.createRequestId(),
		chainID,
		signSource,
		destination,
	)

	resp, r, err := m.client.TransactionsAPI.CreateMessageSignTransaction(m.getCtx()).MessageSignParams(messageSignParams).Execute()
	return m.formatResponse(resp, r, err)
}
