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
	"sync"
	"testing"
	"time"

	ethcmn "github.com/dappledger/AnnChain/genesis/eth/common"
)

var abis string = `[
	{
		"constant": false,
		"inputs": [
			{
				"name": "number",
				"type": "int256"
			}
		],
		"name": "testC",
		"outputs": [
			{
				"name": "",
				"type": "int256"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_per",
		"outputs": [
			{
				"name": "age",
				"type": "uint8"
			},
			{
				"name": "name",
				"type": "string"
			},
			{
				"name": "addr",
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
		"name": "queryPersonAge",
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
		"inputs": [],
		"name": "_numC",
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
		"constant": true,
		"inputs": [],
		"name": "queryPersonName",
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
				"name": "number",
				"type": "uint256"
			}
		],
		"name": "testB",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "number",
				"type": "uint256"
			}
		],
		"name": "testA",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "_numB",
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
		"name": "queryPersonAddress",
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
		"name": "_owner",
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
		"name": "_numA",
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
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]`

var bytcode string = "608060405234801561001057600080fd5b506009600190815560028054600160a060020a03191633179055600019600355600860009081556040805180820190915260048082527f6b656c65000000000000000000000000000000000000000000000000000000006020909201918252919261007d9290919061010a565b50600280548282018054600160a060020a031916600160a060020a0390921691909117905581546004805460ff191660ff9092169190911781556001838101805485946100de93600593926000199181161561010002919091011604610188565b5060029182015491018054600160a060020a031916600160a060020a039092169190911790555061021a565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061014b57805160ff1916838001178555610178565b82800160010185558215610178579182015b8281111561017857825182559160200191906001019061015d565b506101849291506101fd565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106101c15780548555610178565b8280016001018555821561017857600052602060002091601f016020900482015b828111156101785782548255916001019190600101906101e2565b61021791905b808211156101845760008155600101610203565b90565b6104ed806102296000396000f3006080604052600436106100ae5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166330650a3281146100b357806335b873cd146100dd57806336142497146101925780637cea649d146101bd5780639b5c7961146101d25780639cad65241461025c5780639d4853c414610274578063a36000d31461028c578063ac9d8b1e146102a1578063b2bdfa7b146102df578063c0abc010146102f4575b600080fd5b3480156100bf57600080fd5b506100cb600435610309565b60408051918252519081900360200190f35b3480156100e957600080fd5b506100f2610311565b6040805160ff8516815273ffffffffffffffffffffffffffffffffffffffff831691810191909152606060208083018281528551928401929092528451608084019186019080838360005b8381101561015557818101518382015260200161013d565b50505050905090810190601f1680156101825780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b34801561019e57600080fd5b506101a76103c8565b6040805160ff9092168252519081900360200190f35b3480156101c957600080fd5b506100cb6103d1565b3480156101de57600080fd5b506101e76103d7565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610221578181015183820152602001610209565b50505050905090810190601f16801561024e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561026857600080fd5b506100cb60043561046d565b34801561028057600080fd5b506100cb600435610475565b34801561029857600080fd5b506100cb61047d565b3480156102ad57600080fd5b506102b6610483565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b3480156102eb57600080fd5b506102b661049f565b34801561030057600080fd5b506100cb6104bb565b600381905590565b600480546005805460408051602060026101006001861615026000190190941693909304601f810184900484028201840190925281815260ff909416949392918301828280156103a25780601f10610377576101008083540402835291602001916103a2565b820191906000526020600020905b81548152906001019060200180831161038557829003601f168201915b5050506002909301549192505073ffffffffffffffffffffffffffffffffffffffff1683565b60045460ff1690565b60035481565b60058054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156104635780601f1061043857610100808354040283529160200191610463565b820191906000526020600020905b81548152906001019060200180831161044657829003601f168201915b5050505050905090565b600181905590565b600081905590565b60015481565b60065473ffffffffffffffffffffffffffffffffffffffff1690565b60025473ffffffffffffffffffffffffffffffffffffffff1681565b600054815600a165627a7a72305820dc2f4e4b2d1cdc5ec7830c60152ea0d7e4203b44af3b1fb692aea3511ed8291c0029"

var superAddr string = "0x65188459a1dc65984a0c7d4a397ed3986ed0c853"
var superPriv string = "7cb4880c2d4863f88134fd01a250ef6633cc5e01aeba4c862bedbf883a148ba8"

var testAddr string = "0xc409aaf73698fdb5995c4d85f6033d5e90d2f2bd"
var testPriv string = "64d18eb7061dff419581c1af98201b76c7ab6db538b1cf65123c470ccc6d5929"

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

//**************************Account TEST *************************************

func TestCreateAccount(t *testing.T) {

	chash, code, err := client.CreateAccount(GetNonce(superAddr), superPriv, "100", "memo", superAddr, testAddr, "1000")

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

	time.Sleep(time.Second)
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

	param, err := NewExecuteContractParam("1", "10000000", "0", "GetRand", abis, []interface{}{1, 200})

	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(superAddr), superPriv, "100", "", superAddr, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)

	time.Sleep(time.Second)
}

func TestQueryExecuteReceipt(t *testing.T) {

	result, code, err := client.QueryReceipt(contractHash)

	bResult, _ := hex.DecodeString(result.Result)

	outs, err := unpackResultToArray("testA", abis, bResult)

	t.Log(result, outs, code, err)
}

func TestExecuteContract2(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	param, err := NewExecuteContractParam("1", "8000000", "0", "testB", abis, []interface{}{299})
	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(superAddr), superPriv, "100", "", superAddr, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)

	time.Sleep(time.Second)
}

func TestQueryExecuteReceipt2(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryReceipt(contractHash)

	bResult, _ := hex.DecodeString(result.Result)

	outs, err := unpackResultToArray("testB", abis, bResult)

	t.Log(result, outs, code, err)
}

func TestExecuteContract3(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	param, err := NewExecuteContractParam("1", "8000000", "0", "testC", abis, []interface{}{-199})
	if err != nil {
		t.Log(err)
		return
	}

	result, code, err := client.ExcuteContract(GetNonce(superAddr), superPriv, "100", "", superAddr, contractAddress, param)

	contractHash = result

	t.Log(result, code, err)

	time.Sleep(time.Second)
}

func TestQueryExecuteReceipt3(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryReceipt(contractHash)

	bResult, _ := hex.DecodeString(result.Result)

	outs, err := unpackResultToArray("testC", abis, bResult)

	t.Log(result, outs, code, err)
}

func TestQueryContract(t *testing.T) {

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "GetResult", abis, []interface{}{})

	t.Log(result, ethcmn.ToHex(result.(ethcmn.Address).Bytes()), code, err)
}

func TestQueryContract2(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "_numA", abis, []interface{}{})

	t.Log(result, code, err)
}

func TestQueryContract3(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "_numB", abis, []interface{}{})

	t.Log(result, code, err)
}

func TestQueryContract4(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "queryPersonAge", abis, []interface{}{})

	t.Log(result, code, err)
}

func TestQueryContract5(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "queryPersonAddress", abis, []interface{}{})

	t.Log(result, ethcmn.ToHex(result.(ethcmn.Address).Bytes()), code, err)
}

func TestQueryContract6(t *testing.T) {

	client := NewAnnChainClient("tcp://127.0.0.1:46657")

	result, code, err := client.QueryContract(superPriv, superAddr, contractAddress, "_numC", abis, []interface{}{})

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

//******************************Stressing Test*****************************************************

func BenchmarkCreateAddress(t *testing.B) {

	for i := 0; i < routines; i++ {

		newPriv, newAddr := GenerateKey()

		_, _, err := client.CreateAccount(GetNonce(superAddr), superPriv, "0", "memo", superAddr, newAddr, "100000")

		if err == nil {
			mapAddress[newAddr] = newPriv
		}

		time.Sleep(time.Millisecond * 500)
	}

	t.Log("success number:", len(mapAddress))
}

func BenchmarkPayment(t *testing.B) {

	var waitGroup sync.WaitGroup

	for address, privkey := range mapAddress {

		waitGroup.Add(1)

		go func(address, privkey string, waitGroup sync.WaitGroup) {

			for i := 0; i < sustained; i++ {

				client.Payment(GetNonce(address), privkey, "0", "", address, superAddr, "1")

				time.Sleep(time.Second)

			}

			waitGroup.Done()

		}(address, privkey, waitGroup)
	}

	waitGroup.Wait()
}

//*************************************************************************************************
