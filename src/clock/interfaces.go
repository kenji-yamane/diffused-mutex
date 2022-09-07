package clock

type LogicalClock interface {
	InternalEvent()
	ExternalEvent(externalClockStr string)
	GetProcessID(externalClockStr string) int
	GetClockStr() string
}
