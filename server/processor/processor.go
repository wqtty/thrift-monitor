package processor

import (
	"log"

	"unsafe"

	"fmt"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type TMultiplexedProcessor struct {
	serviceProcessorMap map[string]thrift.TProcessor
	DefaultProcessor    thrift.TProcessor
}

type LogProcessor struct {
	realProcessor thrift.TProcessor
}

func (p *LogProcessor) RegisterProcessor(processor thrift.TProcessor) {
	p.realProcessor = processor
}

func (p *LogProcessor) Process(in, out thrift.TProtocol) (bool, thrift.TException) {
	multiplexedProcessor := (*TMultiplexedProcessor)(unsafe.Pointer(p.realProcessor.(*thrift.TMultiplexedProcessor)))
	name, typeId, seqid, err := in.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if typeId != thrift.CALL && typeId != thrift.ONEWAY {
		return false, fmt.Errorf("Unexpected message type %v", typeId)
	}
	//extract the service name
	v := strings.SplitN(name, thrift.MULTIPLEXED_SEPARATOR, 2)
	if len(v) != 2 {
		if multiplexedProcessor.DefaultProcessor != nil {
			smb := thrift.NewStoredMessageProtocol(in, name, typeId, seqid)
			return multiplexedProcessor.DefaultProcessor.Process(smb, out)
		}
		return false, fmt.Errorf("Service name not found in message name: %s.  Did you forget to use a TMultiplexProtocol in your client?", name)
	}
	actualProcessor, ok := multiplexedProcessor.serviceProcessorMap[v[0]]
	if !ok {
		return false, fmt.Errorf("Service name not found: %s.  Did you forget to call registerProcessor()?", v[0])
	}
	log.Print("Got req:", v[0], ".", v[1])
	smb := thrift.NewStoredMessageProtocol(in, v[1], typeId, seqid)
	rslt, e := actualProcessor.Process(smb, out)
	log.Print("intercepted by LogProcessor, after process")
	return rslt, e
}
