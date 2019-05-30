package sdk

import (
	"bytes"
	"errors"
	"math/big"
	"strings"

	"github.com/dappledger/AnnChain-go-sdk/abi"
	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/crypto"
	"github.com/dappledger/AnnChain-go-sdk/rlp"
	"github.com/dappledger/AnnChain-go-sdk/types"
	"github.com/dappledger/AnnChain-go-sdk/utils"
)

func (gs *GoSDK) GasLimit() uint64 {
	return GasLimit
}

const (
	broadcast_tx_async  = "broadcast_tx_async"
	broadcast_tx_commit = "broadcast_tx_commit"
)

func (contract *ContractCreate) checkArgs() ([]byte, error) {

	if contract.PrivKey == "" {
		return nil, errors.New("account privkey is empty.")
	}

	if strings.Index(contract.Code, "0x") == 0 {
		contract.Code = contract.Code[2:]
	}
	if strings.Index(contract.PrivKey, "0x") == 0 {
		contract.PrivKey = contract.PrivKey[2:]
	}
	abiJson, err := abi.JSON(strings.NewReader(contract.ABI))
	if err != nil {
		return nil, err
	}
	initParam, err := parseData("", &abiJson, contract.Params)
	if err != nil {
		return nil, err
	}
	if initParam != "" && strings.Index(initParam, "0x") == 0 {
		initParam = initParam[2:]
	}

	data := common.Hex2Bytes(contract.Code + initParam)
	return data, nil
}
func (contractMethod *ContractMethod) checkArgs(isBatch bool) (abiJson abi.ABI, data []byte, err error) {

	if contractMethod.PrivKey == "" {
		err = errors.New("account privkey is empty.")
		return
	}

	if strings.Index(contractMethod.Contract, "0x") == 0 {
		contractMethod.Contract = contractMethod.Contract[2:]
	}
	if strings.Index(contractMethod.PrivKey, "0x") == 0 {
		contractMethod.PrivKey = contractMethod.PrivKey[2:]
	}
	abiJson, err = abi.JSON(strings.NewReader(contractMethod.ABI))
	if err != nil {
		return
	}
	data, err = PackCalldata(&abiJson, contractMethod.Method, contractMethod.Params, isBatch)
	return
}

func (gs *GoSDK) contractCreate(contract *ContractCreate) (map[string]interface{}, error) {

	data, err := contract.checkArgs()
	if err != nil {
		return nil, err
	}

	privBytes := common.Hex2Bytes(contract.PrivKey)
	addrBytes, err := gs.getAddrBytes(privBytes)
	if err != nil {
		return nil, err
	}

	nonce := contract.Nonce
	if nonce == 0 {
		nonce, err = gs.getNonce(common.Bytes2Hex(addrBytes))
		if err != nil {
			return nil, err
		}
	}

	tx := types.NewContractCreation(nonce, big.NewInt(0), gs.GasLimit(), big.NewInt(0), data)

	signer, sig, err := gs.signTx(privBytes, tx)
	if err != nil {
		return nil, err
	}

	sigTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return nil, err
	}
	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return nil, err
	}

	rpcResult := new(types.ResultBroadcastTxCommit)
	err = gs.sendTxCall("broadcast_tx_commit", txBytes, rpcResult)
	if err != nil {
		return nil, err
	}
	hash := rpcResult.TxHash // the same as sigTx.Hash()

	contractAddr := crypto.CreateAddress(common.BytesToAddress(addrBytes), sigTx.Nonce())

	response := map[string]interface{}{
		"tx":       hash,
		"contract": contractAddr.Hex(),
	}
	return response, nil
}

func (gs *GoSDK) contractCall(contractMethod *ContractMethod, funcType string, isBatch bool) (string, error) {

	_, data, err := contractMethod.checkArgs(isBatch)
	if err != nil {
		return "", err
	}

	privBytes := common.Hex2Bytes(contractMethod.PrivKey)

	nonce := contractMethod.Nonce
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
	toAddress := common.HexToAddress(contractMethod.Contract)

	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gs.GasLimit(), big.NewInt(0), data)

	signer, sig, err := gs.signTx(privBytes, tx)
	if err != nil {
		return "", err
	}
	var sigTx *types.Transaction
	sigTx, err = tx.WithSignature(signer, sig)
	if err != nil {
		return "", err
	}

	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return "", err
	}
	var hash string
	if strings.Contains(funcType, "commit") {
		rpcResult := new(types.ResultBroadcastTxCommit)
		err = gs.sendTxCall(funcType, txBytes, rpcResult)
		hash = rpcResult.TxHash
	} else {
		rpcResult := new(types.ResultBroadcastTx)
		err = gs.sendTxCall(funcType, txBytes, rpcResult)
		hash = rpcResult.TxHash
	}
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (gs *GoSDK) contractRead(contractMethod *ContractMethod) (interface{}, error) {

	abiJson, data, err := contractMethod.checkArgs(false)
	if err != nil {
		return nil, err
	}

	privBytes := common.Hex2Bytes(contractMethod.PrivKey)

	nonce := contractMethod.Nonce
	if nonce == 0 {
		addrBytes, err := gs.getAddrBytes(privBytes)
		if err != nil {
			return nil, err
		}
		nonce, err = gs.getNonce(common.Bytes2Hex(addrBytes))
		if err != nil {
			return nil, err
		}
	}
	toAddress := common.HexToAddress(contractMethod.Contract)

	tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gs.GasLimit(), big.NewInt(0), data)

	signer, sig, err := gs.signTx(privBytes, tx)
	if err != nil {
		return nil, err
	}
	sigTx, err := tx.WithSignature(signer, sig)
	if err != nil {
		return nil, err
	}

	txBytes, err := rlp.EncodeToBytes(sigTx)
	if err != nil {
		return nil, err
	}

	query := append([]byte{types.QueryType_Contract}, txBytes...)
	rpcResult := new(types.ResultQuery)
	err = gs.sendTxCall("query", query, rpcResult)
	if err != nil {
		return nil, err
	}
	return unpackResult(contractMethod.Method, abiJson, string(rpcResult.Result.Data))
}

func (gs *GoSDK) contractSerialCall(contractMethod *ContractMethod, isCommit bool) (ret []string, err error) {

	contractMethod.ParseHex()

	privBytes := common.Hex2Bytes(contractMethod.PrivKey)

	nonce := contractMethod.Nonce
	if nonce == 0 {
		addrBytes, err := gs.getAddrBytes(privBytes)
		if err != nil {
			return nil, err
		}
		nonce, err = gs.getNonce(common.Bytes2Hex(addrBytes))
		if err != nil {
			return nil, err
		}
	}
	toAddress := common.HexToAddress(contractMethod.Contract)

	var abiJson abi.ABI
	abiJson, err = abi.JSON(strings.NewReader(contractMethod.ABI))
	if err != nil {
		return nil, err
	}

	var bb bytes.Buffer
	var count uint32
	err = dealBatchParams(&abiJson, contractMethod.Method, contractMethod.Params, func(calldata []byte) (err error) {
		tx := types.NewTransaction(nonce, toAddress, big.NewInt(0), gs.GasLimit(), big.NewInt(0), calldata)
		var sigTx *types.Transaction
		var txBytes []byte
		signer, sig, err := gs.signTx(privBytes, tx)
		if err != nil {
			return err
		}
		sigTx, err = tx.WithSignature(signer, sig)
		if err != nil {
			return err
		}
		txBytes, err = rlp.EncodeToBytes(sigTx)
		if err != nil {
			return err
		}
		utils.WriteBytes(&bb, txBytes)
		nonce++
		count++
		return nil
	})
	if err != nil {
		return nil, err
	}
	allBytes := bb.Bytes()
	allBytes = append(utils.Uint32Bytes(count), allBytes...)
	txsBytes := types.SerialTag
	txsBytes = append(txsBytes, allBytes...)

	if isCommit {
		rpcResult := new(types.ResultBroadcastTxCommit)
		err = gs.sendTxCall("broadcast_tx_commit", txsBytes, rpcResult)
		if err != nil {
			return nil, err
		}

		txsHash := serialTxsHash(rpcResult.TxHash, int(count))

		return txsHash, nil
	}

	rpcResult := new(types.ResultBroadcastTx)
	err = gs.sendTxCall("broadcast_tx_async", txsBytes, rpcResult)
	if err != nil {
		return nil, err
	}

	txsHash := serialTxsHash(rpcResult.TxHash, int(count))

	return txsHash, nil
}

func (gs *GoSDK) signTx(privBytes []byte, tx *types.Transaction) (signer types.Signer, sig []byte, err error) {

	signer = new(types.HomesteadSigner)

	privkey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, nil, err
	}

	sig, err = crypto.Sign(signer.Hash(tx).Bytes(), privkey)

	return signer, sig, nil
}

func (gs *GoSDK) getAddrBytes(privBytes []byte) (addrBytes []byte, err error) {

	privkey, err := crypto.ToECDSA(privBytes)
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(privkey.PublicKey)
	addrBytes = addr[:]

	return addrBytes, nil
}

func serialTxsHash(res string, count int) []string {
	txsHash := make([]string, 0)
	hashLen := len(res) / count

	for len(res) >= hashLen {
		hash := res[:hashLen]
		txsHash = append(txsHash, hash)
		res = res[hashLen:]
	}

	return txsHash
}
