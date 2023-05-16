//go:build !libsecp256k1
// +build !libsecp256k1

package secp256k1

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"gitee.com/aqchain/go-ethereum/crypto/sha3"
	secp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

type TestCase = struct {
	msg        string
	pubkey     types.PubKey
	sig        string
	expectPass bool
	tcName     string
}

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

func hashPersonalMessage(data []byte) []byte {
	msg := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))

	fmt.Printf("message: %x\n", data)
	fmt.Printf("prehash: %x\n", msg)
	h := sha3.NewKeccak256() // Create a new instance of the hash
	h.Write(msg)             // Write the data you wanted hashed to it
	out := h.Sum(nil)        // Finalize the checksum
	fmt.Printf("hash: %x\n", out)
	return out
}

// Test thing
func TestGenerateRandomEIP191(t *testing.T) {
	const NUM_TEST_CASES = 10

	cases := make([]TestCase, NUM_TEST_CASES)

	for i := range cases {
		privKey := make([]byte, 32)
		// Reads in random bytes to private key
		_, err := rand.Read(privKey)
		privKeyECDSA, err := crypto.ToECDSA(privKey)
		privKeyCosmos := PrivKey{Key: privKey}
		pubKeyCosmos := privKeyCosmos.PubKey()
		require.NoError(t, err)

		// Message is just going to be something random
		message := make([]byte, 64)
		_, err = rand.Read(message)

		// Computes H(EIP191(message))
		hash := hashPersonalMessage([]byte(message))
		sig, err := crypto.Sign(hash, privKeyECDSA)
		sig[64] += 27
		testcase := TestCase{
			msg:        eip191MsgTransform(string(message)),
			pubkey:     pubKeyCosmos,
			sig:        hex.EncodeToString(sig),
			expectPass: true,
			tcName:     fmt.Sprintf("Test case %d", i),
		}

		cases[i] = testcase
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
		fmt.Println("Test name:", tc.tcName)
		fmt.Println("Actual signature: ", hex.EncodeToString(sigBz))
		require.Equal(t, tc.expectPass, sigVerif, "Verify signature didn't act as expected, tc %v", tc.tcName)
	}
}

// TODO: Cleanup this test, separate into multiple tests
// TODO: Make a test in x/auth for testing initializing pubkey from metamask.
// TODO: Make a test in x/auth for accepting multiple sig types
// TODO: Update cgo implementation
func estEthSignatureVerification(t *testing.T) {
	metamaskDemoPrivkeyBz, err := hex.DecodeString("4545454545454545454545454545454545454545454545454545454545454545")
	metamaskDemoPrivkey := PrivKey{Key: metamaskDemoPrivkeyBz}
	metamaskDemoPubkey := metamaskDemoPrivkey.PubKey()
	require.NoError(t, err)
	// This message appears to not be the message getting signed
	metamaskDemoSignMsg := `{
  "account_number": "0",
  "chain_id": "testing",
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
  ],
  "sequence": "0"
}`

	// TODO: Reget this signature, not sure if I have the right one or not.
	metamaskDemoSig := "8625b77a3352aa62608528757765969413f16b8d1f0242b96cda41f2f45dc18756e100e49eaa7d737929a088bf49ce3d5e11ef33a0e074f28b74c401770eb15d1b"
	//// Metamask test vector setup from https://github.com/MetaMask/eth-sig-util/blob/main/src/personal-sign.test.ts
	metamask69PrivkeyBz, _ := hex.DecodeString("6969696969696969696969696969696969696969696969696969696969696969")
	metamask69Privkey := PrivKey{Key: metamask69PrivkeyBz}
	metamask69Pubkey := metamask69Privkey.PubKey()
	fmt.Println("Pubkey: ", metamask69Pubkey)
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
		{eip191MsgTransform(metamaskDemoSignMsg), metamaskDemoPubkey, metamaskDemoSig, true, "metamask Demo EIP 191 test case"},
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
		fmt.Println("Test name:", tc.tcName)
		fmt.Println("Actual signature: ", hex.EncodeToString(sigBz))
		require.Equal(t, tc.expectPass, sigVerif, "Verify signature didn't act as expected, tc %v", tc.tcName)
	}
}
