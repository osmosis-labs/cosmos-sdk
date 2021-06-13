package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/bank/simulation"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abonormal scenarios are not tested here.
func TestRandomizedGenState(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          cdc,
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: 1000,
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)

	var bankGenesis types.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &bankGenesis)

	require.Equal(t, true, bankGenesis.Params.GetDefaultSendEnabled())
	require.Len(t, bankGenesis.Params.GetSendEnabled(), 1)
	require.Len(t, bankGenesis.Balances, 3)
	require.Equal(t, "cosmos1ghekyjucln7y67ntx7cf27m9dpuxxemn4c8g4r", bankGenesis.Balances[2].GetAddress().String())
	require.Equal(t, simulation.GenesisCoins("stake", 1000).String(), bankGenesis.Balances[2].GetCoins().String())
	require.Equal(t, "6000stake,3000stake00,3000stake01,3000stake02,3000stake03,3000stake04,3000stake05,3000stake06,3000stake07,3000stake08,3000stake09,3000stake10,3000stake11,3000stake12,3000stake13,3000stake14,3000stake15,3000stake16,3000stake17,3000stake18,3000stake19", bankGenesis.Supply.String())
}

// TestRandomizedGenState tests abnormal scenarios of applying RandomizedGenState.
func TestRandomizedGenState1(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	s := rand.NewSource(1)
	r := rand.New(s)

	// all these tests will panic
	tests := []struct {
		simState module.SimulationState
		panicMsg string
	}{
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{}, "invalid memory address or nil pointer dereference"},
		{ // panic => reason: incomplete initialization of the simState
			module.SimulationState{
				AppParams: make(simtypes.AppParams),
				Cdc:       cdc,
				Rand:      r,
			}, "assignment to entry in nil map"},
	}

	for _, tt := range tests {
		require.Panicsf(t, func() { simulation.RandomizedGenState(&tt.simState) }, tt.panicMsg)
	}
}
