package crypto

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
	default:
		return PubKeyLenEd25519
	}
}

func SetNodePrivKey(cryptoType string, data []byte) PrivKey {
	switch cryptoType {
	default:
		return setNodePrivKey_ed25519(data)
	}
}

func zeroBytes(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(0)
	}
}
