package main

import (
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/clock"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/messages"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/network"
	"net"
	"os"
	"strconv"
	"time"
)

func initConnections(myId int, ports []string) map[int]*net.UDPConn {
	connections := make(map[int]*net.UDPConn)
	for idx, port := range ports {
		if idx+1 == myId {
			continue
		}
		conn := network.UdpConnect(port)
		connections[idx+1] = conn
	}
	return connections
}

func closeConnections(connections map[int]*net.UDPConn) {
	for _, conn := range connections {
		err := conn.Close()
		customerror.CheckError(err)
	}
}

func main() {
	if len(os.Args) < 4 {
		customerror.CheckError(fmt.Errorf("not enough ports given as arguments"))
	}
	myId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		customerror.CheckError(fmt.Errorf("first argument should be a number representing the sequential process ID"))
	}
	ports := os.Args[2:len(os.Args)]

	connections := initConnections(myId, ports)
	defer closeConnections(connections)

	csConn := network.UdpConnect(src.SharedResourcePort)
	defer func() {
		err = csConn.Close()
		customerror.CheckError(err)
	}()

	terminalCh := make(chan string)
	go src.ReadInput(terminalCh)

	serverCh := make(chan string)
	go network.Serve(serverCh, ports[myId-1])

	var logicalClock clock.LogicalClock
	logicalClock = clock.NewScalarClock(myId)
	replyManager := src.NewReplyManager(len(ports))
	state := src.Released
	var requestTimestamp string
	for {
		select {
		case command, valid := <-terminalCh:
			if !valid {
				break
			}
			switch command {
			case strconv.Itoa(myId):
				logicalClock.InternalEvent()
			case src.ConsumeCmd:
				fmt.Println(state)
				switch state {
				case src.Released:
					state = src.Wanted
					logicalClock.InternalEvent()
					requestTimestamp = logicalClock.GetClockStr()
					for id := 0; id < len(ports); id++ {
						if id+1 == myId {
							continue
						}
						network.UdpSend(connections[id+1], messages.BuildRequestMessage(myId, logicalClock))
					}
				case src.Wanted:
					fmt.Println("x ignored")
				case src.Held:
					fmt.Println("x ignored")
				}
			default:
				fmt.Println("invalid command, ignoring...")
			}
		case msg, valid := <-serverCh:
			if !valid {
				break
			}
			parsedMsg, err := messages.ParseMessage(msg)
			if err != nil {
				fmt.Println("invalid message, ignoring...")
			}
			logicalClock.ExternalEvent(parsedMsg.ClockStr)

			switch messages.MessageType(parsedMsg.Text) {
			case messages.Request:
				switch state {
				case src.Released:
					network.UdpSend(connections[parsedMsg.SenderId], messages.BuildReplyMessage(myId, logicalClock))
				case src.Wanted:
					selectedId, err := logicalClock.CompareClocks(requestTimestamp, parsedMsg.ClockStr, parsedMsg.SenderId)
					if err != nil {
						fmt.Println("invalid message, ignoring...")
						break
					}
					if selectedId == myId {
						replyManager.EnqueueProcess(parsedMsg.SenderId)
					} else {
						network.UdpSend(connections[parsedMsg.SenderId], messages.BuildReplyMessage(myId, logicalClock))
					}
				case src.Held:
					replyManager.EnqueueProcess(parsedMsg.SenderId)
				}
			case messages.Reply:
				if !replyManager.ReceiveReply() {
					break
				}
				state = src.Held
				network.UdpSend(csConn, messages.BuildConsumeMessage(myId, logicalClock))
				fmt.Println("entered cs")
				time.Sleep(5 * time.Second)
				state = src.Released
				fmt.Println("left cs")
				processesToReply := replyManager.Dequeue()
				for _, id := range processesToReply {
					network.UdpSend(connections[id], messages.BuildReplyMessage(myId, logicalClock))
				}
				replyManager.Reset()
			case messages.Consume:
				fmt.Printf("received %s, but I'm not a shared resource, ignoring...\n", messages.Consume)
			default:
			}
		default:
		}
		time.Sleep(time.Second * 1)
	}
}
