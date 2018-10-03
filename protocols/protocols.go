package protocols

import (
	"github.com/scionproto/scion/go/lib/snet"
)

type ScionGenericConfig struct {
	Interactive  *bool
	Sciond       *string
	Dispatcher   *string
	SciondFromIA *bool
}

type Server interface {
	Run(config *ScionGenericConfig, serverAddr *snet.Addr) error
	Stop() error
}

type Client interface {
	Run(config *ScionGenericConfig, serverAddr, clientAddr *snet.Addr) error
	Stop() error
}
