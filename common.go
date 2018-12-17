// Copyright 2017 ZhongAn Information Technology Services Co.,Ltd.
//
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

package annchain

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/util/abi"
	at "github.com/dappledger/AnnChain/angine/types"
	ethcmn "github.com/dappledger/AnnChain/genesis/eth/common"
	"github.com/dappledger/AnnChain/genesis/eth/rlp"
	"github.com/dappledger/AnnChain/genesis/types"
)

func (c *AnnChainClient) signAndEncodeTx(tx *types.Transaction, privkey *ecdsa.PrivateKey) ([]byte, string, at.CodeType, error) {

	var (
		sigTx *types.Transaction
		err   error
	)

	sigTx = tx

	if privkey != nil {
		sigTx, err = tx.Sign(privkey)
		if err != nil {
			return nil, "", at.CodeType_InvalidTx, err
		}
	}
	txBytes, err := rlp.EncodeToBytes(sigTx.Data)
	if err != nil {
		return nil, "", at.CodeType_WrongRLP, err
	}

	return txBytes, sigTx.Hash().Hex(), at.CodeType_OK, nil
}

func packCreateContractData(bytecode, abis string, params []interface{}) (string, error) {

	jAbi, err := abi.JSON(strings.NewReader(abis))
	if err != nil {
		return "", err
	}

	byteCode := ethcmn.Hex2Bytes(bytecode)
	if len(byteCode) == 0 {
		return "", fmt.Errorf("bytecode is null")
	}

	if len(params) > 0 {
		args, err := abi.ParseArgs("", jAbi, params)
		if err != nil {
			return "", err
		}
		packData, err := jAbi.Pack("", args...)
		if err != nil {
			return "", err
		}
		byteCode = append(byteCode, packData...)
	}

	return ethcmn.Bytes2Hex(byteCode), nil

}

func packExecuteContractData(abis, funcname string, params []interface{}) (string, error) {

	jAbi, err := abi.JSON(strings.NewReader(abis))
	if err != nil {
		return "", err
	}

	args, err := abi.ParseArgs(funcname, jAbi, params)
	if err != nil {
		return "", err
	}

	packData, err := jAbi.Pack(funcname, args...)
	if err != nil {
		return "", err
	}

	return ethcmn.Bytes2Hex(packData), nil
}

func unpackResultToArray(funcname, abis string, output []byte) (interface{}, error) {

	abiDef, err := abi.JSON(strings.NewReader(abis))
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return nil, nil
	}
	m, ok := abiDef.Methods[funcname]
	if !ok {
		return nil, errors.New("No such method")
	}
	if len(m.Outputs) == 0 {
		return nil, errors.New("method " + m.Name + " doesn't have any returns")
	}
	if len(m.Outputs) == 1 {
		var result interface{}
		d := ethcmn.ParseData(output)
		if err := abiDef.Unpack(&result, funcname, d); err != nil {
			return nil, err
		}
		return result, nil
	}
	var result []interface{}
	d := ethcmn.ParseData(output)
	if err := abiDef.Unpack(&result, funcname, d); err != nil {
		return nil, err
	}
	return result, nil
}
