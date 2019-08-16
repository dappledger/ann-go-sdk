package crypto

import (
	"github.com/dappledger/AnnChain-go-sdk/go-hash"
)

var (
	node_crypto_type = CryptoTypeZhongAn //default value;
	CryptoType       = CryptoTypeZhongAn
)

func GetNodeCryptoType() string {
	return node_crypto_type
}

func NodeInit(crypto string) {
	switch crypto {
	case CryptoTypeZhongAn:
		node_crypto_type = crypto
		hash.ConfigHasher(hash.HashTypeRipemd160)
	default:
		hash.ConfigHasher(hash.HashTypeRipemd160)
	}
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
	return PubKeyLenEd25519
}

func SetNodePrivKey(cryptoType string, data []byte) PrivKey {
	return setNodePrivKey_ed25519(data)
}

func SetNodePubkey(data []byte) PubKey {
	return setNodePubkey_ed25519(data)
}

func GetNodeSigBytes(sig Signature) []byte {
	gsig := sig.(SignatureEd25519)
	return gsig[:]
}

func GetNodePubkeyBytes(pkey PubKey) []byte {
	gpkey := pkey.(PubKeyEd25519)
	return gpkey[:]
}
