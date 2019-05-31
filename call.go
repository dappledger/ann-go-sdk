package sdk

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/rlp"
	"github.com/dappledger/AnnChain-go-sdk/types"
)

const (
	// angine takes query id from 0x01 to 0x2F
	QueryTxExecution = 0x01
)

func (gs *GoSDK) getNonce(addr string) (uint64, error) {
	if !common.IsHexAddress(addr) {
		return 0, fmt.Errorf("Invalid address(is not hex) %s", addr)
	}
	if strings.Index(addr, "0x") == 0 {
		addr = addr[2:]
	}

	address := common.Hex2Bytes(addr)
	query := append([]byte{types.QueryType_Nonce}, address...)
	rpcResult := new(types.ResultQuery)
	err := gs.sendTxCall("query", query, rpcResult)
	if err != nil {
		return 0, err
	}
	nonce := new(uint64)
	rlp.DecodeBytes(rpcResult.Result.Data, nonce)
	return *nonce, nil
}

func (gs *GoSDK) block(hashstr string) ([]string, int, error) {

	arryTxs := make([]string, 0)

	if strings.Index(hashstr, "0x") == 0 {
		hashstr = hashstr[2:]
	}

	hash := common.Hex2Bytes(hashstr)
	query := append([]byte{types.QueryType_BlockHash}, hash...)
	res := new(types.ResultQuery)
	err := gs.sendTxCall("query", query, res)
	if err != nil {
		return nil, 0, err
	}
	if 0 != res.Result.Code {
		return nil, 0, fmt.Errorf(string(res.Result.Log))
	}
	common.Bytes2Hex(res.Result.Data)
	var blockHashs common.Hashs
	err = rlp.DecodeBytes(res.Result.Data, &blockHashs)
	if err != nil {
		return nil, 0, err
	}
	for _, txhash := range blockHashs {
		arryTxs = append(arryTxs, txhash.Hex())
	}

	return arryTxs, len(arryTxs), nil
}

func (gs *GoSDK) receipt(hashstr string) (*types.ReceiptForStorage, error) {
	if strings.Index(hashstr, "0x") == 0 {
		hashstr = hashstr[2:]
	}

	hash := common.Hex2Bytes(hashstr)
	query := append([]byte{types.QueryType_Receipt}, hash...)
	res := new(types.ResultQuery)
	err := gs.sendTxCall("query", query, res)
	if err != nil {
		return nil, err
	}
	if 0 != res.Result.Code {
		return nil, fmt.Errorf(string(res.Result.Log))
	}
	common.Bytes2Hex(res.Result.Data)
	receiptForStorage := new(types.ReceiptForStorage)
	err = rlp.DecodeBytes(res.Result.Data, receiptForStorage)
	if err != nil {
		return nil, err
	}
	return receiptForStorage, nil
}

func (gs *GoSDK) balance(addr string) (*big.Int, error) {
	if !common.IsHexAddress(addr) {
		return big.NewInt(0), fmt.Errorf("Invalid address(is not hex) %s", addr)
	}
	if strings.Index(addr, "0x") == 0 {
		addr = addr[2:]
	}

	address := common.Hex2Bytes(addr)
	query := append([]byte{types.QueryType_Balance}, address...)
	rpcResult := new(types.ResultQuery)
	err := gs.sendTxCall("query", query, rpcResult)
	if err != nil {
		return big.NewInt(0), err
	}
	balance := big.NewInt(0)
	rlp.DecodeBytes(rpcResult.Result.Data, balance)
	return balance, nil
}

func (gs *GoSDK) txSigned(tx string, isAsync bool) (string, error) {
	if strings.HasPrefix(tx, "0x") {
		tx = tx[2:]
	}
	txBys := common.Hex2Bytes(tx)
	rpcResult := new(types.ResultBroadcastTxCommit)
	method := "broadcast_tx_commit"
	if isAsync {
		method = "broadcast_tx_async"
	}
	err := gs.sendTxCall(method, txBys, rpcResult)
	if err != nil {
		return "", err
	}
	hash := rpcResult.TxHash
	return hash, nil
}
