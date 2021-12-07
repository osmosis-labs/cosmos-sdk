package secp256k1

import (
	"fmt"
	"math/big"

	secp256k1 "github.com/btcsuite/btcd/btcec"

	"golang.org/x/crypto/sha3"
)

func isECRecoverByteValid(v byte, pubkey *secp256k1.PublicKey) bool {
	// EIP-191 allows 27, 28
	if v != 27 && v != 28 {
		return false
	}
	// 27 represents even, 28 represents odd
	// See https://ethereum.github.io/yellowpaper/paper.pdf, page 25. Its defined
	// two sentences above footnote 6.
	// When using the above, we verify all GETH signatures correctly.
	EIP_191_ECRecoverStandard := true
	// _HOWEVER_ we are aiming to verify EIP-191 signatures.
	// Inexplicably, every wallet when making an EIP-191 signature has the
	// v_parity bit flipped (so 28 when you'd ex)
	v_parity := 1 - (v % 2)
	v_parity_bool := v_parity == 1
	if EIP_191_ECRecoverStandard {
		v_parity = 1 - v_parity
		v_parity_bool = !v_parity_bool
	}
	actual_parity := pubkey.Y.Bit(0)
	actual_parity_bool := actual_parity == 1
	fmt.Println(v, v_parity, actual_parity)
	return v_parity_bool == actual_parity_bool
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
