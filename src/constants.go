package src

type MessageType string

type State string

const (
	Request MessageType = "request"
	Reply   MessageType = "reply"
	Consume MessageType = "consume"

	Released State = "released"
	Wanted   State = "wanted"
	Held     State = "held"

	ConsumeCmd = "x"
)

func (m MessageType) String() string {
	return string(m)
}
