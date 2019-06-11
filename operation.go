package sdk

import (
	"errors"
	"math/big"
	"strings"
	"time"

	"github.com/dappledger/AnnChain/eth/common"
	"github.com/dappledger/AnnChain/eth/core/types"
	"github.com/dappledger/AnnChain/eth/crypto"
	"github.com/dappledger/AnnChain/eth/rlp"
	gtypes "github.com/dappledger/AnnChain/gemmill/types"
)

func (gs *GoSDK) get(key string) ([]byte, error) {

	if strings.HasPrefix(key, "0x") {
		key = key[2:]
	}

	bytesKey := common.Hex2Bytes(key)

	result := new(gtypes.ResultQuery)

	err := gs.sendTxCall("query", bytesKey, result)
	if err != nil {
		return nil, err
	}

	if result.Result.Code != gtypes.CodeType_OK {
		return nil, errors.New(result.Result.Error())
	}

	receipt := new(types.SReceipt)

	if err = rlp.DecodeBytes(result.Result.Data, receipt); err != nil {
		return nil, err
	}

	return receipt.Value, nil

}

func (gs *GoSDK) put(privKey string, value []byte, typ CommitType) (string, error) {

	if strings.Index(privKey, "0x") == 0 {
		privKey = privKey[2:]
	}

	privBytes := common.Hex2Bytes(privKey)

	addrBytes, err := gs.getAddrBytes(privBytes)
	if err != nil {
		return "", err
	}

	tx := types.NewTransaction(common.Bytes2Hex(addrBytes), big.NewInt(time.Now().UnixNano()), value)

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

	rpcCommitResult := new(gtypes.ResultBroadcastTxCommit)

	if typ == TypeSyn {
		err = gs.sendTxCall("broadcast_tx_commit", txBytes, rpcCommitResult)
		if err != nil {
			return "", err
		}
		return rpcCommitResult.TxHash, nil
	}

	rpcAsynResult := new(gtypes.ResultBroadcastTx)

	err = gs.sendTxCall("broadcast_tx_async", txBytes, rpcAsynResult)
	if err != nil {
		return "", err
	}

	return sigTx.Hash().Hex(), nil
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
