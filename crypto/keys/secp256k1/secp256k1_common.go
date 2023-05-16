package secp256k1

import (
	"bytes"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"

	"golang.org/x/crypto/sha3"
)

func checkPubkeyMatchesEcrecover(hash []byte, sig []byte, pubkey *secp256k1.PublicKey) bool {
	// EIP-191 allows 27, 28 for the ecrecover byte
	if sig[64] != 27 && sig[64] != 28 {
		return false
	}
	// 27 represents even, 28 represents odd
	// See https://ethereum.github.io/yellowpaper/paper.pdf, page 25. Its defined
	// two sentences above footnote 6.

	// GETH's Ecrecover only works for signatures with last byte 0 <= v < 4, so we need
	// to subtract 27 from the last byte and add it back afterwards.  Ecrecover returns a
	// pubkey bytearray in uncompressed format: 0x04 || pubkey.X || pubkey.Y

	sig[64] -= 27
	recovered_pubkey, _ := crypto.Ecrecover(hash, sig)
	sig[64] += 27

	// Serializes the supplied pubkey in uncompressed format and compares for equality
	actual_pubkey_uncompressed := pubkey.SerializeUncompressed()
	return bytes.Equal(recovered_pubkey, actual_pubkey_uncompressed)
}

func sha3Hash(msg []byte) []byte {
	// And feed the bytes into our hash
	hash := sha3.NewLegacyKeccak256()
	hash.Write(msg)
	sum := hash.Sum(nil)

	return sum
}

// Read Signature struct from R || S. Caller needs to ensure
// that len(sigStr) >= 64.
func signatureFromBytes(sigStr []byte) *secp256k1.Signature {
	return &secp256k1.Signature{
		R: new(big.Int).SetBytes(sigStr[:32]),
		S: new(big.Int).SetBytes(sigStr[32:64]),
	}
}
