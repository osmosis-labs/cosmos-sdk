package vesting

import (
	"context"

	"github.com/hashicorp/go-metrics"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

type msgServer struct {
	keeper.AccountKeeper
	types.BankKeeper
}

// NewMsgServerImpl returns an implementation of the vesting MsgServer interface,
// wrapping the corresponding AccountKeeper and BankKeeper.
func NewMsgServerImpl(k keeper.AccountKeeper, bk types.BankKeeper) types.MsgServer {
	return &msgServer{AccountKeeper: k, BankKeeper: bk}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) CreateVestingAccount(goCtx context.Context, msg *types.MsgCreateVestingAccount) (*types.MsgCreateVestingAccountResponse, error) {
	from, err := s.AccountKeeper.AddressCodec().StringToBytes(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'from' address: %s", err)
	}

	to, err := s.AccountKeeper.AddressCodec().StringToBytes(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'to' address: %s", err)
	}

	if err := validateAmount(msg.Amount); err != nil {
		return nil, err
	}

	if msg.EndTime <= 0 {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid end time")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.BankKeeper.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	if s.BankKeeper.BlockedAddr(to) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	if acc := s.AccountKeeper.GetAccount(ctx, to); acc != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.ToAddress)
	}

	baseAccount := authtypes.NewBaseAccountWithAddress(to)
	baseAccount = s.AccountKeeper.NewAccount(ctx, baseAccount).(*authtypes.BaseAccount)
	baseVestingAccount, err := types.NewBaseVestingAccount(baseAccount, msg.Amount.Sort(), msg.EndTime)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	var vestingAccount sdk.AccountI
	if msg.Delayed {
		vestingAccount = types.NewDelayedVestingAccountRaw(baseVestingAccount)
	} else {
		vestingAccount = types.NewContinuousVestingAccountRaw(baseVestingAccount, ctx.BlockTime().Unix())
	}

	s.AccountKeeper.SetAccount(ctx, vestingAccount)

	defer func() {
		telemetry.IncrCounter(1, "new", "account")

		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "create_vesting_account"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	if err = s.BankKeeper.SendCoins(ctx, from, to, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgCreateVestingAccountResponse{}, nil
}

func (s msgServer) CreatePermanentLockedAccount(goCtx context.Context, msg *types.MsgCreatePermanentLockedAccount) (*types.MsgCreatePermanentLockedAccountResponse, error) {
	from, err := s.AccountKeeper.AddressCodec().StringToBytes(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'from' address: %s", err)
	}

	to, err := s.AccountKeeper.AddressCodec().StringToBytes(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'to' address: %s", err)
	}

	if err := validateAmount(msg.Amount); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := s.BankKeeper.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, err
	}

	if s.BankKeeper.BlockedAddr(to) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	if acc := s.AccountKeeper.GetAccount(ctx, to); acc != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.ToAddress)
	}

	baseAccount := authtypes.NewBaseAccountWithAddress(to)
	baseAccount = s.AccountKeeper.NewAccount(ctx, baseAccount).(*authtypes.BaseAccount)
	vestingAccount, err := types.NewPermanentLockedAccount(baseAccount, msg.Amount)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	s.AccountKeeper.SetAccount(ctx, vestingAccount)

	defer func() {
		telemetry.IncrCounter(1, "new", "account")

		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "create_permanent_locked_account"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	if err = s.BankKeeper.SendCoins(ctx, from, to, msg.Amount); err != nil {
		return nil, err
	}

	return &types.MsgCreatePermanentLockedAccountResponse{}, nil
}

func (s msgServer) CreatePeriodicVestingAccount(goCtx context.Context, msg *types.MsgCreatePeriodicVestingAccount) (*types.MsgCreatePeriodicVestingAccountResponse, error) {
	from, err := s.AccountKeeper.AddressCodec().StringToBytes(msg.FromAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'from' address: %s", err)
	}

	to, err := s.AccountKeeper.AddressCodec().StringToBytes(msg.ToAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid 'to' address: %s", err)
	}

	if msg.StartTime < 1 {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid start time of %d, length must be greater than 0", msg.StartTime)
	}

	var totalCoins sdk.Coins
	for i, period := range msg.VestingPeriods {
		if period.Length < 1 {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid period length of %d in period %d, length must be greater than 0", period.Length, i)
		}

		if err := validateAmount(period.Amount); err != nil {
			return nil, err
		}

		totalCoins = totalCoins.Add(period.Amount...)
	}

	if s.BankKeeper.BlockedAddr(to) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if acc := s.AccountKeeper.GetAccount(ctx, to); acc != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.ToAddress)
	}

	if err := s.BankKeeper.IsSendEnabledCoins(ctx, totalCoins...); err != nil {
		return nil, err
	}

	baseAccount := authtypes.NewBaseAccountWithAddress(to)
	baseAccount = s.AccountKeeper.NewAccount(ctx, baseAccount).(*authtypes.BaseAccount)
	vestingAccount, err := types.NewPeriodicVestingAccount(baseAccount, totalCoins.Sort(), msg.StartTime, msg.VestingPeriods)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	s.AccountKeeper.SetAccount(ctx, vestingAccount)

	defer func() {
		telemetry.IncrCounter(1, "new", "account")

		for _, a := range totalCoins {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "create_periodic_vesting_account"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	if err = s.BankKeeper.SendCoins(ctx, from, to, totalCoins); err != nil {
		return nil, err
	}

	return &types.MsgCreatePeriodicVestingAccountResponse{}, nil
}

func (s msgServer) CreateClawbackVestingAccount(goCtx context.Context, msg *types.MsgCreateClawbackVestingAccount) (*types.MsgCreateClawbackVestingAccountResponse, error) {
	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, err
	}
	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, err
	}
	if err := sdk.VerifyAddressFormat(from); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address: %s", err)
	}

	if err := sdk.VerifyAddressFormat(to); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address: %s", err)
	}

	lockupCoins := sdk.NewCoins()
	for i, period := range msg.LockupPeriods {
		if period.Length < 1 {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid period length of %d in period %d, length must be greater than 0", period.Length, i)
		}
		lockupCoins = lockupCoins.Add(period.Amount...)
	}

	vestingCoins := sdk.NewCoins()
	for i, period := range msg.VestingPeriods {
		if period.Length < 1 {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid period length of %d in period %d, length must be greater than 0", period.Length, i)
		}
		vestingCoins = vestingCoins.Add(period.Amount...)
	}

	// If both schedules are present, the must describe the same total amount.
	// IsEqual can panic, so use (a == b) <=> (a <= b && b <= a).
	if len(msg.LockupPeriods) > 0 && len(msg.VestingPeriods) > 0 &&
		!(lockupCoins.IsAllLTE(vestingCoins) && vestingCoins.IsAllLTE(lockupCoins)) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "vesting and lockup schedules must have same total coins")
	}

	return &types.MsgCreateClawbackVestingAccountResponse{}, nil
}

func (s msgServer) Clawback(goCtx context.Context, msg *types.MsgClawback) (*types.MsgClawbackResponse, error) {
	funder, err := sdk.AccAddressFromBech32(msg.GetFunderAddress())
	if err != nil {
		return nil, err
	}
	if err := sdk.VerifyAddressFormat(funder); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid funder address: %s", err)
	}

	addr, err := sdk.AccAddressFromBech32(msg.GetAddress())
	if err != nil {
		return nil, err
	}
	if err := sdk.VerifyAddressFormat(addr); err != nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address: %s", err)
	}

	if msg.GetDestAddress() != "" {
		dest, err := sdk.AccAddressFromBech32(msg.GetDestAddress())
		if err != nil {
			return nil, err
		}
		if err := sdk.VerifyAddressFormat(dest); err != nil {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination address: %s", err)
		}
	}

	return &types.MsgClawbackResponse{}, nil
}

func validateAmount(amount sdk.Coins) error {
	if !amount.IsValid() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	if !amount.IsAllPositive() {
		return sdkerrors.ErrInvalidCoins.Wrap(amount.String())
	}

	return nil
}
