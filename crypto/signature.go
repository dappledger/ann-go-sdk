package crypto

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ZZMarquis/gm/sm2"
	"github.com/dappledger/AnnChain-go-sdk/wire"
)

const (
	SignKeyLenEd25519 = 64
	SignKeyLenGmsm2   = 64
)

type Signature interface {
	Bytes() []byte
	IsZero() bool
	String() string
	Equals(Signature) bool
	KeyString() string
}

type SignatureEd25519 [64]byte

func (sig SignatureEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ Signature }{sig})
}

func (sig SignatureEd25519) IsZero() bool { return len(sig) == 0 }

func (sig SignatureEd25519) String() string { return fmt.Sprintf("/%X.../", sig[:6]) }

func (sig SignatureEd25519) Equals(other Signature) bool {
	if otherEd, ok := other.(SignatureEd25519); ok {
		return bytes.Equal(sig[:], otherEd[:])
	} else {
		return false
	}
}

func (sig SignatureEd25519) KeyString() string {
	return fmt.Sprintf("%X", sig[:])
}

// SignatureGmsm2 Implements Signature
type SignatureGmsm2 [64]byte

func PaddedBigBytes(bi *big.Int, buf []byte) {
	n := len(buf)
	bb := bi.Bytes()
	offset := n - len(bb)
	if offset < 0 {
		copy(buf, bb)
	} else {
		copy(buf[offset:], bb)
	}
}

func SM2SignToSignatureGmsm2(sdata []byte) SignatureGmsm2 {
	r, s, err := sm2.UnmarshalSign(sdata)
	var sig SignatureGmsm2
	if err != nil {
		return sig
	}
	PaddedBigBytes(r, sig[:32])
	PaddedBigBytes(s, sig[32:])
	return sig
}

func SignatureGmsm2ToSM2Sign(sig SignatureGmsm2) []byte {
	r := new(big.Int).SetBytes(sig[:32])
	s := new(big.Int).SetBytes(sig[32:])
	smdata, _ := sm2.MarshalSign(r, s)
	return smdata
}

func (sig SignatureGmsm2) Bytes() []byte {
	return wire.BinaryBytes(struct{ Signature }{sig})
}

func (sig SignatureGmsm2) IsZero() bool { return len(sig) == 0 }

func Fingerprint(slice []byte) []byte {
	fingerprint := make([]byte, 6)
	copy(fingerprint, slice)
	return fingerprint
}
func (sig SignatureGmsm2) String() string { return fmt.Sprintf("/%X.../", Fingerprint(sig[:])) }

func (sig SignatureGmsm2) Equals(other Signature) bool {
	if otherEd, ok := other.(SignatureGmsm2); ok {
		return bytes.Equal(sig[:], otherEd[:])
	}
	return false
}

func (sig SignatureGmsm2) KeyString() string {
	return fmt.Sprintf("%X", sig[:])
}
