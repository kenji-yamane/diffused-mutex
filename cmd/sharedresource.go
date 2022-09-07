package main

import (
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src"
	"github.com/kenji-yamane/distributed-mutual-exclusion-sample/src/network"
)

func main() {
	serverCh := make(chan string)
	network.Serve(serverCh, src.SharedResourcePort)
}
