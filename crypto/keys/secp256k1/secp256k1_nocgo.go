//go:build !libsecp256k1

package secp256k1

import (
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec"

	"github.com/tendermint/tendermint/crypto"
)

// used to reject malleable signatures
// see:
//  - https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
//  - https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/crypto.go#L39
var secp256k1halfN = new(big.Int).Rsh(secp256k1.S256().N, 1)

// Sign creates an ECDSA signature on curve Secp256k1, using SHA256 on the msg.
// The returned signature will be of the form R || S (in lower-S form).
func (privKey *PrivKey) Sign(msg []byte) ([]byte, error) {
	priv, _ := secp256k1.PrivKeyFromBytes(secp256k1.S256(), privKey.Key)
	sig, err := priv.Sign(crypto.Sha256(msg))
	if err != nil {
		return nil, err
	}
	sigBytes := serializeSig(sig)
	return sigBytes, nil
}

// VerifySignature verifies signatures of the forms:
// * (R || S), using sha256 hashing
// * (R || S || V), using keccak256 hashing
// where V is an ECRecover byte compatible with Ethereum.
// EIP-191 signatures from metamask return V as one of {27, 28}
// It rejects signatures which are not in lower-S form.
func (pubKey *PubKey) VerifySignature(msg []byte, sigStr []byte) bool {
	sigLen := len(sigStr)
	// Invalid signature length
	if sigLen != 64 && sigLen != 65 {
		return false
	}
	pub, err := secp256k1.ParsePubKey(pubKey.Key, secp256k1.S256())
	if err != nil {
		return false
	}
	// parse the signature:
	signature := signatureFromBytes(sigStr)
	// Reject malleable signatures. libsecp256k1 does this check but btcec doesn't.
	// see: https://github.com/ethereum/go-ethereum/blob/f9401ae011ddf7f8d2d95020b7446c17f8d98dc1/crypto/signature_nocgo.go#L90-L93
	// This check is also part of the Ethereum signature verification specification,
	// see: https://ethereum.github.io/yellowpaper/paper.pdf, Berlin version, page 29
	if signature.S.Cmp(secp256k1halfN) > 0 {
		return false
	}
	if sigLen == 64 { // Use sha256 hashing for (R,S) signature form
		return signature.Verify(crypto.Sha256(msg), pub)
	} else {
		// Ensure that "v" is canonical
		v := sigStr[64]
		ecRecoverValid := isECRecoverByteValid(v, pub)
		if !ecRecoverValid {
			return false
		}
		// Use keccak256 hashing for (R,S) signature form
		validSig := signature.Verify(sha3Hash(msg), pub)
		return validSig
	}
}

// Serialize signature to R || S.
// R, S are padded to 32 bytes respectively.
func serializeSig(sig *secp256k1.Signature) []byte {
	rBytes := sig.R.Bytes()
	sBytes := sig.S.Bytes()
	sigBytes := make([]byte, 64)
	// 0 pad the byte arrays from the left if they aren't big enough.
	copy(sigBytes[32-len(rBytes):32], rBytes)
	copy(sigBytes[64-len(sBytes):64], sBytes)
	return sigBytes
}
