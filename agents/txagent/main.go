package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// few global variables

var (
	doneDir       = "done"
	configuration configDef
	pkgTtl        int64
	cycleDuration time.Duration
	buffer        bufferDef
)

// Functions

func createDoneDirectory() {
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
	createDoneDirectory()
	pkgTtl = configuration.Ttl.Seconds + 60*configuration.Ttl.Minutes + 3600*configuration.Ttl.Hours
	cycleDuration = time.Duration(
		configuration.Cycle.MiliSec * int64(time.Second),
	)
	buffer.content = make(packagesByValidTill)
	// TODO: Load buffer from file
	go checkForNewFilesToTransmit()
	go serveStatus()
	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
	// TODO: Save buffer to file
}
