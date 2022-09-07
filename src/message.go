package src

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/clock"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
)

type messageSerializer struct {
	Text     string `json:"text"`
	ClockStr string `json:"clock_str"`
}

func buildRequestMessage(c clock.LogicalClock) string {
	return buildMessage(c, Request)
}

func buildMessage(c clock.LogicalClock, messageType MessageType) string {
	m := messageSerializer{
		Text:     messageType.String(),
		ClockStr: c.GetClockStr(),
	}
	mStr, err := json.Marshal(m)
	if err != nil {
		customerror.CheckError(fmt.Errorf("error serializing"))
	}
	return string(mStr)
}

func parseMessage(msg string) (messageSerializer, error) {
	var msgSerializer messageSerializer
	err := json.Unmarshal([]byte(msg), &msgSerializer)
	return msgSerializer, err
}
