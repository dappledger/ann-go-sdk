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
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/dappledger/ann-go-sdk/common"
	"github.com/dappledger/ann-go-sdk/rlp"
	"github.com/dappledger/ann-go-sdk/types"
)

const max_payload_size = 4000

type Tx struct {
	AccountBase
	To      string   `json:"to"` //to address;
	Payload string   `json:"payload"`
	Value   *big.Int `json:"value"`
}

func (gs *GoSDK) sendTx(sendTx *Tx, funcType string) (hash string, err error) {
	if sendTx.PrivKey == "" {
		return "", fmt.Errorf("account privkey is empty.")
	}

	if sendTx.To == "" {
		return "", fmt.Errorf("to address is empty.")
	}

	if strings.Index(sendTx.PrivKey, "0x") == 0 {
		sendTx.PrivKey = sendTx.PrivKey[2:]
	}
	if strings.Index(sendTx.To, "0x") == 0 {
		sendTx.To = sendTx.To[2:]
	}

	privBytes := common.Hex2Bytes(sendTx.PrivKey)

	nonce := sendTx.Nonce
	if nonce == 0 {
		addrBytes, err := gs.getAddrBytes(privBytes)
		if err != nil {
			return "", err
		}
		nonce, err = gs.getNonce(common.Bytes2Hex(addrBytes))
		if err != nil {
			return "", err
		}
	}

	to := common.HexToAddress(sendTx.To)
	value := sendTx.Value

	payload := sendTx.Payload
	if len(payload) > max_payload_size {
		err = fmt.Errorf("payload length must be less than 4000")
		return "", err
	}
	data := []byte(payload)

	tx := types.NewTransaction(nonce, to, value, gs.GasLimit(), big.NewInt(0), data)

	signer, sig, err := gs.signTx(privBytes, tx)
	if err != nil {
		return "", err
	}
	sigTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return "", err
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return "", err
	}
	rpcResult := new(types.ResultBroadcastTxCommit)
	err = gs.sendTxCall(funcType, txBytes, rpcResult)
	if err != nil {
		return "", err
	}
	return rpcResult.TxHash, nil
}

func (gs *GoSDK) txPayLoad(txstr string) (string, error) {
	if strings.Index(txstr, "0x") == 0 {
		txstr = txstr[2:]
	}

	txhash, err := hex.DecodeString(txstr)
	if err != nil {
		return "", err
	}
	query := make([]byte, len(txhash)+1)
	query[0] = QueryTxExecution
	copy(query[1:], txhash)

	query = append([]byte{types.QueryType_PayLoad}, query...)

	res := new(types.ResultQuery)
	err = gs.sendTxCall("query", query, res)
	if err != nil {
		return "", err
	}

	return string(res.Result.Data), nil
}
