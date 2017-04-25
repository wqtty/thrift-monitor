package logic

import "log"

type TServiceImpl struct {
}

const ServiceName = "tservice"

func (*TServiceImpl)Test() error {
	log.Print("got Test request")
	return nil
}