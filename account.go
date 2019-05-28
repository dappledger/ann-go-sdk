package sdk

import (
	"github.com/dappledger/AnnChain-go-sdk/common"
	"github.com/dappledger/AnnChain-go-sdk/crypto"
)

type Account struct {
	Privkey string `json:"privkey"`
	Address string `json:"address"`
}

func (gs *GoSDK) accountCreate() (Account, error) {
	var account Account

	privkey, err := crypto.GenerateKey()
	if err != nil {
		return Account{}, err
	}

	address := crypto.PubkeyToAddress(privkey.PublicKey)

	account.Privkey = common.Bytes2Hex(crypto.FromECDSA(privkey))
	account.Address = address.Hex()

	return account, nil
}
