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

package smoke

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dappledger/ann-go-sdk"
	"github.com/dappledger/ann-go-sdk/common"
	"github.com/dappledger/ann-go-sdk/crypto"
	"github.com/dappledger/ann-go-sdk/types"
)

const (
	accPriv  = "48deaa73f328f38d5fcb29d076b2b639c8491f97d245fc22e95a86366687903a"
	accAddr  = "28112ca022224ae7757bcd559666be5340ff109a"
	byteCode = `608060405234801561001057600080fd5b5061020c806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806326f34b8b146100725780633c6bb436146100b35780633d4197f0146100de578063862c242e1461010b578063e1cb0e5214610142575b600080fd5b34801561007e57600080fd5b5061009d6004803603810190808035906020019092919050505061016d565b6040518082815260200191505060405180910390f35b3480156100bf57600080fd5b506100c861018d565b6040518082815260200191505060405180910390f35b3480156100ea57600080fd5b5061010960048036038101908080359060200190929190505050610193565b005b34801561011757600080fd5b50610140600480360381019080803590602001909291908035906020019092919050505061019d565b005b34801561014e57600080fd5b506101576101d7565b6040518082815260200191505060405180910390f35b600060016000838152602001908152602001600020600101549050919050565b60005481565b8060008190555050565b8160016000848152602001908152602001600020600001819055508060016000848152602001908152602001600020600101819055505050565b600080549050905600a165627a7a7230582090454576ac53e48f93db33e3502a6c1c9bff38697b8035a76500dc8ab84056b50029`
	abi      = `[{"constant": true,"inputs": [{"name": "_no","type": "uint256"}],"name": "getBatchVal","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": true,"inputs": [],"name": "val","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_val","type": "uint256"}],"name": "setVal","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": false,"inputs": [{"name": "_no","type": "uint256"},{"name": "_val","type": "uint256"}],"name": "setBatchVal","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": true,"inputs": [],"name": "getVal","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"}]`
)

func ExpectHexEqual(t *testing.T, hex1, hex2 string) {
	if strings.HasPrefix(hex1, "0x") || strings.HasPrefix(hex1, "0X") {
		hex1 = hex1[2:]
	}
	hex1 = strings.ToUpper(hex1)
	if strings.HasPrefix(hex2, "0x") || strings.HasPrefix(hex2, "0X") {
		hex2 = hex2[2:]
	}
	hex2 = strings.ToUpper(hex2)
	assert.Equal(t, hex1, hex2)
}

var (
	cas = []string{
		"5CECE8180C49B637ED8447F40063D962CB78A2854EDE3EFB6E50B9D2BBC77D4B62D09588461E764E3B14D332FBBA8C7EC1EDA1AE4BE056DCC1392246BF31D522",
		"F0B19CD1A93F4239B238C19D66C5A16883976949E7D002A16948EFC6E8F7E38F34306B2DD3A43E90DD3B32A3FD1FD8AA45811F9E06A3358E2E15506C7B3A8A56",
		"8B90818748BE5E11E10929EF193E5129A6389DC0410718F59FC51DEACE6BBC59C0CB6B46F66F731B9F6A6C202D0DD618AA9034F06C198FE31FDE9F3DE2188479",
	}
	opnode_pub = "E866267CFD46C8C1BC3649E4D15302861FFD44F29C0524F509B6877DFB6A15EB"
	opnode_prv = "883708FDCFDA9DA9BD5923E03647CE8C48AA25D8D926BE1AE061FC919930CD11E866267CFD46C8C1BC3649E4D15302861FFD44F29C0524F509B6877DFB6A15EB"
)

func NodeSign(priv string, data []byte) types.SigInfo {
	privK := crypto.SetNodePrivKey("ZA", common.FromHex(priv))
	pub := privK.PubKey()
	s := privK.Sign(data)
	return types.SigInfo{
		common.FromHex(pub.KeyString()),
		common.FromHex(s.KeyString()),
	}
}

func TestNode(t *testing.T) {
	client := sdk.New("localhost:46657", sdk.ZaCryptoType)
	acc, err := client.AccountCreate()
	assert.Nil(t, err)
	nonce, err := client.Nonce(acc.Address)
	assert.Nil(t, err)
	accbase := sdk.AccountBase{
		acc.Privkey,
		nonce,
	}
	opcmds := []types.ValidatorCmd{types.ValidatorCmdRemoveNode, types.ValidatorCmdAddPeer, types.ValidatorCmdUpdateNode}
	powers := []int64{0, 0, 100}
	//remove node;
	for idx, opcmd := range opcmds {
		power := powers[idx]
		//
		data, err := client.MakeNodeOpMsg(opnode_pub, power, accbase, opcmd)
		assert.Nil(t, err)
		sinfo := NodeSign(opnode_prv, data)
		var caSinfo []types.SigInfo
		for _, pk := range cas {
			caSinfo = append(caSinfo, NodeSign(pk, data))
		}
		req, err := client.MakeNodeContractRequest(data, sinfo.Signature, caSinfo, accbase)
		assert.Nil(t, err)
		//
		_, err = client.ContractCall(req)
		assert.Nil(t, err)
		//
		vals, err := client.Validators()
		assert.Nil(t, err)
		d, err := json.MarshalIndent(vals, "", "\t")
		assert.Nil(t, err)
		//
		time.Sleep(time.Second * 4)
		fmt.Printf("%s:\n%s\n\n", opcmd, string(d))
		accbase.Nonce++
	}
}

func TestZA(t *testing.T) {
	client := sdk.New("localhost:46657", sdk.ZaCryptoType)
	//health
	isHealth, err := client.CheckHealth()
	assert.Nil(t, err)
	assert.Equal(t, true, isHealth)
	//nonce
	nonce0, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	// evm contract
	var arg = sdk.ContractCreate{
		AccountBase: sdk.AccountBase{
			PrivKey: accPriv,
			Nonce:   nonce0,
		},
		ABI:  abi,
		Code: byteCode,
	}
	result, err := client.ContractCreate(&arg)
	assert.Nil(t, err)
	contractId := result["contract"].(string)

	time.Sleep(2 * time.Second)
	nonce1, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	assert.Equal(t, nonce0+1, nonce1)
	params := []interface{}{168}
	var callArg = sdk.ContractMethod{
		AccountBase: sdk.AccountBase{
			PrivKey: accPriv,
			Nonce:   nonce1,
		},
		ABI:      abi,
		Contract: contractId,
		Method:   "setVal",
		Params:   params,
	}
	txHash, err := client.ContractCall(&callArg)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
	//receipt check;
	receipt, err := client.Receipt(txHash)
	assert.Nil(t, err)
	hashs, _, err := client.GetTransactionsHashByHeight(receipt.Height)
	assert.Nil(t, err)
	ExpectHexEqual(t, receipt.TxHash.Hex(), hashs[receipt.TransactionIndex])
	ExpectHexEqual(t, receipt.From.Hex(), accAddr)
	ExpectHexEqual(t, receipt.To.Hex(), contractId)

	assert.Equal(t, txHash, receipt.TxHash.Hex())
	assert.Nil(t, err)
	nonce2, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	assert.Equal(t, nonce1+1, nonce2)
	callArg.Method = "getVal"
	callArg.Params = nil
	callArg.Nonce = nonce2
	resp, err := client.ContractRead(&callArg)
	assert.Nil(t, err)
	res := resp.([]interface{})
	assert.Equal(t, big.NewInt(168), res[0].(*big.Int))

	nonce3, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	assert.Equal(t, nonce3, nonce2)
	params = []interface{}{169}
	callArg.Method = "setVal"
	callArg.Params = params
	txHash, err = client.ContractCall(&callArg)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
	_, err = client.Receipt(txHash)
	assert.Nil(t, err)
	nonce4, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	assert.Equal(t, nonce3+1, nonce4)
	callArg.Method = "getVal"
	callArg.Params = nil
	callArg.Nonce = nonce4
	resp, err = client.ContractRead(&callArg)
	assert.Nil(t, err)
	res = resp.([]interface{})
	assert.Equal(t, big.NewInt(169), res[0].(*big.Int))
}
