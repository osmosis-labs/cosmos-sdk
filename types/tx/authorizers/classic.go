package authorizers

import (
	txsigning "cosmossdk.io/x/tx/signing"
	"github.com/cometbft/cometbft/libs/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type ClassicAuthorizer struct {
	accountKeeper   authsigning.AccountKeeper
	signModeHandler *txsigning.HandlerMap
}

func (c ClassicAuthorizer) SetAccountKeeper(ak authsigning.AccountKeeper) {
	c.accountKeeper = ak
}

func (c ClassicAuthorizer) SetSignModeHandler(sm *txsigning.HandlerMap) {
	c.signModeHandler = sm
}

var _ sdk.Authorizer = &ClassicAuthorizer{}

type ClassicAuthData struct {
	Signer    []byte
	Signature signing.SignatureV2
	txData    txsigning.TxData
}

func (c ClassicAuthorizer) GetAuthorizationData(tx sdk.Tx) [][]byte {
	sigTx, ok := tx.(authsigning.Tx)
	if !ok {
		return [][]byte{}
	}

	// stdSigs contains the sequence number, account number, and signatures.
	// When simulating, this would just be a 0-length slice.
	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return [][]byte{}
	}

	signers, err := sigTx.GetSigners()
	if err != nil {
		return [][]byte{}
	}

	// check that signer length and signature length are the same
	if len(sigs) != len(signers) {
		return [][]byte{}
	}

	adaptableTx, ok := tx.(authsigning.V2AdaptableTx)
	if !ok {
		return [][]byte{}
	}
	txData := adaptableTx.GetSigningTxData()

	// marshall signers in an array of bytes. [(i,j) for i,j in zip(signers, sigs)]
	result := make([][]byte, len(signers))
	for i, signer := range signers {
		classicAuthData := ClassicAuthData{
			Signer:    signer,
			Signature: sigs[i], // ToDo: better marshaling
			txData:    txData,  // ToDo: Should this be repeated for every sig?
		}
		result[i], err = json.Marshal(&classicAuthData)
	}
	return result
}

func (c ClassicAuthorizer) Authorize(ctx sdk.Context, msg []sdk.Msg, authorizationData [][]byte) bool {
	return true
	if c.accountKeeper == nil || c.signModeHandler == nil {
		return false
	}

	// ToDo: this is a hack to try to mimic the behavior of the original code while keeping the data as one per message
	classicAuthData := make([]ClassicAuthData, len(authorizationData))
	signers, signatures := make([][]byte, len(authorizationData)), make([]signing.SignatureV2, len(authorizationData))
	txData := make([]txsigning.TxData, len(authorizationData))
	for i, authData := range authorizationData {
		err := json.Unmarshal(authData, &classicAuthData[i])
		if err != nil {
			return false
		}
		signers[i] = classicAuthData[i].Signer
		signatures[i] = classicAuthData[i].Signature
		txData[i] = classicAuthData[i].txData
	}

	// ToDo: Deal with simulate. Should this be part of the iface?
	err := authsigning.ValidateSignature(ctx, signatures, signers, txData[0], false, c.accountKeeper, c.signModeHandler)
	if err != nil {
		return false
	}
	return true

}

func (c ClassicAuthorizer) ConfirmExecution(ctx sdk.Context, msg []sdk.Msg, authorized bool, authorizationData [][]byte) bool {
	// To be executed in the post handler
	return true
}
