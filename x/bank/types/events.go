package types

// bank module event types
const (
	EventTypeTransfer = "transfer"

	AttributeKeyRecipient = "recipient"
	AttributeKeySender    = "sender"

	AttributeValueCategory = ModuleName
)

var (
	AttributeKeyRecipientBytes = []byte(AttributeKeyRecipient)
	AttributeKeySenderBytes    = []byte(AttributeKeySender)
)
