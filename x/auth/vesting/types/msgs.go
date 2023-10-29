package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateVestingAccount{}
	_ sdk.Msg = &MsgCreateClawbackVestingAccount{}
	_ sdk.Msg = &MsgClawback{}
	_ sdk.Msg = &MsgCreatePermanentLockedAccount{}
	_ sdk.Msg = &MsgCreatePeriodicVestingAccount{}
)

// NewMsgCreateVestingAccount returns a reference to a new MsgCreateVestingAccount.
func NewMsgCreateVestingAccount(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins, endTime int64, delayed bool) *MsgCreateVestingAccount {
	return &MsgCreateVestingAccount{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
		EndTime:     endTime,
		Delayed:     delayed,
	}
}

func NewMsgCreateClawbackVestingAccount(fromAddr, toAddr sdk.AccAddress, startTime int64, lockupPeriods, vestingPeriods []Period, merge bool) *MsgCreateClawbackVestingAccount {
	return &MsgCreateClawbackVestingAccount{
		FromAddress:    fromAddr.String(),
		ToAddress:      toAddr.String(),
		StartTime:      startTime,
		LockupPeriods:  lockupPeriods,
		VestingPeriods: vestingPeriods,
		Merge:          merge,
	}
}

func NewMsgClawback(funderAddress, address, destAddress string) *MsgClawback {
	return &MsgClawback{
		FunderAddress: funderAddress,
		Address:       address,
		DestAddress:   destAddress,
	}
}

// NewMsgCreatePermanentLockedAccount returns a reference to a new MsgCreatePermanentLockedAccount.
func NewMsgCreatePermanentLockedAccount(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) *MsgCreatePermanentLockedAccount {
	return &MsgCreatePermanentLockedAccount{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
	}
}

// NewMsgCreatePeriodicVestingAccount returns a reference to a new MsgCreatePeriodicVestingAccount.
func NewMsgCreatePeriodicVestingAccount(fromAddr, toAddr sdk.AccAddress, startTime int64, periods []Period) *MsgCreatePeriodicVestingAccount {
	return &MsgCreatePeriodicVestingAccount{
		FromAddress:    fromAddr.String(),
		ToAddress:      toAddr.String(),
		StartTime:      startTime,
		VestingPeriods: periods,
	}
}
