package main

import (
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/clock"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
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

	terminalCh := make(chan string)
	go src.ReadInput(terminalCh)

	serverCh := make(chan string)
	go network.Serve(serverCh, ports[myId-1])

	var logicalClock clock.LogicalClock
	logicalClock = clock.NewScalarClock(myId)
	state := src.Released
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
				if state != src.Released {
					fmt.Println("x ignored")
					break
				}
				state = src.Wanted
				logicalClock.InternalEvent()
				for id := 0; id < len(ports); id++ {
					if id+1 == myId {
						continue
					}
					network.UdpSend(connections[id+1], src.BuildRequestMessage(myId, logicalClock))
				}
			default:
				fmt.Println("invalid command, ignoring...")
			}
		case msg, valid := <-serverCh:
			if !valid {
				break
			}
			parsedMsg, err := src.ParseMessage(msg)
			if err != nil {
				fmt.Println("invalid message, ignoring...")
			}
			logicalClock.ExternalEvent(parsedMsg.ClockStr)

			switch src.MessageType(parsedMsg.Text) {
			case src.Request:
				network.UdpSend(connections[parsedMsg.SenderId], src.BuildReplyMessage(myId, logicalClock))
			case src.Reply:
			default:
			}
		default:
		}
		time.Sleep(time.Second * 1)
	}
}
