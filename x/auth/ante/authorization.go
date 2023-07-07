package ante

import (
	errorsmod "cosmossdk.io/errors"
	txsigning "cosmossdk.io/x/tx/signing"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/authorizers"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type ClassicAuthorizationDecorator struct {
	accountKeeper   authsigning.AccountKeeper
	signModeHandler *txsigning.HandlerMap
}

func NewClassicAuthorizationDecorator(ak AccountKeeper, signModeHandler *txsigning.HandlerMap) ClassicAuthorizationDecorator {
	return ClassicAuthorizationDecorator{
		accountKeeper:   ak,
		signModeHandler: signModeHandler,
	}
}

func (vbd ClassicAuthorizationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}

	authorizer := tx.GetAuthorizer()
	classicAuthorizer, ok := authorizer.(authorizers.ClassicAuthorizer)
	if !ok {
		return ctx, errorsmod.New("auth", 2, "invalid authorizer")
	}
	classicAuthorizer.SetAccountKeeper(vbd.accountKeeper)
	classicAuthorizer.SetSignModeHandler(vbd.signModeHandler)

	if !authorizer.Authorize(ctx, tx.GetMsgs(), authorizer.GetAuthorizationData(tx)) {
		// TODO: better errors
		return ctx, errorsmod.New("auth", 3, "unauthorized")
	}

	return next(ctx, tx, simulate)
}
