package main

import (
	"time"

	"log"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/wqtty/go-test/tserver"
)

const serverAddress = "127.0.0.1:7777"

func main() {
	clientSocket, err := thrift.NewTSocket(serverAddress)
	if err != nil {
		log.Fatal("CAN NOT CONNECT TO PAYMENT SERVICE! YOU CAN NOT CREATE ORDER!", err)
		return
	}
	clientSocket.SetTimeout(time.Duration(5) * time.Second)
	transport := thrift.NewTBufferedTransport(clientSocket, 1024)
	protocol := thrift.NewTBinaryProtocolTransport(transport)
	mp := thrift.NewTMultiplexedProtocol(protocol, "tservice")
	if err := transport.Open(); err != nil {
		log.Fatal("transport.Open failed, trying to connect to ", serverAddress, " err:", err)
		return
	}
	client := tserver.NewTserviceClientProtocol(transport, mp, mp)
	client.Test()
}
