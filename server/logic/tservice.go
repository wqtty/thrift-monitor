package logic

import (
	"log"
	"time"
)

type TServiceImpl struct {
}

const ServiceName = "tservice"

func (*TServiceImpl) Test() error {
	log.Print("got Test request")
	time.Sleep(10 * time.Millisecond)
	return nil
}
