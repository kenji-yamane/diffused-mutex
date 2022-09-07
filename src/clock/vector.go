package clock

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/math"
)

type VectorClock struct {
	id    int
	ticks []int
}

func NewVectorClock(id, numProcesses int) LogicalClock {
	return &VectorClock{
		id:    id,
		ticks: make([]int, numProcesses),
	}
}

func (c *VectorClock) InternalEvent() {
	c.ticks[c.id-1]++
	c.echoClock()
}

func (c *VectorClock) ExternalEvent(externalClockStr string) {
	externalClock, err := parse(externalClockStr)
	if err != nil {
		fmt.Println("invalid clock string, ignoring...")
		return
	}
	c.ticks[c.id-1]++
	for idx, tick := range externalClock.Ticks {
		if idx+1 == c.id {
			continue
		}
		c.ticks[idx] = math.Max(c.ticks[idx], tick)
	}
	c.echoClock()
}

func (c *VectorClock) GetClockStr() string {
	clockStr, err := c.serialize()
	if err != nil {
		fmt.Println("customerror serializing clock")
	}
	return clockStr
}

func (c *VectorClock) echoClock() {
	fmt.Println("logical clock: ", c.ticks)
}

type vectorClockSerializer struct {
	Id    int   `json:"id"`
	Ticks []int `json:"ticks"`
}

func (c *VectorClock) serialize() (string, error) {
	jsonClock, err := json.Marshal(vectorClockSerializer{
		Id:    c.id,
		Ticks: c.ticks,
	})
	return string(jsonClock), err
}

func parse(jsonClock string) (vectorClockSerializer, error) {
	var otherClock vectorClockSerializer
	err := json.Unmarshal([]byte(jsonClock), &otherClock)
	return otherClock, err
}
