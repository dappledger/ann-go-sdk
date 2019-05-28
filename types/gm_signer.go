package types

import (
	"errors"
	"math/big"

	"github.com/ZZMarquis/gm/sm2"
	"github.com/ZZMarquis/gm/sm3"
	"github.com/dappledger/AnnChain-go-sdk/common"
)

type Sm2Signer struct{}

func (s Sm2Signer) Equal(s2 Signer) bool {
	_, ok := s2.(Sm2Signer)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format.
func (ss Sm2Signer) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes(sig[64:])
	return r, s, v, nil
}

func (ss Sm2Signer) Sender(tx *Transaction) (common.Address, error) {
	return recoverPlainSm2(ss.Hash(tx), tx.data.R, tx.data.S, tx.data.V)
}

func (ss Sm2Signer) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
	})
}

func recoverPlainSm2(txHash common.Hash, r, s, v *big.Int) (common.Address, error) {
	//Verify,sig = sig+pubBytes
	vBytes := v.Bytes()
	pubBytes := vBytes[len(vBytes)-64:]
	pubkey, err := sm2.RawBytesToPublicKey(pubBytes)
	if err != nil {
		return common.Address{}, err
	}
	rBytes, sBytes := r.Bytes(), s.Bytes()
	sig := make([]byte, len(vBytes)) // len(r+s+v) - len(pubBytes) = len(v)
	copy(sig[32-len(rBytes):32], rBytes)
	copy(sig[64-len(sBytes):64], sBytes)
	copy(sig[64:], vBytes[:len(vBytes)-64])
	ok := sm2.Verify(pubkey, nil, txHash.Bytes(), sig)
	if !ok {
		return common.Address{}, errors.New("Sm2 Verify failed")
	}

	d := sm3.New()
	d.Write(pubBytes)
	hash := d.Sum(nil)[12:]
	var addr common.Address
	copy(addr[:], hash)
	return addr, nil
}
