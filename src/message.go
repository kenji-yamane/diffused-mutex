package src

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/clock"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
)

type MessageSerializer struct {
	SenderId int    `json:"id"`
	Text     string `json:"text"`
	ClockStr string `json:"clock_str"`
}

func BuildReplyMessage(myId int, c clock.LogicalClock) string {
	return buildMessage(myId, c, Reply)
}

func BuildRequestMessage(myId int, c clock.LogicalClock) string {
	return buildMessage(myId, c, Request)
}

func buildMessage(myId int, c clock.LogicalClock, messageType MessageType) string {
	m := MessageSerializer{
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

func ParseMessage(msg string) (MessageSerializer, error) {
	var msgSerializer MessageSerializer
	err := json.Unmarshal([]byte(msg), &msgSerializer)
	return msgSerializer, err
}
