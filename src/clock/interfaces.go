package clock

type LogicalClock interface {
	InternalEvent()
	ExternalEvent(externalClockStr string)
	GetClockStr() string
}
