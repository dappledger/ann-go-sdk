package crypto

import (
	"bytes"
	"fmt"

	"github.com/dappledger/ann-go-sdk/ed25519"
	"github.com/dappledger/ann-go-sdk/ed25519/extra25519"
	"github.com/dappledger/ann-go-sdk/wire"
)

const (
	// A series of combination of ciphers
	// ZA includes ed25519,ecdsa,ripemd160,keccak256,secretbox
	CryptoTypeZhongAn = "ZA"
)

const (
	PrivKeyLenEd25519 = 64
)

// PrivKey is part of PrivAccount and state.PrivValidator.
type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) Signature
	PubKey() PubKey
	Equals(PrivKey) bool
	KeyString() string
}

// Types of PrivKey implementations
const (
	PrivKeyTypeEd25519 = byte(0x01)
)

// for wire.readReflect
var _ = wire.RegisterInterface(
	struct{ PrivKey }{},
	wire.ConcreteType{PrivKeyEd25519{}, PrivKeyTypeEd25519},
)

func PrivKeyFromBytes(privKeyBytes []byte) (privKey PrivKey, err error) {
	err = wire.ReadBinaryBytes(privKeyBytes, &privKey)
	return
}

//-------------------------------------

// Implements PrivKey
type PrivKeyEd25519 [PrivKeyLenEd25519]byte

func (privKey PrivKeyEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeyEd25519) Sign(msg []byte) Signature {
	privKeyBytes := [PrivKeyLenEd25519]byte(privKey)
	signatureBytes := ed25519.Sign(&privKeyBytes, msg)
	return SignatureEd25519(*signatureBytes)
}

func (privKey PrivKeyEd25519) PubKey() PubKey {
	privKeyBytes := [PrivKeyLenEd25519]byte(privKey)
	return PubKeyEd25519(*ed25519.MakePublicKey(&privKeyBytes))
}

func (privKey PrivKeyEd25519) Equals(other PrivKey) bool {
	if otherEd, ok := other.(PrivKeyEd25519); ok {
		return bytes.Equal(privKey[:], otherEd[:])
	} else {
		return false
	}
}

func (privKey PrivKeyEd25519) KeyString() string {
	return fmt.Sprintf("%X", privKey[:])
}

func (privKey PrivKeyEd25519) ToCurve25519() *[32]byte {
	keyCurve25519 := new([32]byte)
	privKeyBytes := [PrivKeyLenEd25519]byte(privKey)
	extra25519.PrivateKeyToCurve25519(keyCurve25519, &privKeyBytes)
	return keyCurve25519
}

func (privKey PrivKeyEd25519) String() string {
	return fmt.Sprintf("PrivKeyEd25519{*****}")
}

// Deterministically generates new priv-key bytes from key.
func (privKey PrivKeyEd25519) Generate(index int) PrivKeyEd25519 {
	newBytes := wire.BinarySha256(struct {
		PrivKey [PrivKeyLenEd25519]byte
		Index   int
	}{privKey, index})
	var newKey [PrivKeyLenEd25519]byte
	copy(newKey[:], newBytes)
	return PrivKeyEd25519(newKey)
}

func GenPrivKeyEd25519() PrivKeyEd25519 {
	privKeyBytes := new([PrivKeyLenEd25519]byte)
	copy(privKeyBytes[:32], CRandBytes(32))
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyEd25519FromSecret(secret []byte) PrivKeyEd25519 {

	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKeyBytes := new([PrivKeyLenEd25519]byte)
	copy(privKeyBytes[:32], privKey32)
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

func GenPrivkeyByBytes(cryptoType string, data []byte) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeZhongAn:
		var ed PrivKeyEd25519
		copy(ed[:], data)
		privkey = ed
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}

func GenPrivkeyByType(cryptoType string) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeZhongAn:
		privkey = GenPrivKeyEd25519()
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}
