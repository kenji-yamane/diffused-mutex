package src

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/clock"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
)

type messageSerializer struct {
	SenderId int    `json:"id"`
	Text     string `json:"text"`
	ClockStr string `json:"clock_str"`
}

func buildReplyMessage(myId int, c clock.LogicalClock) string {
	return buildMessage(myId, c, Reply)
}

func buildRequestMessage(myId int, c clock.LogicalClock) string {
	return buildMessage(myId, c, Request)
}

func buildMessage(myId int, c clock.LogicalClock, messageType MessageType) string {
	m := messageSerializer{
		SenderId: myId,
		Text:     messageType.String(),
		ClockStr: c.GetClockStr(),
	}
	mStr, err := json.MarshalIndent(m, "", "\t")
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
