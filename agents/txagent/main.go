package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// few global variables

var (
	configuration configDef
	pkgTtl        int64
	cycleDuration time.Duration
	procDelay     time.Duration
	connUDP       *net.UDPConn
	addrDestUDP   *net.UDPAddr
)

// Functions

func createDoneDirectory(doneDir string) {
	doneDirPath := path.Join(configuration.Files.WorkDir, doneDir)
	if _, errSt := os.Stat(doneDirPath); errors.Is(errSt, os.ErrNotExist) {
		errMk := os.Mkdir(doneDirPath, os.ModePerm)
		if errMk != nil {
			log.Fatalf(
				"Can not create 'done' directory in working directory %v due to problem: %v",
				configuration.Files.WorkDir,
				errMk,
			)
		}
	}
}

func (c *configDef) Load(confFileName string) {
	confFile, errRead := ioutil.ReadFile(confFileName)
	if errRead != nil {
		log.Fatalf("Problem with reading config file '%s': %v", confFileName, errRead)
	}
	errYaml := yaml.Unmarshal(confFile, c)
	if errYaml != nil {
		log.Fatalf("Problem with config yaml un-marshaling: %v", errYaml)
	}
}

func executableName() string {
	executablePath, errExe := os.Executable()
	if errExe != nil {
		log.Fatalf("Problem with getting my executable name: %v", errExe)
	}
	_, exeName := path.Split(executablePath)
	return exeName
}

// main

func main() {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(
		cancelChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	configuration.Load(executableName() + ".yaml")
	fmt.Printf("%+v\n", configuration)
	createDoneDirectory(configuration.Files.DoneSubDir)
	pkgTtl = configuration.Ttl.Seconds + 60*configuration.Ttl.Minutes + 3600*configuration.Ttl.Hours
	cycleDuration = time.Duration(configuration.Cycle.MiliSec) * time.Millisecond

	procDelay = time.Duration(configuration.Files.ProcessAfterSec) * time.Second
	fmt.Printf("%+v, %+v\n", cycleDuration, procDelay)

	diodeIp := net.ParseIP(configuration.Network.DiodeIp)
	if diodeIp == nil {
		// IP value from config can not be parsed as IP. Assuming that given string is env variable
		diodeIpEnv := os.Getenv(
			configuration.Network.DiodeIp,
		)
		if diodeIpEnv == "" {
			log.Fatal(
				fmt.Sprintf(
					"Configuration contains diode IP as '%s' value. Such environment variable is not found or not set.",
					configuration.Network.DiodeIp,
				),
			)
		}
		diodeIp = net.ParseIP(diodeIpEnv)
		if diodeIp == nil {
			log.Fatal(
				fmt.Sprintf(
					"Configuration contains diode IP as '%s' value assummed to be environment variable, but value '%s' of this variable can not to be parsed to IP address",
					configuration.Network.DiodeIp,
					diodeIpEnv,
				),
			)
		}
	}

	diodePort, portErr := strconv.Atoi(configuration.Network.DiodePort)
	if portErr != nil {
		// Port value from config is not integer. Assuming that given string is env variable
		diodePortEnv := os.Getenv(
			configuration.Network.DiodePort,
		)
		if diodePortEnv == "" {
			log.Fatal(
				fmt.Sprintf(
					"Configuration contains diode port as '%s' value. Such environment variable is not found or not set.",
					configuration.Network.DiodePort,
				),
			)
		}
		diodePort, portErr = strconv.Atoi(diodePortEnv)
		if portErr != nil {
			log.Fatal(
				fmt.Sprintf(
					"Configuration contains diode port as '%s' value assummed to be environment variable, but value '%s' of this variable can not to be parsed to integer port number",
					configuration.Network.DiodePort,
					diodePortEnv,
				),
			)
		}
	}

	addrDestUDP = &net.UDPAddr{
		IP:   diodeIp,
		Port: diodePort,
	}

	addrLocUDP := &net.UDPAddr{Port: int(configuration.Network.LocalPort)}

	var errCon error
	connUDP, errCon = net.ListenUDP("udp", addrLocUDP)
	if errCon != nil {
		log.Fatal("UDP connection error:", errCon)
	}
	go checkForNewFilesToTransmit()
	go serveStatus()
	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
}
