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
	"encoding/json"
	"fmt"
	"time"

	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/types"
)

func (gs *GoSDK) checkHealth() (bool, error) {
	rpcResult := new(types.ResultHealthInfo)
	err := gs.sendTxCall("healthinfo", nil, rpcResult)
	if err != nil {
		return false, err
	}
	//same as rendez-api
	return (200 == rpcResult.Status), nil
}

func (gs *GoSDK) Validators() (*types.ResultValidators, error) {
	rpcResult := new(types.ResultValidators)
	err := gs.JsonRPCCall("validators", []byte{}, rpcResult)
	if err != nil {
		return nil, err
	}

	return rpcResult, nil
}

const (
	AdminTo     = "0x0000000000000000000000000000000002000000" //contract addr;
	AdminMethod = "changenode"
	AdminABI    = `[
		{
			"constant": false,
			"inputs": [
				{
					"name": "txdata",
					"type": "bytes"
				}
			],
			"name": "changenode",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`
)

func (gs *GoSDK) makeNodeOpMsg(ndpub string, power int64, acc AccountBase, op types.ValidatorCmd) ([]byte, error) {
	pubdata := common.FromHex(ndpub)
	addrBytes, err := gs.getAddrBytes(common.FromHex(acc.PrivKey))
	if err != nil {
		return nil, err
	}
	vAttr := types.ValidatorAttr{
		PubKey: pubdata,
		Cmd:    op,
		Addr:   addrBytes,
		Nonce:  acc.Nonce,
	}
	if op == types.ValidatorCmdUpdateNode && power >= 0 {
		vAttr.Power = power
	}
	return json.Marshal(&vAttr)
}

func (gs *GoSDK) makeNodeContractRequest(opmsg []byte, selfSign []byte, casigns []types.SigInfo, acc AccountBase) (*ContractMethod, error) {
	opcmd := &types.AdminOPCmd{
		"changeValidator",
		opmsg,
		time.Now(),
		selfSign,
		casigns,
	}
	d, err := json.Marshal(opcmd)
	if err != nil {
		return nil, err
	}
	d = types.TagAdminOPTx(d)
	return &ContractMethod{
		AccountBase: acc,
		Contract:    AdminTo,
		ABI:         AdminABI,
		Method:      AdminMethod,
		Params:      []interface{}{fmt.Sprintf("0x%x", d)},
	}, nil
}
