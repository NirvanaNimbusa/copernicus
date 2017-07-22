package core

import "github.com/btcboost/secp256k1-go/secp256k1"

type Signature secp256k1.EcdsaSignature

func (sig *Signature) toLibEcdsaSignature() *secp256k1.EcdsaSignature {
	return (*secp256k1.EcdsaSignature)(sig)
}

func (sig *Signature) Serialize() []byte {
	_, serializedSig, _ := secp256k1.EcdsaSignatureSerializeDer(secp256k1Context, sig.toLibEcdsaSignature())
	return serializedSig
}

func (sig *Signature) Verify(hash []byte, pubKey *PublicKey) bool {
	correct, _ := secp256k1.EcdsaVerify(secp256k1Context, sig.toLibEcdsaSignature(),
		hash, (*secp256k1.PublicKey)(pubKey))
	return correct == 1
}

func ParseDERSignature(signature []byte) (*Signature, error) {
	_, sig, err := secp256k1.EcdsaSignatureParseDer(secp256k1Context, signature)
	if err != nil {
		return nil, err
	}
	return (*Signature)(sig), nil
}

func ParseSignature(signature []byte) (*Signature, error) {
	_, sig, err := secp256k1.EcdsaSignatureParseCompact(secp256k1Context, signature)
	if err != nil {
		return nil, err
	}
	return (*Signature)(sig), nil
}

func Sign(privKey *PrivateKey, hash []byte) ([]byte, error) {
	_, signature, err := secp256k1.EcdsaSign(secp256k1Context, hash, privKey.D.Bytes())
	if err != nil {
		return nil, err
	}
	_, ret, _ := secp256k1.EcdsaSignatureSerializeDer(secp256k1Context, signature)
	return ret, nil
}

/**
 * IsValidSignatureEncoding  A canonical signature exists of: <30> <total len> <02> <len R> <R> <02> <len S> <S> <hashtype>
 * Where R and S are not negative (their first byte has its highest bit not set), and not
 * excessively padded (do not start with a 0 byte, unless an otherwise negative number follows,
 * in which case a single 0 byte is necessary and even required).
 *
 * See https://bitcointalk.org/index.php?topic=8392.msg127623#msg127623
 *
 * This function is consensus-critical since BIP66.
 */

func IsValidSignatureEncoding(signs []byte) bool {
	// Format: 0x30 [total-length] 0x02 [R-length] [R] 0x02 [S-length] [S] [sighash]
	// * total-length: 1-byte length descriptor of everything that follows,
	//   excluding the sighash byte.
	// * R-length: 1-byte length descriptor of the R value that follows.
	// * R: arbitrary-length big-endian encoded R value. It must use the shortest
	//   possible encoding for a positive integers (which means no null bytes at
	//   the start, except a single one when the next byte has its highest bit set).
	// * S-length: 1-byte length descriptor of the S value that follows.
	// * S: arbitrary-length big-endian encoded S value. The same rules apply.
	// * sighash: 1-byte value indicating what data is hashed (not part of the DER
	//   signature)
	signsLen := len(signs)
	if signsLen < 9 {
		return false
	}
	if signsLen > 73 {
		return false
	}
	if signs[0] != 0x30 {
		return false
	}
	if int(signs[1]) != (signsLen - 3) {
		return false
	}
	lenR := signs[3]
	if int(5+lenR) >= signsLen {
		return false
	}
	lenS := signs[5+lenR]
	if int(lenR+lenS+7) != signsLen {
		return false
	}
	if signs[2] != 0x02 {
		return false
	}
	if lenR == 0 {
		return false
	}
	if (signs[4] & 0x80) == 1 {
		return false
	}
	if lenR > 1 && (signs[4] == 0x00) && (signs[5]&0x80) != 1 {
		return false
	}
	if signs[lenR+4] != 0x02 {
		return false
	}
	if lenS == 0 {
		return false
	}
	if signs[lenR+6]&0x80 == 1 {
		return false
	}
	if lenS > 1 && (signs[lenR+6] == 0x00) && (signs[lenR+7]&0x80) != 1 {
		return false
	}
	return true

}
