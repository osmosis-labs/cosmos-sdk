package secp256k1

import (
	"bytes"
	"math/big"

	"gitee.com/aqchain/go-ethereum/crypto"
	secp256k1 "github.com/btcsuite/btcd/btcec"

	"golang.org/x/crypto/sha3"
)

func isECRecoverByteValid(v byte, msg []byte, sig []byte, pubkey *secp256k1.PublicKey) bool {
	// EIP-191 allows 27, 28
	if v != 27 && v != 28 {
		return false
	}
	// 27 represents even, 28 represents odd
	// See https://ethereum.github.io/yellowpaper/paper.pdf, page 25. Its defined
	// two sentences above footnote 6.
	// When using the above, we verify all GETH signatures correctly.
	// _HOWEVER_ we are aiming to verify EIP-191 signatures.
	// Inexplicably, every wallet when making an EIP-191 signature has the
	// v_parity bit flipped (so 28 when you'd expect 27)
	recovered, _ := crypto.Ecrecover(msg, sig)
	uncompressed := pubkey.SerializeUncompressed()

	return bytes.Equal(recovered, uncompressed)
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
