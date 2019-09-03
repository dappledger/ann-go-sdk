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
	"math/big"
	"strings"

	"github.com/dappledger/ann-go-sdk/common"
)

type AccountBase struct {
	PrivKey string `json:"privkey"`
	Nonce   uint64 `json:"nonce"`
}

func (ab *AccountBase) ParseHex() {
	if strings.Index(ab.PrivKey, "0x") == 0 {
		ab.PrivKey = ab.PrivKey[2:]
	}
}

type ContractCreate struct {
	AccountBase
	Code   string        `json:"code"`
	ABI    string        `json:"abiDefinition"`
	Params []interface{} `json:"params"`
}

type ContractMethod struct {
	AccountBase
	Contract string        `json:"contract"`
	ABI      string        `json:"abiDefinition"`
	Method   string        `json:"method"`
	Params   []interface{} `json:"params"`
}

func (cm *ContractMethod) ParseHex() {
	if strings.Index(cm.Contract, "0x") == 0 {
		cm.Contract = cm.Contract[2:]
	}
	cm.AccountBase.ParseHex()
}

type QueryResult struct {
	Value string `json:"value"`
	Msg   string `json:"msg"`
}

// RPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
type RPCTransaction struct {
	BlockHash        []byte          `json:"blockHash"`
	BlockHeight      uint64          `json:"blockHeight"`
	From             common.Address  `json:"from"`
	Gas              uint64          `json:"gas"`
	GasPrice         *big.Int        `json:"gasPrice"`
	Hash             common.Hash     `json:"hash"`
	Input            []byte          `json:"input"`
	Nonce            uint64          `json:"nonce"`
	To               *common.Address `json:"to"`
	TransactionIndex uint64          `json:"transactionIndex"`
	Value            *big.Int        `json:"value"`
	V                *big.Int        `json:"v"`
	R                *big.Int        `json:"r"`
	S                *big.Int        `json:"s"`
}
