<!--
order: 2
-->

# State

## Minter

The minter is a space for holding current rewards information.

 - Minter: `0x00 -> ProtocolBuffer(minter)`

+++ https://github.com/cosmos/cosmos-sdk/blob/v0.40.0-rc7/proto/cosmos/mint/v1beta1/mint.proto#L8-L19

## Params

Minting params are held in the global params store. 

 - Params: `mint/params -> amino(params)`

```go
type Params struct {
	MintDenom           string  // type of coin to mint
	InflationRateChange sdk.Dec // maximum annual change in inflation rate
	InflationMax        sdk.Dec // maximum inflation rate
	InflationMin        sdk.Dec // minimum inflation rate
	GoalBonded          sdk.Dec // goal of percent bonded atoms
	BlocksPerYear       int64   // expected blocks per year
}
```
