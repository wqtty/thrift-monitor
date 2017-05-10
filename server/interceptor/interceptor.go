package interceptor

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unsafe"

	"git.apache.org/thrift.git/lib/go/thrift"
)

//TMultiplexedProcessor is copied from thrift v0.9.3 git.apache.org/thrift.git/lib/go/thrift/multiplexed_protocol.go
type TMultiplexedProcessor struct {
	serviceProcessorMap map[string]thrift.TProcessor
	DefaultProcessor    thrift.TProcessor
}

//LogInterceptor is the interceptor that logs all RPCs invoked by thrift multiplexed processor
//it should implements TProcessor interface defined in git.apache.org/thrift.git/lib/go/thrift/processor.go
type LogInterceptor struct {
	realProcessor *thrift.TMultiplexedProcessor
}

//RegisterProcessor registers the multiplexed processor for later use
func (p *LogInterceptor) RegisterProcessor(processor *thrift.TMultiplexedProcessor) {
	p.realProcessor = processor
}

//Process satisfies TProcessor interface, the code below is copied from git.apache.org/thrift.git/lib/go/thrift/multiplexed_protocol.go
//and slightly changed, i.e. add logging
func (p *LogInterceptor) Process(in, out thrift.TProtocol) (bool, thrift.TException) {
	multiplexedProcessor := (*TMultiplexedProcessor)(unsafe.Pointer(p.realProcessor))
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
	start := time.Now()
	rslt, e := actualProcessor.Process(smb, out)
	elapsed := time.Since(start)
	log.Print("method:", v[0], ".", v[1], " has taken ", elapsed.Nanoseconds()/time.Millisecond.Nanoseconds(), " ms")
	return rslt, e
}
