<!--
order: 3
-->

# End Block

At each `EndBlock`, the fees received are transferred to the distribution `ModuleAccount`, as it's the account the one who keeps track of the flow of coins in (as in this case) and out the module. Unlike the mainline cosmos-SDK, we do not implement any sort of proposer bonus, and the community pool doesn't have a stream of tx fees. The fees are distributed proportionally by voting power to all bonded validators independent of whether they voted (social distribution). Note the social distribution is applied to proposer validator in addition to the proposer reward.

All provision rewards are added to a provision reward pool which validator holds individually (`ValidatorDistribution.ProvisionsRewardPool`).

```go
func AllocateTokens(feesCollected sdk.Coins, feePool FeePool, proposer ValidatorDistribution, 
              sumPowerPrecommitValidators, totalBondedTokens, communityTax, 
              proposerCommissionRate sdk.Dec)

     SendCoins(FeeCollectorAddr, DistributionModuleAccAddr, feesCollected)
     feesCollectedDec = MakeDecCoins(feesCollected)

     // communityTax param currently set to 0 on Osmosis-1
     communityFunding = feesCollectedDec * communityTax
     feePool.CommunityFund += communityFunding

     poolReceived = feesCollectedDec - communityFunding
     feePool.Pool += poolReceived

     SetValidatorDistribution(proposer)
     SetFeePool(feePool)
```
