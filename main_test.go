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
	"fmt"
	"testing"
	"time"
)

var abis string = `[
	{
		"constant": true,
		"inputs": [],
		"name": "GetResult",
		"outputs": [
			{
				"name": "",
				"type": "int256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "a",
				"type": "int256"
			}
		],
		"name": "Add",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "a",
				"type": "int256"
			},
			{
				"indexed": false,
				"name": "result",
				"type": "int256"
			}
		],
		"name": "onAdd",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "",
				"type": "int256"
			}
		],
		"name": "onResult",
		"type": "event"
	}
]`

var bytcode string = "608060405234801561001057600080fd5b50610140806100206000396000f30060806040526004361061004b5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416639a7d9af18114610050578063ff108fcb14610077575b600080fd5b34801561005c57600080fd5b50610065610091565b60408051918252519081900360200190f35b34801561008357600080fd5b5061008f600435610097565b005b60005490565b6000805482019081905560408051838152602081019290925280517f2e81aa73693352ce9bdefb8e614e9e5cf9e05ea2a3fde2326285a70c70e1954e9281900390910190a160005460408051918252517f3c401639fc256a3b5bf32b9fe59c87252b2d137bc6e1097e7e14c0ea9190cb6a9181900360200190a1505600a165627a7a72305820c2965db3127c35ebb622a1f6e27cca74992cbd1e549b3fcc6ad35b924514a8fd0029"

var superAddr string = "0x65188459a1dc65984a0c7d4a397ed3986ed0c853"
var superPriv string = "7cb4880c2d4863f88134fd01a250ef6633cc5e01aeba4c862bedbf883a148ba8"

var testAddr string = "0xc409aaf73698fdb5995c4d85f6033d5e90d2f2bd"
var testPriv string = "64d18eb7061dff419581c1af98201b76c7ab6db538b1cf65123c470ccc6d5929"

var testValidatorPub string = "BDCB2DB0320FECF1D771852E49FE65D47AD82043028318D286C185418B571980"
var testSigs string = "3d7e15d350673620aaa578a79d5accaa5c1a6f4932f47a5e402a425f3893aec3f73303221652bda71ba0e95438c40705db2207d1081dd66d8c7468c283409a06"

var (
	newAddr, newPriv string
	nonce            uint64
	hash             string

	contractAddress string
	contractReceipt string
	contractHash    string
	height          uint64

	mapAddress map[string]string
	client     *AnnChainClient

	routines  int = 10
	sustained int = 100
)

func init() {
	client = NewAnnChainClient("tcp://127.0.0.1:46657")
	mapAddress = make(map[string]string, routines)
}

func TestGen(t *testing.T) {

	newPriv, newAddr = GenerateKey()

	t.Log(newPriv, newAddr, len(newAddr))

}

func GetNonce(source string) uint64 {

	nonce, _, _ := client.QueryNonce(source)

	return nonce
}

//**************************Validator Equity TEST*****************************
func TestRequestSpecialOP(t *testing.T) {

	client = NewAnnChainClient("tcp://127.0.0.1:46657")

	chash, code, err := client.RequestSpecialOP(superPriv, testValidatorPub, testSigs, "tcp://127.0.0.1:46657", superAddr, true, 1)

	hash = chash

	t.Log(hash, code, err)

	time.Sleep(time.Second)
}

//**************************Account TEST *************************************

func TestCreateAccounts(t *testing.T) {

	nonce := GetNonce(superAddr)

	for i := 0; i < 5; i++ {

		_, newAddr = GenerateKey()

		chash, code, err := client.CreateAccount(nonce, superPriv, "0", "memo", superAddr, newAddr, "1000")

		t.Log(chash, code, err)

		nonce++
	}
}

func TestCreateAccount(t *testing.T) {

	chash, code, err := client.CreateAccount(GetNonce(superAddr), superPriv, "0", "memo", superAddr, testAddr, "1000")

	hash = chash

	t.Log(hash, code, err)

	time.Sleep(time.Second)
}
func TestQueryAccount(t *testing.T) {

	result, code, err := client.QueryAccount(superAddr)

	t.Log(result, code, err)
}

//****************************************************************************

//**************************ManageData TEST***********************************

func TestManageData(t *testing.T) {

	var datas map[string]ManageDataValueParam

	datas = make(map[string]ManageDataValueParam)

	datas["8"] = ManageDataValueParam{Value: "zhaoyang", Category: "B"}
	datas["9"] = ManageDataValueParam{Value: "fanhongyue", Category: "B"}

	result, code, err := client.ManageData(GetNonce(superAddr), superPriv, "100", "memo", superAddr, datas)

	t.Log(result, code, err)

	hash = result

	time.Sleep(time.Second)
}

func TestQueryAccountManageDatas(t *testing.T) {

	result, code, err := client.QueryAccountManageDatas(superAddr, "asc", 10, 0)

	t.Log(result, code, err)
}

func TestQueryAccountManageData(t *testing.T) {

	result, code, err := client.QueryAccountManageData(superAddr, "3")

	t.Log(result, code, err)
}

func TestQueryAccountCategoryManageData(t *testing.T) {

	result, code, err := client.QueryAccountCategoryManageData(superAddr, "B")

	t.Log(result, code, err)
}

//*********************************************************************************************

//************************************Ledger TEST *********************************************

func TestQueryTransactions(t *testing.T) {

	result, code, err := client.QueryTransactions("desc", 10, 0)

	t.Log(result, code, err)
}

func TestQueryTransaction(t *testing.T) {

	result, code, err := client.QueryTransaction(hash)

	t.Log(result, code, err)

	time.Sleep(time.Second)
}

func TestQueryAccountTransactions(t *testing.T) {

	result, code, err := client.QueryAccountTransactions(superAddr, "asc", 10, 0)

	t.Log(result, code, err)
}

func TestQueryLedgerTransactions(t *testing.T) {

	result, code, err := client.QueryLedgerTransactions(23, "asc", 10, 0)

	t.Log(result, code, err)
}

//*********************************************************************************************

//************************************Contract TEST********************************************

func TestCreateContract(t *testing.T) {

	param, err := NewCreateContractParam("1", "8000000", "0", bytcode, abis, []interface{}{})

	if err != nil {
		t.Log(err)
		return
	}

	ccontractHash, code, err := client.CreateContract(GetNonce(superAddr), superPriv, "100", "", superAddr, param)

	contractHash = ccontractHash

	t.Log(contractHash, code, err)

	time.Sleep(time.Second * 4)
}

func TestQueryCreateReceipt(t *testing.T) {

	result, code, err := client.QueryReceipt(contractHash)

	contractAddress = result.ContractAddress

	t.Log(result, code, err)
}

func TestQueryContractExist(t *testing.T) {

	result, code, err := client.QueryContractExist(contractAddress)

	t.Log(result, code, err)
}

func TestExecuteContract(t *testing.T) {

	param, err := NewExecuteContractParam("1", "8000000", "0", "Add", abis, []interface{}{1})

	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(superAddr), superPriv, "0", "", superAddr, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)
}

func OnMyEvent(inputs []interface{}) {
	fmt.Println(inputs)
}

func TestListenEvent(t *testing.T) {

	err := client.ListenEvent(contractHash, abis, OnMyEvent, 10)
	if err != nil {
		t.Log(err)
	}

}

func TestQueryContract(t *testing.T) {

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "GetResult", abis, []interface{}{})

	t.Log(result, code, err)
}

//*************************************************************************************************

//*************************************Payment TEST************************************************

func TestPayment(t *testing.T) {

	result, code, err := client.Payment(GetNonce(superAddr), superPriv, "0", "memo", superAddr, testAddr, "100")

	t.Log(result, code, err)

	hash = result

	time.Sleep(time.Second)
}

func TestQueryPayments(t *testing.T) {

	result, code, err := client.QueryPayments("asc", 100, 0)

	t.Log(result, code, err)
}

func TestQueryAccountPayments(t *testing.T) {

	result, code, err := client.QueryAccountPayments(superAddr, "asc", 100, 0)

	t.Log(result, code, err)
}

func TestQueryPayment(t *testing.T) {

	result, code, err := client.QueryPayment(hash)

	t.Log(result, code, err)
}

//*************************************************************************************************
