package sdk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dappledger/AnnChain-go-sdk/common/hexutil"
	"github.com/dappledger/AnnChain-go-sdk/crypto"
	gtypes "github.com/dappledger/AnnChain/gemmill/types"
)

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

	vAttr := &gtypes.ValidatorAttr{
		pubkey,
		uint64(power),
		isCa,
	}

	scmd := &gtypes.SpecialOPCmd{}
	scmd.CmdType = gtypes.SpecialOP_ChangeValidator
	scmd.Time = time.Now()
	if err := scmd.LoadMsg(vAttr); err != nil {
		return err
	}

	scmd.PubKey, _ = hex.DecodeString(privKey.PubKey().KeyString())
	signMessage, _ := json.Marshal(scmd)
	scmd.Signature, _ = hex.DecodeString(privKey.Sign(signMessage).KeyString())
	cmdBytes, _ := json.Marshal(scmd)

	rpcResult := new(gtypes.ResultRequestSpecialOP)
	err := gs.JsonRPCCall("request_special_op", gtypes.TagSpecialOPTx(cmdBytes), rpcResult)
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

	scmd := &gtypes.SpecialOPCmd{}
	scmd.CmdType = gtypes.SpecialOP_Disconnect
	scmd.Time = time.Now()
	if err := scmd.LoadMsg(pubkey); err != nil {
		return err
	}

	scmd.PubKey, _ = hex.DecodeString(privKey.PubKey().KeyString())
	signMessage, _ := json.Marshal(scmd)
	scmd.Signature, _ = hex.DecodeString(privKey.Sign(signMessage).KeyString())
	cmdBytes, _ := json.Marshal(scmd)

	rpcResult := new(gtypes.ResultRequestSpecialOP)
	err := gs.JsonRPCCall("request_special_op", gtypes.TagSpecialOPTx(cmdBytes), rpcResult)
	if err != nil {
		return err
	}

	fmt.Println(*rpcResult)

	return nil
}

func (gs *GoSDK) Validators() (*gtypes.ResultValidators, error) {
	rpcResult := new(gtypes.ResultValidators)
	err := gs.JsonRPCCall("validators", []byte{}, rpcResult)
	if err != nil {
		return nil, err
	}

	return rpcResult, nil
}
