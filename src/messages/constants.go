package messages

type MessageType string

const (
	Request MessageType = "request"
	Reply   MessageType = "reply"
	Consume MessageType = "consume"
)

func (m MessageType) String() string {
	return string(m)
}
