package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"

	"github.com/scionproto/scion/go/lib/log"
	"github.com/scionproto/scion/go/lib/sciond"
	"github.com/scionproto/scion/go/lib/snet"
	"github.com/xabarass/sperf/protocols"
	udpPerf "github.com/xabarass/sperf/protocols/udp"
)

type SperfConfig struct {
	serverAddr  *string
	clientAddr  *string
	useUdp      *bool
	scionConfig protocols.ScionGenericConfig
}

func loadConfiguration() (SperfConfig, error) {
	var config SperfConfig

	parser := argparse.NewParser("sperf", "Measure SCION network performance")
	// Server address always needs to be specified. If client address is speicified, that menas we are running client mode
	config.serverAddr = parser.String("s", "server", &argparse.Options{Required: true, Help: "Server SCION address"})
	config.clientAddr = parser.String("c", "client", &argparse.Options{Required: false, Help: "Client IP address"})
	// Specify which kind of traffic to use
	config.useUdp = parser.Flag("u", "udp", &argparse.Options{Required: false, Help: "Use raw UDP packets instead of QUIC"})
	// SCION specific flags
	config.scionConfig.Interactive = parser.Flag("i", "interactive", &argparse.Options{Required: false, Help: "Use interactive mode in choosing paths"})
	config.scionConfig.Sciond = parser.String("", "sciond", &argparse.Options{Required: false, Help: "Path to sciond socket"})
	config.scionConfig.Dispatcher = parser.String("", "dispatcher", &argparse.Options{Required: false, Help: "Path to dispatcher socket"})
	config.scionConfig.SciondFromIA = parser.Flag("", "sciondFromIA", &argparse.Options{Required: false, Help: "SCIOND socket path from IA address"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	return config, err
}

func runServer(server *protocols.Server, sc *SperfConfig) error {
	return nil
}

func initNetwork(config *protocols.ScionGenericConfig, localAddr *snet.Addr) error {
	var sciondPath string
	if *config.SciondFromIA {
		sciondPath = sciond.GetDefaultSCIONDPath(&localAddr.IA)
	} else {

		sciondPath = sciond.GetDefaultSCIONDPath(nil)
	}

	return snet.Init(localAddr.IA, sciondPath, *config.Dispatcher)
}

func checkError(what string, err error) {
	if err != nil {
		log.Crit(err.Error())
		os.Exit(1)
	}
}

func main() {
	config, err := loadConfiguration()
	if err != nil {
		os.Exit(1)
	}

	log.AddLogConsFlags()
	defer log.LogPanicAndExit()

	var localAddr *snet.Addr
	var runServer bool
	if *config.clientAddr != "" {
		localAddr, err = snet.AddrFromString(*config.clientAddr)
		checkError("Parsing client address", err)
		runServer = false
	} else {
		localAddr, err = snet.AddrFromString(*config.serverAddr)
		checkError("Parsing server address", err)
		runServer = true
	}

	err = initNetwork(&config.scionConfig, localAddr)
	checkError("Initializing network", err)

	if runServer {
		server := udpPerf.CreateServer()
		server.Run(&config.scionConfig, localAddr)

	} else {
		log.Debug("Not running the server")
	}
}
