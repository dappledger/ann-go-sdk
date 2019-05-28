package crypto

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/ZZMarquis/gm/sm2"

	"github.com/dappledger/AnnChain-go-sdk/ed25519"
	"github.com/dappledger/AnnChain-go-sdk/ed25519/extra25519"
	"github.com/dappledger/AnnChain-go-sdk/wire"
)

const (
	// A series of combination of ciphers
	// ZA includes ed25519,ecdsa,ripemd160,keccak256,secretbox
	// GM includes SM2,SM3,SM4
	CryptoTypeZhongAn = "ZA"
	CryptoTypeGM      = "GM"
)

type PrivKey interface {
	Bytes() []byte
	Sign(msg []byte) Signature
	PubKey() PubKey
	Equals(PrivKey) bool
	KeyString() string
}

type PrivKeyEd25519 [64]byte

func (privKey PrivKeyEd25519) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeyEd25519) Sign(msg []byte) Signature {
	privKeyBytes := [64]byte(privKey)
	signatureBytes := ed25519.Sign(&privKeyBytes, msg)
	return SignatureEd25519(*signatureBytes)
}

func (privKey PrivKeyEd25519) PubKey() PubKey {
	privKeyBytes := [64]byte(privKey)
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
	privKeyBytes := [64]byte(privKey)
	extra25519.PrivateKeyToCurve25519(keyCurve25519, &privKeyBytes)
	return keyCurve25519
}

func (privKey PrivKeyEd25519) String() string {
	return ("PrivKeyEd25519{*****}")
}

// Deterministically generates new priv-key bytes from key.
func (privKey PrivKeyEd25519) Generate(index int) PrivKeyEd25519 {
	newBytes := wire.BinarySha256(struct {
		PrivKey [64]byte
		Index   int
	}{privKey, index})
	var newKey [64]byte
	copy(newKey[:], newBytes)
	return PrivKeyEd25519(newKey)
}

func GenPrivKeyEd25519() PrivKeyEd25519 {
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], CRandBytes(32))
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

// NOTE: secret should be the output of a KDF like bcrypt,
// if it's derived from user input.
func GenPrivKeyEd25519FromSecret(secret []byte) PrivKeyEd25519 {
	privKey32 := Sha256(secret) // Not Ripemd160 because we want 32 bytes.
	privKeyBytes := new([64]byte)
	copy(privKeyBytes[:32], privKey32)
	ed25519.MakePublicKey(privKeyBytes)
	return PrivKeyEd25519(*privKeyBytes)
}

//-------------------------------------
// PrivKeyGmsm2 Implements PrivKey
type PrivKeyGmsm2 [32]byte

func (privKey PrivKeyGmsm2) Bytes() []byte {
	return wire.BinaryBytes(struct{ PrivKey }{privKey})
}

func (privKey PrivKeyGmsm2) Sign(msg []byte) Signature {
	sm2key, err := sm2.RawBytesToPrivateKey(privKey[:])
	if err != nil {
		panic(fmt.Errorf("Sign failed 1: %v", err))
	}
	sdata, err := sm2.Sign(sm2key, nil, msg)
	if err != nil {
		panic(fmt.Errorf("Sign failed 2: %v", err))
	}

	return SM2SignToSignatureGmsm2(sdata)
}

func (privKey PrivKeyGmsm2) PubKey() PubKey {
	priv, err := sm2.RawBytesToPrivateKey(privKey[:])
	if err != nil {
		panic(fmt.Errorf("Sign failed 1: %v", err))
	}
	pub := new(sm2.PublicKey)
	pub.Curve = priv.Curve
	pub.X, pub.Y = priv.Curve.ScalarBaseMult(priv.D.Bytes())
	pubBytes := pub.GetRawBytes()
	if len(pubBytes) != 64 {
		panic(fmt.Errorf("PrivKeyGmsm2.PubKey error(len=%d)", len(pubBytes)))
	}
	var pubKey PubKeyGmsm2
	copy(pubKey[:], pubBytes)
	return pubKey
}

func (privKey PrivKeyGmsm2) Equals(other PrivKey) bool {
	if otherSecp, ok := other.(PrivKeyGmsm2); ok {
		return bytes.Equal(privKey[:], otherSecp[:])
	} else {
		return false
	}
}

func (privKey PrivKeyGmsm2) String() string {
	return ("PrivKeyGmsm2{*****}")
}

func (privKey PrivKeyGmsm2) KeyString() string {
	return fmt.Sprintf("%X", privKey[:])
}

func GenPrivKeyGmsm2() PrivKeyGmsm2 {
	priv, _, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	privdata := priv.GetRawBytes()
	if len(privdata) != 32 {
		panic(fmt.Errorf("GenPrivKeyGmsm2 error(len=%d)", len(privdata)))
	}
	var privKey PrivKeyGmsm2
	copy(privKey[:], privdata)
	return privKey
}

func GenPrivkeyByBytes(cryptoType string, data []byte) (PrivKey, error) {
	var privkey PrivKey
	switch cryptoType {
	case CryptoTypeZhongAn:
		var ed PrivKeyEd25519
		copy(ed[:], data)
		privkey = ed
	case CryptoTypeGM:
		var gm PrivKeyGmsm2
		copy(gm[:], data)
		privkey = gm
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
	case CryptoTypeGM:
		privkey = GenPrivKeyGmsm2()
	default:
		return nil, fmt.Errorf("Unknow crypto type")
	}
	return privkey, nil
}
