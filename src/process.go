package src

import (
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/clock"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/network"
	"net"
	"os"
	"strconv"
	"time"
)

func initConnections(myID int, ports []string) map[int]*net.UDPConn {
	connections := make(map[int]*net.UDPConn)
	for idx, port := range ports {
		if idx+1 == myID {
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

// Execute parses arguments into process information and
// executes accordingly
func Execute() {
	if len(os.Args) < 4 {
		customerror.CheckError(fmt.Errorf("not enough ports given as arguments"))
	}
	myID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		customerror.CheckError(fmt.Errorf("first argument should be a number representing the sequential process ID"))
	}
	ports := os.Args[2:len(os.Args)]

	connections := initConnections(myID, ports)
	defer closeConnections(connections)

	terminalCh := make(chan string)
	go readInput(terminalCh)

	serverCh := make(chan string)
	go network.Serve(serverCh, ports[myID-1])

	var logicalClock clock.LogicalClock
	logicalClock = clock.NewScalarClock(myID)
	state := Released
	for {
		select {
		case command, valid := <-terminalCh:
			if !valid {
				break
			}
			switch command {
			case strconv.Itoa(myID):
				logicalClock.InternalEvent()
			case ConsumeCmd:
				if state != Released {
					fmt.Println("x ignored")
					break
				}
				state = Wanted
				logicalClock.InternalEvent()
				for id := 0; id < len(ports); id++ {
					if id+1 == myID {
						continue
					}
					network.UdpSend(connections[id+1], buildRequestMessage(logicalClock))
				}
			default:
				fmt.Println("invalid command, ignoring...")
			}
		case msg, valid := <-serverCh:
			if !valid {
				break
			}
			parsedMsg, err := parseMessage(msg)
			if err != nil {
				fmt.Println("invalid message, ignoring...")
			}
			logicalClock.ExternalEvent(parsedMsg.ClockStr)

			switch MessageType(parsedMsg.Text) {
			case Request:
				id := logicalClock.GetProcessID(parsedMsg.ClockStr)
				network.UdpSend(connections[id], buildReplyMessage(logicalClock))
			case Reply:
			default:
			}
		default:
		}
		time.Sleep(time.Second * 1)
	}
}
