package message

type MessageType int

const (
	InvitationType MessageType = iota
	ConnectionRequestType
	ConnectionResponseType
	ConnectionACKType

	ProposalCredentialType
	OfferCredentialType
	RequestCredentialType
	IssueCredentialType
	CredentialACKType

	RequestPresentationType
	PresentationType
	PresentationACKType
	GeneralMsgType
)

type Message struct {
	MessageType
	Content   interface{}
	JsonBytes []byte
}
