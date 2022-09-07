package clock

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/math"
)

type ScalarClock struct {
	id    int
	ticks int
}

func NewScalarClock(id int) LogicalClock {
	return &ScalarClock{
		id:    id,
		ticks: 0,
	}
}

func (c *ScalarClock) InternalEvent() {
	c.ticks++
	c.echoClock()
}

func (c *ScalarClock) ExternalEvent(externalClockStr string) {
	externalClock, err := c.parse(externalClockStr)
	if err != nil {
		fmt.Println("invalid clock string, ignoring...")
		return
	}
	c.ticks = math.Max(externalClock.Ticks, c.ticks) + 1
	c.echoClock()
}

type scalarClockSerializer struct {
	Ticks int `json:"ticks"`
}

func (c *ScalarClock) serialize() (string, error) {
	jsonClock, err := json.MarshalIndent(scalarClockSerializer{
		Ticks: c.ticks,
	}, "", "\t")
	return string(jsonClock), err
}

func (c *ScalarClock) GetClockStr() string {
	clockStr, err := c.serialize()
	if err != nil {
		fmt.Println("customerror serializing clock")
	}
	return clockStr
}

func (c *ScalarClock) echoClock() {
	fmt.Println("logical clock: ", c.ticks)
}

func (c *ScalarClock) parse(jsonClock string) (scalarClockSerializer, error) {
	var otherClock scalarClockSerializer
	err := json.Unmarshal([]byte(jsonClock), &otherClock)
	return otherClock, err
}

func (c *ScalarClock) CompareClocks(requestClockStr string, externalClockStr string, externalId int) (int, error) {
	requestClock, err := c.parse(requestClockStr)
	if err != nil {
		return 0, err
	}
	externalClock, err := c.parse(externalClockStr)
	if err != nil {
		return 0, err
	}
	if requestClock.Ticks < externalClock.Ticks {
		return c.id, nil
	} else if requestClock.Ticks > externalClock.Ticks {
		return externalId, nil
	}
	return math.Min(c.id, externalId), nil
}
