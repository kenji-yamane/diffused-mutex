package src

type State string

const (
	Released State = "released"
	Wanted   State = "wanted"
	Held     State = "held"

	ConsumeCmd = "x"

	SharedResourcePort = "10001"
)
