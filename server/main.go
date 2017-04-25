package main

import (
	"log"

	"github.com/wqtty/thrift-monitor/server/logic"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/wqtty/thrift-monitor/server/interceptor"
	"github.com/wqtty/thrift-monitor/tserver"
)

const serverAddress = "0.0.0.0:7777"

func main() {
	transportFactory := thrift.NewTBufferedTransportFactory(1024) //for node js default
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	//create socket
	serverSocket, err := thrift.NewTServerSocket(serverAddress)
	if err != nil {
		log.Fatal("thrift.NewTServerSocket:", err)
	}
	//register handler
	multiplexedProcessor := thrift.NewTMultiplexedProcessor()
	handler := &logic.TServiceImpl{}
	tServiceProcessor := tserver.NewTserviceProcessor(handler)
	multiplexedProcessor.RegisterProcessor(logic.ServiceName, tServiceProcessor)

	//here is where the interceptor comes in
	logInterceptor := &interceptor.LogInterceptor{}
	logInterceptor.RegisterProcessor(multiplexedProcessor)
	//instead of using the multiplexedProcessor directly, we use our interceptor as the first argument
	server := thrift.NewTSimpleServer4(logInterceptor, serverSocket, transportFactory, protocolFactory)
	err = server.Serve()
}
