package crypto

import (
	"bytes"
	"fmt"

	"github.com/dappledger/ann-go-sdk/wire"
)

const (
	SignKeyLenEd25519 = 64
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
