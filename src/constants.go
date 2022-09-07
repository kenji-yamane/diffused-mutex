package src

type MessageType string

const (
	Request MessageType = "request"
	Reply   MessageType = "reply"
	Consume MessageType = "consume"

	ConsumeCmd = "x"
)

func (m MessageType) String() string {
	return string(m)
}
