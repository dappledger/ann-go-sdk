package crypto

import (
	"crypto/cipher"

	"github.com/ZZMarquis/gm/sm4"
)

//--------------------------------sm2-----------------------------
func setNodePubkey_SM2(data []byte) PubKey {
	msgPubKey := PubKeyGmsm2{}
	copy(msgPubKey[:], data)
	return msgPubKey
}

func setNodePrivKey_SM2(data []byte) PrivKey {
	pk := PrivKeyGmsm2{}
	copy(pk[:], data)
	return pk
}

func setNodeSignature_SM2(data []byte) Signature {
	pk := SignatureGmsm2{}
	copy(pk[:], data)
	return pk
}

//--------------------------------ed25519-----------------------------
func setNodePubkey_ed25519(data []byte) PubKey {
	msgPubKey := PubKeyEd25519{}
	copy(msgPubKey[:], data)
	return msgPubKey
}

func setNodePrivKey_ed25519(data []byte) PrivKey {
	pk := PrivKeyEd25519{}
	copy(pk[:], data)
	return pk
}

func setNodeSignature_ed25519(data []byte) Signature {
	pk := SignatureEd25519{}
	copy(pk[:], data)
	return pk
}

//--------------------------------------------------
func NodePubkeyLen(cryptoType string) int {
	switch cryptoType {
	case CryptoTypeGM:
		return PubKeyLenGmsm2
	default:
		return PubKeyLenEd25519
	}
}

func SetNodePrivKey(cryptoType string, data []byte) PrivKey {
	switch cryptoType {
	case CryptoTypeGM:
		return setNodePrivKey_SM2(data)
	default:
		return setNodePrivKey_ed25519(data)
	}
}

//--------------------------sm4----------------------------------
const (
	BlockSize = sm4.BlockSize
)

type Sm4Cipher struct {
	cipher.Block
}

func NewCipher(key []byte) *Sm4Cipher {
	kbuf := make([]byte, BlockSize)
	copy(kbuf, key)
	subc, _ := sm4.NewCipher(kbuf)
	return &Sm4Cipher{
		subc,
	}
}

func zeroBytes(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(0)
	}
}

func (c *Sm4Cipher) cryptBlock(src []byte, encrypto bool) []byte {
	input := make([]byte, BlockSize)
	nblock := (len(src) + 0xf) >> 4
	outbuf := make([]byte, nblock*BlockSize)
	for i := 0; i < nblock; i++ {
		pos := i * BlockSize
		if i == nblock-1 && len(src[pos:]) < BlockSize { //last;
			zeroBytes(input)
		}
		copy(input, src[pos:])
		output := outbuf[pos : pos+BlockSize]
		if encrypto {
			c.Encrypt(output, input)
		} else {
			c.Decrypt(output, input)
		}
	}
	return outbuf
}

func SM4Encrypto(src, key []byte) []byte {
	c := NewCipher(key)
	return c.cryptBlock(src, true)
}

func SM4Decrypt(src, key []byte) []byte {
	c := NewCipher(key)
	return c.cryptBlock(src, false)
}
