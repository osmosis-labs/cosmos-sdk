//go:build !libsecp256k1
// +build !libsecp256k1

package secp256k1

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/stretchr/testify/require"
)

// Ensure that signature verification works, and that
// non-canonical signatures fail.
// Note: run with CGO_ENABLED=0 or go test -tags !cgo.
func TestSignatureVerificationAndRejectUpperS(t *testing.T) {
	msg := []byte("We have lingered long enough on the shores of the cosmic ocean.")
	for i := 0; i < 500; i++ {
		priv := GenPrivKey()
		sigStr, err := priv.Sign(msg)
		require.NoError(t, err)
		sig := signatureFromBytes(sigStr)
		require.False(t, sig.S.Cmp(secp256k1halfN) > 0)

		pub := priv.PubKey()
		require.True(t, pub.VerifySignature(msg, sigStr))

		// malleate:
		sig.S.Sub(secp256k1.S256().CurveParams.N, sig.S)
		require.True(t, sig.S.Cmp(secp256k1halfN) > 0)
		malSigStr := serializeSig(sig)

		require.False(t, pub.VerifySignature(msg, malSigStr),
			"VerifyBytes incorrect with malleated & invalid S. sig=%v, key=%v",
			sig,
			priv,
		)
	}
}

func eip191MsgTransform(msg string) string {
	return "\x19Ethereum Signed Message:\n" + fmt.Sprintf("%d", len(msg)) + msg
}

func TestEthSignatureVerification(t *testing.T) {
	metamaskDemoPubkey := "BDGnzLx3JYDk5sBm9iNxufEs4BKE3pneA3FacHAV98mK1s0Pykt8ADZSuhO8o/7xyGigjsuyyPDtJLi2ePF2Dm0="
	metamaskDemoPubkeyBz, _ := base64.StdEncoding.DecodeString(metamaskDemoPubkey)
	// This message appears to not be the message getting signed
	metamaskDemoSignMsg := `{
  "chain_id": "testing",
  "account_number": "0",
  "sequence": "0",
  "fee": {
    "amount": [
      {
        "amount": "100",
        "denom": "ucosm"
      }
    ],
    "gas": "250"
  },
  "memo": "Some memo",
  "msgs": [
    {
      "type": "cosmos-sdk/MsgSend",
      "value": {
        "amount": [
          {
            "amount": "1234567",
            "denom": "ucosm"
          }
        ],
        "from_address": "cosmos1tru96ya986ta2lruqeh9fsleca7ucuzpwqjhvr",
        "to_address": "cosmos1tru96ya986ta2lruqeh9fsleca7ucuzpwqjhvr"
      }
    }
  ]
}`
	fmt.Println(sha256.Sum256([]byte(metamaskDemoSignMsg)))
	metamaskDemoPubkey = metamaskDemoPubkey
	metamaskDemoSignMsg = metamaskDemoSignMsg

	metamaskDemoSig := "f0992543f357cfb6e614271b37867377fd8027833579a10967f0282d894e3efd797f950bd342ad519c229740c6017d43f1ebd9186e2baca32dfbd31fa0c4fcf91b"
	//// Metamask test vector setup from https://github.com/MetaMask/eth-sig-util/blob/main/src/personal-sign.test.ts
	metamask69PrivkeyBz, _ := hex.DecodeString("6969696969696969696969696969696969696969696969696969696969696969")
	metamask69Privkey := PrivKey{Key: metamask69PrivkeyBz}
	metamask69Pubkey := metamask69Privkey.PubKey()

	// gethtestmsg, _ := hex.DecodeString("ce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971008")
	// gethtestsig := "90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc931c"
	// gethtestpubkey, _ := hex.DecodeString("04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	// gethtestpubkeyc, _ := hex.DecodeString("02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")

	cases := []struct {
		msg        string
		pubkey     types.PubKey
		sig        string
		expectPass bool
		tcName     string
	}{
		// {string(gethtestmsg), &PubKey{Key: testpubkey}, gethtestsig, true, "geth test vector #0  -- require hashing disabled"},
		// {string(testmsg), &PubKey{testpubkeyc}, testsig, true, "geth test vector #1 -- require hashing disabled"},
		{eip191MsgTransform("hello world"), metamask69Pubkey,
			"ce909e8ea6851bc36c007a0072d0524b07a3ff8d4e623aca4c71ca8e57250c4d0a3fc38fa8fbaaa81ead4b9f6bd03356b6f8bf18bccad167d78891636e1d69561b",
			true, "metamask test vector #1"},
		{eip191MsgTransform(metamaskDemoSignMsg), &PubKey{metamaskDemoPubkeyBz}, metamaskDemoSig, true, "metamask Demo EIP 191 test case"},
	}

	for _, tc := range cases {
		// pubkeyBz, err := base64.StdEncoding.DecodeString(tc.pubkey)
		// require.NoError(t, err, tc.tcName)
		// require.Len(t, pubkeyBz, PubKeySize, tc.tcName)
		// pk := PubKey{Key: pubkeyBz}
		pk := tc.pubkey
		sigBz, err := hex.DecodeString(tc.sig)
		require.NoError(t, err, tc.tcName)
		msgBz := []byte(tc.msg)
		sigVerif := pk.VerifySignature(msgBz, sigBz)
		require.Equal(t, tc.expectPass, sigVerif, "Verify signature didn't act as expected, tc %v", tc.tcName)
	}
}
