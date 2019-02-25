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
	"encoding/hex"
	"fmt"
	"testing"
	"time"
)

var abis string = `[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [
			{
				"name": "",
				"type": "string"
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
				"name": "_spender",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "approve",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
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
				"name": "_from",
				"type": "address"
			},
			{
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transferFrom",
		"outputs": [
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [
			{
				"name": "",
				"type": "uint8"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_owner",
				"type": "address"
			}
		],
		"name": "balanceOf",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "owner",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [
			{
				"name": "",
				"type": "string"
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
				"name": "_to",
				"type": "address"
			},
			{
				"name": "_value",
				"type": "uint256"
			}
		],
		"name": "transfer",
		"outputs": [
			{
				"name": "",
				"type": "address"
			},
			{
				"name": "",
				"type": "uint256"
			},
			{
				"name": "",
				"type": "bool"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_owner",
				"type": "address"
			},
			{
				"name": "_spender",
				"type": "address"
			}
		],
		"name": "allowance",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "user",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "from",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "to",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "value",
				"type": "uint256"
			}
		],
		"name": "Transfer",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": true,
				"name": "owner",
				"type": "address"
			},
			{
				"indexed": true,
				"name": "spender",
				"type": "address"
			},
			{
				"indexed": false,
				"name": "value",
				"type": "uint256"
			}
		],
		"name": "Approval",
		"type": "event"
	}
]`

var bytcode string = "6080604052620f424060005534801561001757600080fd5b50604051602080610ec48339810180604052810190808051906020019092919050505033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060005460026000600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506107d0600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050610d8b806101396000396000f3006080604052600436106100a4576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306fdde03146100a9578063095ea7b31461013957806318160ddd1461019e57806323b872dd146101c9578063313ce5671461024e57806370a082311461027f5780638da5cb5b146102d657806395d89b411461032d578063a9059cbb146103bd578063dd62ed3e1461045c575b600080fd5b3480156100b557600080fd5b506100be6104d3565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100fe5780820151818401526020810190506100e3565b50505050905090810190601f16801561012b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561014557600080fd5b50610184600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061050c565b604051808215151515815260200191505060405180910390f35b3480156101aa57600080fd5b506101b36105fe565b6040518082815260200191505060405180910390f35b3480156101d557600080fd5b50610234600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610607565b604051808215151515815260200191505060405180910390f35b34801561025a57600080fd5b506102636109c6565b604051808260ff1660ff16815260200191505060405180910390f35b34801561028b57600080fd5b506102c0600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506109cb565b6040518082815260200191505060405180910390f35b3480156102e257600080fd5b506102eb610a14565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561033957600080fd5b50610342610a3a565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610382578082015181840152602081019050610367565b50505050905090810190601f1680156103af5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156103c957600080fd5b50610408600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610a73565b604051808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182151515158152602001935050505060405180910390f35b34801561046857600080fd5b506104bd600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ca1565b6040518082815260200191505060405180910390f35b6040805190810160405280600581526020017f546f6b656e00000000000000000000000000000000000000000000000000000081525081565b600081600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925846040518082815260200191505060405180910390a36001905092915050565b60008054905090565b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415151561064457600080fd5b600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054821115151561069257600080fd5b600360008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054821115151561071d57600080fd5b61076f82600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d2890919063ffffffff16565b600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555061080482600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d4190919063ffffffff16565b600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055506108d682600360008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d2890919063ffffffff16565b600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040518082815260200191505060405180910390a3600190509392505050565b601281565b6000600260008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6040805190810160405280600181526020017f540000000000000000000000000000000000000000000000000000000000000081525081565b60008060008073ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff1614151515610ab357600080fd5b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020548411151515610b0157600080fd5b610b5384600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d2890919063ffffffff16565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610be884600260008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610d4190919063ffffffff16565b600260008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508473ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef866040518082815260200191505060405180910390a3848460019250925092509250925092565b6000600360008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b6000828211151515610d3657fe5b818303905092915050565b6000808284019050838110151515610d5557fe5b80915050929150505600a165627a7a723058208d15a8072b13b5043728f4be4a53db02d76d725db16fc0777309e409345024d30029"

var superAddr string = "0x65188459a1dc65984a0c7d4a397ed3986ed0c853"
var superPriv string = "7cb4880c2d4863f88134fd01a250ef6633cc5e01aeba4c862bedbf883a148ba8"

var testAddr1 string = "0xc409aaf73698fdb5995c4d85f6033d5e90d2f2bd"
var testPriv1 string = "64d18eb7061dff419581c1af98201b76c7ab6db538b1cf65123c470ccc6d5929"

var testAddr2 string = "0x5f45ab0be3fc342e3a1033b293af0514d760ecf8"
var testPriv2 string = "d5236f9f29dfd15a8c6c42531f490b4682c6a39ef02f8de7ee2f0edab6037511"

var testAddr3 string = "0xc80bf8e2b390967dc908a4d37674473563036475"
var testPriv4 string = "fc70adae9f98412734748d9906e7524699b1e9d3eadf6b4f9b02df5fb911a4b9"

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

func TestCreateAccount1(t *testing.T) {
	time.Sleep(time.Second * 4)

	chash, code, err := client.CreateAccount(GetNonce(superAddr), superPriv, "0", "memo", superAddr, testAddr1, "1000")

	hash = chash

	t.Log(hash, code, err)

	time.Sleep(time.Second * 4)
}

func TestCreateAccount2(t *testing.T) {
	chash, code, err := client.CreateAccount(GetNonce(superAddr), superPriv, "0", "memo", superAddr, testAddr2, "1000")

	hash = chash

	t.Log(hash, code, err)

	time.Sleep(time.Second * 4)
}

func TestCreateAccount(t *testing.T) {
	chash, code, err := client.CreateAccount(GetNonce(superAddr), superPriv, "0", "memo", superAddr, testAddr3, "1000")

	hash = chash

	t.Log(hash, code, err)

	time.Sleep(time.Second * 3)
}

func TestQueryAccount(t *testing.T) {

	result, code, err := client.QueryAccount(superAddr)

	t.Log(result, code, err)

	result, code, err = client.QueryAccount(testAddr1)

	t.Log(result, code, err)

	result, code, err = client.QueryAccount(testAddr2)

	t.Log(result, code, err)

	result, code, err = client.QueryAccount(testAddr3)

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

	time.Sleep(time.Second * 4)
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

	result, code, err := client.QueryLedgerTransactions(1337, "asc", 10, 0)

	t.Log(result, code, err)
}

//*********************************************************************************************

//************************************Contract TEST********************************************

func TestCreateContract(t *testing.T) {

	param, err := NewCreateContractParam("0", "8000000", "0", bytcode, abis, []interface{}{})

	if err != nil {
		t.Log(err)
		return
	}

	ccontractHash, code, err := client.CreateContract(GetNonce(testAddr1), testPriv1, "0", "", testAddr1, param)

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

func TestQueryContract1(t *testing.T) { // 100 0000

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr1})

	t.Log(result, code, err)
}

func TestQueryContract2(t *testing.T) { // 0

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr2})

	t.Log(result, code, err)
}

func TestQueryContract3(t *testing.T) { // 0

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr3})

	t.Log(result, code, err)
}

func TestExecuteContract1(t *testing.T) {

	param, err := NewExecuteContractParam("0", "8000000", "0", "transfer", abis, []interface{}{testAddr3, 1000})

	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(testAddr1), testPriv1, "0", "", testAddr1, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)

	time.Sleep(time.Second * 4)

	receipt, code, err := client.QueryReceipt(result)

	bPayLoad, err := hex.DecodeString(receipt.Result)
	if err != nil {
		t.Log("hex decode error:", err)
		return
	}
	t.Log(receipt.Result, bPayLoad)

	outs, err := unpackResultToArray("transfer", abis, bPayLoad)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(outs)
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

func TestQueryContract4(t *testing.T) { // 99 9000
	time.Sleep(time.Second * 4)

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr1})

	t.Log(result, code, err)
}

func TestQueryContract5(t *testing.T) { // 0

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr2})

	t.Log(result, code, err)
}

func TestQueryContract6(t *testing.T) { // 1000

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr3})

	t.Log(result, code, err)
}

func TestExecuteContract2(t *testing.T) {

	param, err := NewExecuteContractParam("0", "8000000", "0", "approve", abis, []interface{}{testAddr2, 3000})

	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(testAddr1), testPriv1, "0", "", testAddr1, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)

	time.Sleep(time.Second * 4)
}

func TestQueryContract7(t *testing.T) { // 3000

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "allowance", abis, []interface{}{testAddr1, testAddr2})

	t.Log(result, code, err)
}

func TestQueryContract8(t *testing.T) { // 0

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "allowance", abis, []interface{}{testAddr1, testAddr3})

	t.Log(result, code, err)
}

func TestExecuteContract3(t *testing.T) {

	param, err := NewExecuteContractParam("0", "8000000", "0", "transferFrom", abis, []interface{}{testAddr1, testAddr3, 2000})

	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(testAddr2), testPriv2, "0", "", testAddr2, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)

	time.Sleep(time.Second * 4)

	receipt, code, err := client.QueryReceipt(result)

	t.Log(receipt, code, err)
}

func TestQueryContract9(t *testing.T) { // 1000

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "allowance", abis, []interface{}{testAddr1, testAddr2})

	t.Log(result, code, err)
}

func TestQueryContract10(t *testing.T) { // 0

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "allowance", abis, []interface{}{testAddr1, testAddr3})

	t.Log(result, code, err)
}

func TestQueryContract11(t *testing.T) { // 99 7000

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr1})

	t.Log(result, code, err)
}

func TestQueryContract12(t *testing.T) { // 0

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr2})

	t.Log(result, code, err)
}

func TestQueryContract13(t *testing.T) { // 3000

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "balanceOf", abis, []interface{}{testAddr3})

	t.Log(result, code, err)
}

//*************************************************************************************************

//*************************************Payment TEST************************************************

func TestPayment(t *testing.T) {

	result, code, err := client.Payment(GetNonce(superAddr), superPriv, "0", "memo", superAddr, testAddr1, "100")

	t.Log(result, code, err)

	hash = result

	time.Sleep(time.Second * 4)
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
