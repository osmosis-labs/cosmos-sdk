package types

import (
	"testing"

	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestOsmosisProtoHack(t *testing.T) {
	// ensure things are registered after one call
	RegisterProtoLocallyIfNotRegistered()
	registeredMsg1 := proto.MessageType("cosmos.gov.v1beta1.MsgSubmitProposal")
	require.NotNil(t, registeredMsg1)
	// ensure repeat call doesn't mess anything up
	RegisterProtoLocallyIfNotRegistered()
	registeredMsg2 := proto.MessageType("cosmos.gov.v1beta1.MsgSubmitProposal")
	require.Equal(t, registeredMsg1, registeredMsg2)
}
