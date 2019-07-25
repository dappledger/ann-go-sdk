package sdk

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/rlp"
	"github.com/dappledger/AnnChain-go-sdk/rpc"
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

func (gs *GoSDK) receipt(strHash string) (*types.ReceiptForStorage, error) {

	if strings.Index(strHash, "0x") == 0 {
		strHash = strHash[2:]
	}

	bytHash := common.Hex2Bytes(strHash)

	queryParam := append([]byte{types.QueryType_Receipt}, bytHash...)

	resultQuery := new(types.ResultQuery)

	err := gs.sendTxCall("query", queryParam, resultQuery)
	if err != nil {
		return nil, err
	}

	if 0 != resultQuery.Result.Code {
		return nil, errors.New(resultQuery.Result.Log)
	}

	receipt := new(types.ReceiptForStorage)

	err = rlp.DecodeBytes(resultQuery.Result.Data, receipt)
	if err != nil {
		return nil, err
	}

	resultTrans, tx, err := gs.getTxByHash(bytHash)
	if err != nil {
		return nil, err
	}

	from, err := types.Sender(new(types.HomesteadSigner), tx)
	if err != nil {
		return nil, err
	}

	receipt.TxIndex = resultTrans.TransactionIndex
	receipt.Height = resultTrans.BlockHeight
	receipt.BlockHashHex = common.ToHex(resultTrans.BlockHash)
	receipt.From = from
	receipt.Timestamp = time.Unix(int64(resultTrans.Timestamp/uint64(time.Second)), int64(resultTrans.Timestamp%uint64(time.Second)))
	if tx.To() == nil {
		receipt.To = common.Address{}
	} else {
		receipt.To = *tx.To()
	}

	return receipt, nil
}

func (gs *GoSDK) getTxByHash(hash []byte) (*ResultTransaction, *types.Transaction, error) {

	res := new(types.ResultQuery)

	err := gs.sendTxCall("transaction", hash, res)

	if err != nil {
		return nil, nil, err
	}

	if 0 != res.Result.Code {
		return nil, nil, errors.New(res.Result.Log)
	}

	resultTrans := &ResultTransaction{}

	err = rlp.DecodeBytes(res.Result.Data, &resultTrans)
	if err != nil {
		return nil, nil, err
	}

	tx := &types.Transaction{}

	err = rlp.DecodeBytes(resultTrans.RawTransaction, tx)
	if err != nil {
		return nil, nil, err
	}
	return resultTrans, tx, nil
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

func (gs *GoSDK) getTransactionsHashByHeight(height uint64) (hashs []string, total int, err error) {
	res := new(ResultBlock)
	clientJSON := rpc.NewClientJSONRPC(gs.rpcAddr)
	var _params []interface{}
	_params = []interface{}{height}
	_, err = clientJSON.Call("block", _params, res)
	if err != nil {
		return nil, 0, err
	}

	total = int(res.Block.Header.NumTxs)
	blockdata := res.Block.Data
	if len(blockdata.Txs)+len(blockdata.ExTxs) != total {
		err = fmt.Errorf("logic err:NumTxs<%d> != len(txs)<%d>", total, len(blockdata.Txs)+len(blockdata.ExTxs))
		return nil, 0, err
	}
	hashs = make([]string, total)
	for idx, tx := range blockdata.Txs {
		hashs[idx] = fmt.Sprintf("%x", tx.Hash())
	}
	base := len(blockdata.Txs)
	for idx, tx := range blockdata.ExTxs {
		hashs[idx+base] = fmt.Sprintf("%x", tx.Hash())
	}
	return hashs, total, nil
}
