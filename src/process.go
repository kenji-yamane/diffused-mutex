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
	logicalClock = clock.NewVectorClock(myID, len(ports))
	for {
		select {
		case command, valid := <-terminalCh:
			if !valid {
				break
			}
			id, err := strconv.Atoi(command)
			if err != nil {
				fmt.Println("invalid command, ignoring...")
				break
			}
			if id < 1 || id > len(ports) {
				fmt.Println("given id does not exist in context, ignoring...")
				break
			}
			logicalClock.InternalEvent()
			if id != myID {
				network.UdpSend(connections[id], logicalClock.GetClockStr())
			}
		case msg, valid := <-serverCh:
			if !valid {
				break
			}
			logicalClock.ExternalEvent(msg)
		default:
		}
		time.Sleep(time.Second * 1)
	}
}

