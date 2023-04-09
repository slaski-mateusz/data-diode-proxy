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
	createDoneDirectory(configuration.Files.DoneSubDir)
	pkgTtl = configuration.Ttl.Seconds + 60*configuration.Ttl.Minutes + 3600*configuration.Ttl.Hours
	cycleDuration = time.Duration(configuration.Cycle.MiliSec) * time.Millisecond

	procDelay = time.Duration(configuration.Files.ProcessAfterSec) * time.Second
	fmt.Printf("%+v, %+v\n", cycleDuration, procDelay)

	addrDestUDP = &net.UDPAddr{
		IP:   net.IP{127, 0, 0, 1},
		Port: 2345,
	}
	addrLocUDP := &net.UDPAddr{Port: 1234}
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
