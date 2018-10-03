package udpprotocol

import (
	"fmt"
	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/xabarass/sperf/protocols"
)

type UdpServer struct {
}

func (s *UdpServer) Run(config *protocols.ScionGenericConfig, serverAddr *snet.Addr) error {
	log.Debug("Starting UDP server")

	serverConnection, err := snet.ListenSCION("udp4", serverAddr)
	receivePacketBuffer := make([]byte, 2500)

	for {
		n, sc, err := serverConnection.ReadFromSCION(receivePacketBuffer)
		if err != nil {
			log.Error("Error receiving data")
			continue
		}
		if n < 1 {
			log.Error("Error receiving data")
			continue
		}
		recv := string(receivePacketBuffer[:n])
		fmt.Println(recv)
		fmt.Println(sc.String())
	}
	return err
}

func (s *UdpServer) Stop() error {
	return nil
}

func CreateServer() protocols.Server {
	return &UdpServer{}
}
