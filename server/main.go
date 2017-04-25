package main

import (
	"log"

	"github.com/wqtty/go-test/server/logic"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/wqtty/go-test/tserver"
	"github.com/wqtty/go-test/server/processor"
)

func main() {
	transportFactory := thrift.NewTBufferedTransportFactory(1024) //for node js default
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	//create socket
	serverSocket, err := thrift.NewTServerSocket("0.0.0.0:7777")
	if err != nil {
		log.Fatal("thrift.NewTServerSocket:", err)
	}
	//register handler
	multiplexedProcessor := thrift.NewTMultiplexedProcessor()
	//order
	handler := &logic.TServiceImpl{}
	orderProcessor := tserver.NewTserviceProcessor(handler)
	multiplexedProcessor.RegisterProcessor(logic.ServiceName, orderProcessor)
	logProcessor := &processor.LogProcessor{}
	logProcessor.RegisterProcessor(multiplexedProcessor)
	//create server
	server := thrift.NewTSimpleServer4(logProcessor, serverSocket, transportFactory, protocolFactory)
	err = server.Serve()
}
