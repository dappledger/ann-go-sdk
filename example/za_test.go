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
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/dappledger/ann-go-sdk/rlp"
	"github.com/stretchr/testify/assert"

	sdk "github.com/dappledger/ann-go-sdk"
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

func TestKV(t *testing.T) {
	client := sdk.New("localhost:46657", sdk.ZaCryptoType)

	nonce1, err := client.Nonce(accAddr)
	assert.Nil(t, err)

	var arg = sdk.KVTx{
		AccountBase: sdk.AccountBase{
			PrivKey: accPriv,
			Nonce:   nonce1,
		},
		Key:   []byte("key1"),
		Value: []byte("value1"),
	}

	sig, err := KVSignature(&arg)
	assert.Nil(t, err)
	_, err = client.PutSignature(sig)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	value1, err := client.Get([]byte("key1"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("value1"), value1)

	arg.Nonce, err = client.Nonce(accAddr)
	assert.Nil(t, err)
	arg.Key = []byte("key2")
	arg.Value = []byte("value2")
	_, err = client.Put(&arg)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	value2, err := client.Get([]byte("key2"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("value2"), value2)

	arg.Nonce, err = client.Nonce(accAddr)
	assert.Nil(t, err)
	arg.Key = []byte("key3")
	arg.Value = []byte("value3")
	_, err = client.Put(&arg)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	value3, err := client.Get([]byte("key3"))
	assert.Nil(t, err)
	assert.Equal(t, []byte("value3"), value3)

	arg.Nonce, err = client.Nonce(accAddr)
	assert.Nil(t, err)
	arg.Key = []byte("key3")
	arg.Value = []byte("value3")
	_, err = client.Put(&arg)
	assert.NotNil(t, err)
	assert.True(t, true, strings.HasPrefix(err.Error(), "duplicate key"))

	kvs, err := client.GetWithPrefix([]byte("k"), []byte("key1"), 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(kvs))
	assert.Equal(t, &sdk.KVResult{Key: []byte("key2"), Value: []byte("value2")}, kvs[0])
	assert.Equal(t, &sdk.KVResult{Key: []byte("key3"), Value: []byte("value3")}, kvs[1])
	for _, kv := range kvs {
		t.Log(string(kv.Key), string(kv.Value))
	}
}

func KVSignature(kvTx *sdk.KVTx) (string, error) {
	if kvTx.PrivKey == "" {
		return "", fmt.Errorf("account privkey is empty")
	}

	if strings.Index(kvTx.PrivKey, "0x") == 0 {
		kvTx.PrivKey = kvTx.PrivKey[2:]
	}

	privBytes := common.Hex2Bytes(kvTx.PrivKey)
	kvBytes, err := rlp.EncodeToBytes(&sdk.KVResult{Key: kvTx.Key, Value: kvTx.Value})
	if err != nil {
		return "", err
	}

	txdata := append(sdk.KVTxType, kvBytes...)
	tx := types.NewTransaction(kvTx.Nonce, common.Address{}, big.NewInt(0), sdk.GasLimit, big.NewInt(0), txdata)
	signer, sig, err := signTx(privBytes, tx)
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

	return common.Bytes2Hex(txBytes), nil
}

func signTx(privBytes []byte, tx *types.Transaction) (signer types.Signer, sig []byte, err error) {
	signer = new(types.HomesteadSigner)

	privkey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, nil, err
	}

	sig, err = crypto.Sign(signer.Hash(tx).Bytes(), privkey)

	return signer, sig, nil
}

func TestPayloadTx(t *testing.T) {
	client := sdk.New("localhost:46657", sdk.ZaCryptoType)

	nonce1, err := client.Nonce(accAddr)
	assert.Nil(t, err)

	var arg = sdk.Tx{
		AccountBase: sdk.AccountBase{
			PrivKey: accPriv,
			Nonce:   nonce1,
		},
		Payload: "value1",
	}

	sig, err := PayloadTxSignature(&arg)
	assert.Nil(t, err)
	txHash, err := client.TranscationSignature(sig)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)

	value1, err := client.TransactionPayLoad(txHash)
	assert.Nil(t, err)
	assert.Equal(t, "value1", value1)

	arg.Nonce, err = client.Nonce(accAddr)
	assert.Nil(t, err)
	arg.Payload = "value2"
	txHash2, err := client.Transaction(&arg)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
	value2, err := client.TransactionPayLoad(txHash2)
	assert.Equal(t, "value2", value2)
}

func PayloadTxSignature(payloadTx *sdk.Tx) (string, error) {
	if payloadTx.PrivKey == "" {
		return "", fmt.Errorf("account privkey is empty")
	}

	if strings.Index(payloadTx.PrivKey, "0x") == 0 {
		payloadTx.PrivKey = payloadTx.PrivKey[2:]
	}

	privBytes := common.Hex2Bytes(payloadTx.PrivKey)

	payload := payloadTx.Payload
	data := []byte(payload)

	tx := types.NewTransaction(payloadTx.Nonce, common.Address{}, payloadTx.Value, sdk.GasLimit, big.NewInt(0), data)

	signer, sig, err := signTx(privBytes, tx)
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

	return common.Bytes2Hex(txBytes), nil
}


func TestHttps(t *testing.T) {
	goSdk := sdk.New("baas-poc3.zhongan.io",sdk.ZaCryptoType)
	h,err := goSdk.LastHeight()
	assert.Error(t,err)
	if h!=0 {
		t.Fatal(fmt.Sprintf("height must be zero"))
	}

	goSdk,err =  sdk.NewSDk("https://baas-poc3.zhongan.io",sdk.ZaCryptoType)
	assert.NoError(t,err)
	h,err = goSdk.LastHeight()
	assert.NoError(t,err)
	if h<10 {
		t.Fatal(fmt.Sprintf("height is to small %d",h))
	}
}
func TestPendingNonce(t *testing.T) {
	client := sdk.New("localhost:46657", sdk.ZaCryptoType)

	var txsHash []string
	for i := 0; i < 10; i++ {
		nonce, err := client.PendingNonce(accAddr)
		assert.Nil(t, err)

		var arg = sdk.Tx{
			AccountBase: sdk.AccountBase{
				PrivKey: accPriv,
				Nonce:   nonce,
			},
			Payload: "value"+ fmt.Sprintf("%v", i),
		}

		sig, err := PayloadTxSignature(&arg)
		assert.Nil(t, err)
		// async
		txHash, err := client.TranscationSignatureAsync(sig)
		assert.Nil(t, err)

		txsHash = append(txsHash, txHash)
	}

	time.Sleep(3 * time.Second)

	for n := 0; n < 10; n++ {
		value, err := client.TransactionPayLoad(txsHash[n])
		assert.Nil(t, err)
		assert.Equal(t, "value"+ fmt.Sprintf("%v", n), value)
	}
}