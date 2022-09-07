package network

import (
	"fmt"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/customerror"
	"net"
)

func UdpConnect(port string) *net.UDPConn {
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+port)
	customerror.CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	customerror.CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	customerror.CheckError(err)

	return conn
}

func UdpSend(conn *net.UDPConn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(msg, err)
	} else {
		fmt.Println("sending message: ", msg)
	}
}
