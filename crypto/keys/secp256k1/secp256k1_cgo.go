//go:build libsecp256k1

package secp256k1

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1/internal/secp256k1"
)

// Sign creates an ECDSA signature on curve Secp256k1, using SHA256 on the msg.
func (privKey *PrivKey) Sign(msg []byte) ([]byte, error) {
	rsv, err := secp256k1.Sign(crypto.Sha256(msg), privKey.Key)
	if err != nil {
		return nil, err
	}
	// we do not need v  in r||s||v:
	rs := rsv[:len(rsv)-1]
	return rs, nil
}

// VerifySignature validates the signature.
// The msg will be hashed prior to signature verification.
func (pubKey *PubKey) VerifySignature(msg []byte, sig []byte) bool {
	var msgHash []byte
	// if len(sig) == 64, use Sha256
	// if len(sig) == 65, check & remove ECRecover byte, use Sha3
	if len(sig) == 64 {
		msgHash = crypto.Sha256(msg)
	} else if len(sig) == 65 {
		msgHash = sha3Hash(msg)
		v := sig[64]
		ecRecoverValid := isECRecoverByteValid(v, pub)
		if !ecRecoverValid {
			return false
		}
		sig = sig[:64]
	} else {
		return false
	}
	return secp256k1.VerifySignature(pubKey.Key, msgHash, sig)
}
