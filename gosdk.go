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
	"crypto/tls"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"strings"

	"github.com/dappledger/ann-go-sdk/common"
	"github.com/dappledger/ann-go-sdk/types"
)

const GasLimit = 1000000000

type CyrptoType string

const (
	ZaCryptoType CyrptoType = "ZA"
)

type GoSDK struct {
	rpcAddr    string
	cryptoType CyrptoType
	client     *http.Client
}

func (gs *GoSDK) Url() string {
	return gs.rpcAddr
}

//Deprecated, use NewSDk instead
func New(rpcAddr string, cryptoType CyrptoType) *GoSDK {
	return &GoSDK{
		rpcAddr:    rpcAddr,
		cryptoType: cryptoType,
		client:     nil,
	}
}

func NewSDk(rpcAddr string, cryptoType CyrptoType) (*GoSDK, error) {
	var client *http.Client
	parts := strings.SplitN(rpcAddr, "://", 2)
	if len(parts) == 1 {
		//tcp
	} else if len(parts) != 2 {
		return nil, fmt.Errorf("unsupported rpc %v", rpcAddr)
	} else {
		switch parts[0] {
		case "tcp":
		case "http":
			_, err := url.Parse(rpcAddr)
			if err != nil {
				return nil, err
			}
			client = http.DefaultClient
		case "https":
			_, err := url.Parse(rpcAddr)
			if err != nil {
				return nil, err
			}
			client = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		}
	}
	return &GoSDK{
		rpcAddr:    rpcAddr,
		cryptoType: cryptoType,
		client:     client,
	}, nil
}

func (gs *GoSDK) JsonRPCCall(method string, params []byte, result interface{}) error {
	return gs.sendTxCall(method, params, result)
}

//--------------------------------Account---------------------
func (gs *GoSDK) AccountCreate() (Account, error) {
	return gs.accountCreate()
}

//--------------------------------Tx-------------------------

func (gs *GoSDK) Nonce(addr string) (uint64, error) {
	return gs.getNonce(addr)
}

func (gs *GoSDK) PendingNonce(addr string) (uint64, error) {
	return gs.getPendingNonce(addr)
}

func (gs *GoSDK) CheckHealth() (bool, error) {
	return gs.checkHealth()
}

func (gs *GoSDK) Receipt(hash string) (*types.ReceiptForDisplay, error) {
	return gs.receipt(hash)
}

func (gs *GoSDK) Balance(addr string) (*big.Int, error) {
	return gs.balance(addr)
}

func (gs *GoSDK) GetTransactionsHashByHeight(height uint64) ([]string, int, error) {
	return gs.getTransactionsHashByHeight(height)
}

func (gs *GoSDK) GetTransactionByHash(hash string) (*RPCTransaction, error) {
	_, ethtx, err := gs.getTxByHash(common.FromHex(hash))
	return ethtx, err
}

//--------------------------------EVM-------------------------

func (gs *GoSDK) ContractCreate(contract *ContractCreate) (map[string]interface{}, error) {
	return gs.contractCreate(contract)
}

func (gs *GoSDK) ContractCall(contractMethod *ContractMethod) (string, error) {
	return gs.contractCall(contractMethod, broadcast_tx_commit)
}

func (gs *GoSDK) ContractAsync(contractMethod *ContractMethod) (string, error) {
	return gs.contractCall(contractMethod, broadcast_tx_async)
}

func (gs *GoSDK) ContractRead(contractMethod *ContractMethod) (interface{}, error) {
	return gs.contractRead(contractMethod, 0)
}

func (gs *GoSDK) ContractReadByHeight(contractMethod *ContractMethod, height uint64) (interface{}, error) {
	return gs.contractRead(contractMethod, height)
}

//----------------------payload---------------------------------------
func (gs *GoSDK) Transaction(sendTx *Tx) (string, error) {
	return gs.sendTx(sendTx, "broadcast_tx_commit")
}

func (gs *GoSDK) TransactionAsync(sendTx *Tx) (string, error) {
	return gs.sendTx(sendTx, "broadcast_tx_async")
}

func (gs *GoSDK) TransactionPayLoad(txhash string) (string, error) {
	return gs.txPayLoad(txhash)
}

//----------------------KV---------------------------------------
func (gs *GoSDK) Put(kvTx *KVTx) (string, error) {
	return gs.kvPut(kvTx, "broadcast_tx_commit")
}

func (gs *GoSDK) PutAsync(kvTx *KVTx) (string, error) {
	return gs.kvPut(kvTx, "broadcast_tx_async")
}

func (gs *GoSDK) PutSignature(txSigned string) (string, error) {
	return gs.kvTxSigned(txSigned, "broadcast_tx_commit")
}

func (gs *GoSDK) PutSignatureAsync(txSigned string) (string, error) {
	return gs.kvTxSigned(txSigned, "broadcast_tx_async")
}

func (gs *GoSDK) Get(key []byte) ([]byte, error) {
	return gs.kvGet(key)
}

func (gs *GoSDK) GetWithPrefix(prefix, lastKey []byte, limit uint32) ([]*KVResult, error) {
	return gs.kvGetWithPrefix(prefix, lastKey, limit)
}

func (gs *GoSDK)GetKeyValueUpdateHistory(key []byte,pageNo uint32 ,pageSize uint32) (*ValueHistoryResult, error) {
	return gs.getKeyValueUpdateHistory(key,pageNo,pageSize)
}

//---------------------------txSigned-------------------------------------------------
func (gs *GoSDK) TranscationSignature(txSigned string) (string, error) {
	return gs.txSigned(txSigned, false)
}

func (gs *GoSDK) TranscationSignatureAsync(txSigned string) (string, error) {
	return gs.txSigned(txSigned, true)
}

//-------------------------------------node----------------------------------------------------\
func (gs *GoSDK) MakeNodeOpMsg(ndpub string, power int64, acc AccountBase, op types.ValidatorCmd) ([]byte, error) {
	return gs.makeNodeOpMsg(ndpub, power, acc, op)
}

func (gs *GoSDK) MakeNodeContractRequest(opmsg []byte, selfSign []byte, casigns []types.SigInfo, acc AccountBase) (*ContractMethod, error) {
	return gs.makeNodeContractRequest(opmsg, selfSign, casigns, acc)
}
