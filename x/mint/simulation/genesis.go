package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/mint/types"
)

// Simulation parameter constants
const (
	MaxRewardPerEpoch = "max_reward_per_epoch"
	MinRewardPerEpoch = "min_reward_per_epoch"
)

// GenMaxRewardPerEpoch randomized MaxRewardPerEpoch
func GenMaxRewardPerEpoch(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(20, 2)
}

// GenMinRewardPerEpoch randomized MinRewardPerEpoch
func GenMinRewardPerEpoch(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(7, 2)
}

// RandomizedGenState generates a random GenesisState for mint
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var maxRewardPerEpoch sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MaxRewardPerEpoch, &maxRewardPerEpoch, simState.Rand,
		func(r *rand.Rand) { maxRewardPerEpoch = GenMaxRewardPerEpoch(r) },
	)

	var minRewardPerEpoch sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MinRewardPerEpoch, &minRewardPerEpoch, simState.Rand,
		func(r *rand.Rand) { minRewardPerEpoch = GenMinRewardPerEpoch(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	annualProvisions := sdk.NewDec(500000)
	epochDuration, _ := time.ParseDuration("168h") // 1 week
	epochsPerYear := int64(60 * 60 * 8766 / 5)
	params := types.NewParams(mintDenom, annualProvisions, maxRewardPerEpoch, minRewardPerEpoch, epochDuration, 156, epochsPerYear)

	mintGenesis := types.NewGenesisState(types.InitialMinter(), params, 0, 0)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
