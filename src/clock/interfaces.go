package clock

type LogicalClock interface {
	InternalEvent()
	ExternalEvent(externalClockStr string)
	GetClockStr() string
	CompareClocks(requestClockStr string, externalClockStr string, externalId int) (int, error)
}
