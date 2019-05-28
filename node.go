package sdk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dappledger/AnnChain-go-sdk/common/hexutil"
	"github.com/dappledger/AnnChain-go-sdk/crypto"
	"github.com/dappledger/AnnChain-go-sdk/types"
)

func (gs *GoSDK) checkHealth() (bool, error) {
	rpcResult := new(types.ResultHealthInfo)
	err := gs.sendTxCall("healthinfo", nil, rpcResult)
	if err != nil {
		return false, err
	}

	return (200 == rpcResult.Status), nil
}

func (gs *GoSDK) RemoveValidator(priv, pub string) error {
	return gs.changeValidator(priv, pub, true, 0)
}

func (gs *GoSDK) AddValidator(priv, pub string, isCa bool, power int64) error {
	return gs.changeValidator(priv, pub, isCa, power)
}

func (gs *GoSDK) changeValidator(priv, pub string, isCa bool, power int64) error {
	cryptoType := string(gs.cryptoType)
	data := hexutil.FromHex(priv)
	privKey := crypto.SetNodePrivKey(cryptoType, data)

	pubkey := hexutil.FromHex(pub)

	if len(pubkey) != crypto.NodePubkeyLen(cryptoType) {
		return fmt.Errorf("pubkey format error:need %d's bytes;but %d", crypto.NodePubkeyLen(cryptoType), len(pubkey))
	}

	vAttr := &types.ValidatorAttr{
		pubkey,
		uint64(power),
		isCa,
	}

	scmd := &types.SpecialOPCmd{}
	scmd.CmdType = types.SpecialOP_ChangeValidator
	scmd.Time = time.Now()
	if err := scmd.LoadMsg(vAttr); err != nil {
		return err
	}

	scmd.PubKey, _ = hex.DecodeString(privKey.PubKey().KeyString())
	signMessage, _ := json.Marshal(scmd)
	scmd.Signature, _ = hex.DecodeString(privKey.Sign(signMessage).KeyString())
	cmdBytes, _ := json.Marshal(scmd)

	rpcResult := new(types.ResultRequestSpecialOP)
	err := gs.JsonRPCCall("request_special_op", types.TagSpecialOPTx(cmdBytes), rpcResult)
	if err != nil {
		return err
	}

	fmt.Println(*rpcResult)

	return nil
}

//remove node(pubkey) from this system
func (gs *GoSDK) DisconnectPerr(priv, pub string) error {
	cryptoType := string(gs.cryptoType)
	data := hexutil.FromHex(priv)
	privKey := crypto.SetNodePrivKey(cryptoType, data)

	pubkey := hexutil.FromHex(pub)

	scmd := &types.SpecialOPCmd{}
	scmd.CmdType = types.SpecialOP_Disconnect
	scmd.Time = time.Now()
	if err := scmd.LoadMsg(pubkey); err != nil {
		return err
	}

	scmd.PubKey, _ = hex.DecodeString(privKey.PubKey().KeyString())
	signMessage, _ := json.Marshal(scmd)
	scmd.Signature, _ = hex.DecodeString(privKey.Sign(signMessage).KeyString())
	cmdBytes, _ := json.Marshal(scmd)

	rpcResult := new(types.ResultRequestSpecialOP)
	err := gs.JsonRPCCall("request_special_op", types.TagSpecialOPTx(cmdBytes), rpcResult)
	if err != nil {
		return err
	}

	fmt.Println(*rpcResult)

	return nil
}

func (gs *GoSDK) Validators() (*types.ResultValidators, error) {
	rpcResult := new(types.ResultValidators)
	err := gs.JsonRPCCall("validators", []byte{}, rpcResult)
	if err != nil {
		return nil, err
	}

	return rpcResult, nil
}
