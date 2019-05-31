package smoke

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/dappledger/AnnChain-go-sdk"
	"github.com/stretchr/testify/assert"
)

const (
	accPriv        = "48deaa73f328f38d5fcb29d076b2b639c8491f97d245fc22e95a86366687903a"
	accAddr        = "28112ca022224ae7757bcd559666be5340ff109a"
	byteCode       = `608060405234801561001057600080fd5b5061020c806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806326f34b8b146100725780633c6bb436146100b35780633d4197f0146100de578063862c242e1461010b578063e1cb0e5214610142575b600080fd5b34801561007e57600080fd5b5061009d6004803603810190808035906020019092919050505061016d565b6040518082815260200191505060405180910390f35b3480156100bf57600080fd5b506100c861018d565b6040518082815260200191505060405180910390f35b3480156100ea57600080fd5b5061010960048036038101908080359060200190929190505050610193565b005b34801561011757600080fd5b50610140600480360381019080803590602001909291908035906020019092919050505061019d565b005b34801561014e57600080fd5b506101576101d7565b6040518082815260200191505060405180910390f35b600060016000838152602001908152602001600020600101549050919050565b60005481565b8060008190555050565b8160016000848152602001908152602001600020600001819055508060016000848152602001908152602001600020600101819055505050565b600080549050905600a165627a7a7230582090454576ac53e48f93db33e3502a6c1c9bff38697b8035a76500dc8ab84056b50029`
	abi            = `[{"constant": true,"inputs": [{"name": "_no","type": "uint256"}],"name": "getBatchVal","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": true,"inputs": [],"name": "val","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"},{"constant": false,"inputs": [{"name": "_val","type": "uint256"}],"name": "setVal","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": false,"inputs": [{"name": "_no","type": "uint256"},{"name": "_val","type": "uint256"}],"name": "setBatchVal","outputs": [],"payable": false,"stateMutability": "nonpayable","type": "function"},{"constant": true,"inputs": [],"name": "getVal","outputs": [{"name": "","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"}]`
	caValidatorPri = "C4743CCC9CDDBE2E6F9F990D6255B06A779AAFE369336AA70A79672326045AFCC1F68967AB7F7FBCA19F64B7D329A68AB5B9BCA629BAE6931B41AD5C5A62BF66"
	addPeerPub     = "4B4375A2B7B2424A26E5F630CC458436F8391C941D46AA62EDAB614AE7A3A812"
	isCA           = true
	power          = 100
	blockHash      = "DAFB86A83816B0BCB362DABC19661C6B9836680C"
)

var client *sdk.GoSDK

func init() {
	client = sdk.New("127.0.0.1:46657", sdk.ZaCryptoType)
}

func TestAddValidator(t *testing.T) {
	err := client.AddValidator(caValidatorPri, addPeerPub, isCA, power)
	if err != nil {
		fmt.Println("add validator fail")
	}
}

//func TestRemoveValidator(t *testing.T) {
//	err := client.RemoveValidator(caValidatorPri, addPeerPub)
//	if err != nil {
//		fmt.Println("remove validator fail")
//	}
//}

//func TestAccount(t *testing.T) {
//	account, _ := client.AccountCreate()
//	fmt.Println("account", account)

//	balance, _ := client.Balance(accAddr)
//	fmt.Println("balance:", balance)
//}

func TestZA(t *testing.T) {
	//health
	isHealth, err := client.CheckHealth()
	assert.Nil(t, err)
	assert.Equal(t, true, isHealth)
	//nonce
	nonce0, err := client.Nonce(accAddr)
	fmt.Println("=======2 nonce:", nonce0)
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
	client.ContractCreate(&arg)
	fmt.Println("======= create contract result:", result, err)
	assert.Nil(t, err)
	contractId := result["contract"].(string)

	time.Sleep(2 * time.Second)
	nonce1, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	fmt.Println("======:", nonce0, nonce1)
	//	assert.Equal(t, nonce0+1, nonce1)
	params := []interface{}{178}
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
	fmt.Println("======= txhash:", txHash)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
	receipt, err := client.Receipt(txHash)
	fmt.Println("======= receipt:", receipt, receipt.TxHash.Hex())

	spload, err := client.TransactionPayLoad(txHash)
	fmt.Println("======= payload:", spload, err)

	assert.NotNil(t, receipt.PostState)
	assert.Equal(t, txHash, receipt.TxHash.Hex())
	assert.Nil(t, err)
	nonce2, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	//	assert.Equal(t, nonce1+1, nonce2)
	callArg.Method = "getVal"
	callArg.Params = nil
	callArg.Nonce = nonce2
	resp, err := client.ContractRead(&callArg)
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(178), resp)

	nonce3, err := client.Nonce(accAddr)
	assert.Nil(t, err)
	assert.Equal(t, nonce3, nonce2)
	params = []interface{}{169}
	callArg.Method = "setVal"
	callArg.Params = params
	txHash, err = client.ContractCall(&callArg)
	fmt.Println("======= txhash:", txHash)
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
	fmt.Println("======= read contract:", resp)
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(169), resp)

	txs, count, err := client.Block(blockHash)
	assert.Nil(t, err)
	fmt.Println("====block", txs, count)
	for _, tx := range txs {
		res, err := client.Receipt(tx)
		fmt.Println("==block==receipt:", res, err)
	}

}
