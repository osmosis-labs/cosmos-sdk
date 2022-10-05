package types

import (
	proto "github.com/gogo/protobuf/proto"
)

// This is a very gross hack, to register gov proto in global variables, for SDK tests.
// Normally every proto file registers their proto files in a global index.
// (This is bad design)
// This causes problems for us in Osmosis, as we seek to make Osmosis' code
// the source of truth.
// We just need type definitions and SDK tests to pass.
// So we call this within SDK test helpers, which will register the global variables
// to this repo, if not already done. (Which would happen in Osmosis init)
func RegisterProtoLocallyIfNotRegistered() {
	registeredMsg := proto.MessageType("cosmos.gov.v1beta1.GenesisState")
	if registeredMsg == nil {
		proto.RegisterType((*GenesisState)(nil), "cosmos.gov.v1beta1.GenesisState")

		proto.RegisterFile("cosmos/gov/v1beta1/genesis.proto", fileDescriptor_43cd825e0fa7a627)
	}

	registeredMsg = proto.MessageType("cosmos.gov.v1beta1.MsgSubmitProposal")
	if registeredMsg == nil {
		proto.RegisterType((*MsgSubmitProposal)(nil), "cosmos.gov.v1beta1.MsgSubmitProposal")
		proto.RegisterType((*MsgSubmitProposalResponse)(nil), "cosmos.gov.v1beta1.MsgSubmitProposalResponse")
		proto.RegisterType((*MsgVote)(nil), "cosmos.gov.v1beta1.MsgVote")
		proto.RegisterType((*MsgVoteResponse)(nil), "cosmos.gov.v1beta1.MsgVoteResponse")
		proto.RegisterType((*MsgVoteWeighted)(nil), "cosmos.gov.v1beta1.MsgVoteWeighted")
		proto.RegisterType((*MsgVoteWeightedResponse)(nil), "cosmos.gov.v1beta1.MsgVoteWeightedResponse")
		proto.RegisterType((*MsgDeposit)(nil), "cosmos.gov.v1beta1.MsgDeposit")
		proto.RegisterType((*MsgDepositResponse)(nil), "cosmos.gov.v1beta1.MsgDepositResponse")

		proto.RegisterFile("cosmos/gov/v1beta1/tx.proto", fileDescriptor_3c053992595e3dce)
	}

	registeredMsg = proto.MessageType("cosmos.gov.v1beta1.WeightedVoteOption")
	if registeredMsg == nil {
		proto.RegisterEnum("cosmos.gov.v1beta1.VoteOption", VoteOption_name, VoteOption_value)
		proto.RegisterEnum("cosmos.gov.v1beta1.ProposalStatus", ProposalStatus_name, ProposalStatus_value)
		proto.RegisterType((*WeightedVoteOption)(nil), "cosmos.gov.v1beta1.WeightedVoteOption")
		proto.RegisterType((*TextProposal)(nil), "cosmos.gov.v1beta1.TextProposal")
		proto.RegisterType((*Deposit)(nil), "cosmos.gov.v1beta1.Deposit")
		proto.RegisterType((*Proposal)(nil), "cosmos.gov.v1beta1.Proposal")
		proto.RegisterType((*TallyResult)(nil), "cosmos.gov.v1beta1.TallyResult")
		proto.RegisterType((*Vote)(nil), "cosmos.gov.v1beta1.Vote")
		proto.RegisterType((*DepositParams)(nil), "cosmos.gov.v1beta1.DepositParams")
		proto.RegisterType((*VotingParams)(nil), "cosmos.gov.v1beta1.VotingParams")
		proto.RegisterType((*TallyParams)(nil), "cosmos.gov.v1beta1.TallyParams")
		proto.RegisterType((*ProposalVotingPeriod)(nil), "cosmos.gov.v1beta1.ProposalVotingPeriod")

		proto.RegisterFile("cosmos/gov/v1beta1/gov.proto", fileDescriptor_6e82113c1a9a4b7c)
	}

	registeredMsg = proto.MessageType("cosmos.gov.v1beta1.QueryProposalRequest")
	if registeredMsg == nil {
		proto.RegisterType((*QueryProposalRequest)(nil), "cosmos.gov.v1beta1.QueryProposalRequest")
		proto.RegisterType((*QueryProposalResponse)(nil), "cosmos.gov.v1beta1.QueryProposalResponse")
		proto.RegisterType((*QueryProposalsRequest)(nil), "cosmos.gov.v1beta1.QueryProposalsRequest")
		proto.RegisterType((*QueryProposalsResponse)(nil), "cosmos.gov.v1beta1.QueryProposalsResponse")
		proto.RegisterType((*QueryVoteRequest)(nil), "cosmos.gov.v1beta1.QueryVoteRequest")
		proto.RegisterType((*QueryVoteResponse)(nil), "cosmos.gov.v1beta1.QueryVoteResponse")
		proto.RegisterType((*QueryVotesRequest)(nil), "cosmos.gov.v1beta1.QueryVotesRequest")
		proto.RegisterType((*QueryVotesResponse)(nil), "cosmos.gov.v1beta1.QueryVotesResponse")
		proto.RegisterType((*QueryParamsRequest)(nil), "cosmos.gov.v1beta1.QueryParamsRequest")
		proto.RegisterType((*QueryParamsResponse)(nil), "cosmos.gov.v1beta1.QueryParamsResponse")
		proto.RegisterType((*QueryDepositRequest)(nil), "cosmos.gov.v1beta1.QueryDepositRequest")
		proto.RegisterType((*QueryDepositResponse)(nil), "cosmos.gov.v1beta1.QueryDepositResponse")
		proto.RegisterType((*QueryDepositsRequest)(nil), "cosmos.gov.v1beta1.QueryDepositsRequest")
		proto.RegisterType((*QueryDepositsResponse)(nil), "cosmos.gov.v1beta1.QueryDepositsResponse")
		proto.RegisterType((*QueryTallyResultRequest)(nil), "cosmos.gov.v1beta1.QueryTallyResultRequest")
		proto.RegisterType((*QueryTallyResultResponse)(nil), "cosmos.gov.v1beta1.QueryTallyResultResponse")

		proto.RegisterFile("cosmos/gov/v1beta1/query.proto", fileDescriptor_e35c0d133e91c0a2)
	}
}
