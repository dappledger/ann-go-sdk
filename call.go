// Copyright © 2017 ZhongAn Technology
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sdk

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/dappledger/ann-go-sdk/common"
	"github.com/dappledger/ann-go-sdk/rlp"
	"github.com/dappledger/ann-go-sdk/types"
)

type ResultTransaction struct {
	BlockHash        []byte `json:"block_hash"`
	BlockHeight      uint64 `json:"block_height"`
	TransactionIndex uint64 `json:"transaction_index"`
	RawTransaction   []byte `json:"raw_transaction"`
	Timestamp        uint64 `json:"timestamp"`
}

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

func (gs *GoSDK) receipt(hashstr string) (*types.ReceiptForDisplay, error) {
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
	rt, etx, err := gs.getTxByHash(hash)
	if err != nil {
		return nil, err
	}

	timestamp := time.Unix(int64(rt.Timestamp/uint64(time.Second)), int64(rt.Timestamp%uint64(time.Second)))
	receiptForShow := &types.ReceiptForDisplay{
		PostState:         receiptForStorage.PostState,
		Status:            receiptForStorage.Status,
		CumulativeGasUsed: receiptForStorage.CumulativeGasUsed,
		Bloom:             receiptForStorage.Bloom,
		Logs:              receiptForStorage.Logs,
		TxHash:            receiptForStorage.TxHash,
		ContractAddress:   receiptForStorage.ContractAddress,
		GasUsed:           receiptForStorage.GasUsed,
		TransactionIndex:  rt.TransactionIndex,
		BlockHashHex:      fmt.Sprintf("0x%x", rt.BlockHash),
		Height:            rt.BlockHeight,
		From:              etx.From,
		To:                etx.To,
		Timestamp:         timestamp,
	}

	return receiptForShow, nil
}

func (gs *GoSDK) getTxByHash(hash []byte) (*ResultTransaction, *RPCTransaction, error) {
	res := new(types.ResultQuery)
	err := gs.sendTxCall("transaction", hash, res)
	if err != nil {
		return nil, nil, err
	}
	var rt = &ResultTransaction{}
	data := res.Result.Data
	err = rlp.DecodeBytes(data, &rt)
	if err != nil {
		return nil, nil, err
	}

	ethtx := &types.Transaction{}
	err = rlp.DecodeBytes(rt.RawTransaction, ethtx)
	if err != nil {
		return nil, nil, err
	}

	var (
		signer  types.Signer
		v, r, s *big.Int
	)

	signer = new(types.HomesteadSigner)

	from, err := types.Sender(signer, ethtx)
	if err != nil {
		return nil, nil, err
	}

	v, r, s = ethtx.RawSignatureValues()

	rpc := &RPCTransaction{
		BlockHash:        rt.BlockHash,
		BlockHeight:      rt.BlockHeight,
		From:             from,
		Hash:             common.BytesToHash(hash),
		Input:            ethtx.Data(),
		Nonce:            ethtx.Nonce(),
		To:               ethtx.To(),
		TransactionIndex: rt.TransactionIndex,
		Value:            ethtx.Value(),
		V:                v,
		R:                r,
		S:                s,
	}

	return rt, rpc, nil
}

func (gs *GoSDK) getTransactionsHashByHeight(height uint64) (hashs []string, total int, err error) {
	res := new(types.ResultBlock)
	clientJSON := gs.NewClientJsonRPC()
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
