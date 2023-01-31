package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"
)

// few global variables

var (
	configuration configDef
	pkgTtl        int64
	cycleDuration time.Duration
)

// Configuration file types

type ttlDef struct {
	Seconds int64 `yaml:"seconds"`
	Minutes int64 `yaml:"minutes"`
	Hours   int64 `yaml:"hours"`
}

type cycleDef struct {
	Seconds int64 `yaml:"seconds"`
}

type packetSizeDef struct {
	Bytes int `yaml:"bytes"`
}

type configDef struct {
	Workdir     string        `yaml:"workdir"`
	FilePattern string        `yaml:"fpattern"`
	Cycle       cycleDef      `yaml:"cycle"`
	Ttl         ttlDef        `yaml:"ttl"`
	PacketSize  packetSizeDef `yaml:"packet_size"`
}

// Buffer types

type validTill int64

type packageId int64

type transmittedFileName string

type bufferDef struct {
	sync.Mutex
	content map[validTill]map[transmittedFileName]map[packageId][]byte
}

// Functions

func createDoneDirectory() {
	doneDirPath := path.Join(configuration.Workdir, "done")
	if _, errSt := os.Stat(doneDirPath); errors.Is(errSt, os.ErrNotExist) {
		errMk := os.Mkdir(doneDirPath, os.ModePerm)
		if errMk != nil {
			log.Fatalf(
				"Can not create 'done' directory in working directory %v due to problem: %v",
				configuration.Workdir,
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

func configFileName() string {
	executablePath, errExe := os.Executable()
	if errExe != nil {
		log.Fatalf("Problem with getting my executable name: %v", errExe)
	}
	_, executableName := path.Split(executablePath)
	return executableName
}

func putFileContentToBuffer(fileName string, fileToRead *os.File, buffer *bufferDef) {
	bytesRead := make([]byte, configuration.PacketSize.Bytes)
	var fileContent []byte
	buffer.Lock()
	defer buffer.Unlock()
	var pkgId packageId = 1
	for {
		var amountRead int
		amountRead, _ = fileToRead.Read(bytesRead)
		fmt.Print(amountRead, " ")
		vt := validTill(time.Now().Unix() + int64(pkgTtl))
		buffer.content[vt][transmittedFileName(fileName)][pkgId] = bytesRead
		pkgId++
		if amountRead == 0 {
			break
		}
	}
	fmt.Println(string(fileContent))
}

func removeOutdatedPackagesFromBuffer(buffer *bufferDef) {
	nowSeconds := time.Now().Unix()
	buffer.Lock()
	defer buffer.Unlock()
	for ts := range buffer.content {
		if ts < validTill(nowSeconds-pkgTtl) {
			delete(buffer.content, ts)
		}
	}
}

func checkForNewFilesToTransmit(buffer *bufferDef) {
	for {
		files, errRd := ioutil.ReadDir(configuration.Workdir)
		if errRd != nil {
			log.Fatalf(
				"Can not list working directory '%s'. Error: %v",
				configuration.Workdir,
				errRd,
			)
		}
		var timeToCheck time.Time
		timeToCheck = time.Now().Add(-cycleDuration)
		fmt.Println(time.Now())
		fmt.Println(timeToCheck)
		for _, f := range files {
			fileMatchPattern, _ := regexp.MatchString(configuration.FilePattern, f.Name())
			if !fileMatchPattern {
				continue
			}
			fileIsRegular := f.Mode().IsRegular()
			if !fileIsRegular {
				continue
			}
			fileIsRecentlyModified := f.ModTime().After(timeToCheck)
			if !fileIsRecentlyModified {
				continue
			}
			fmt.Println(f.Name(), f.ModTime(), " ")
			filePath := filepath.Join(configuration.Workdir, f.Name())
			fileToRead, errOpen := os.Open(filePath)
			if errOpen != nil {
				fmt.Printf("can not read file %v %v\n", f.Name(), errOpen)
			} else {
				putFileContentToBuffer(f.Name(), fileToRead, buffer)
			}

		}
		fmt.Println()
		time.Sleep(cycleDuration)
	}
}

// main

func main() {

	cancelChan := make(chan os.Signal, 1)
	signal.Notify(
		cancelChan,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	configuration.Load(configFileName() + ".yaml")
	createDoneDirectory()
	pkgTtl = configuration.Ttl.Seconds + 60*configuration.Ttl.Minutes + 3600*configuration.Ttl.Hours
	cycleDuration = time.Duration(
		configuration.Cycle.Seconds * int64(time.Second),
	)
	var buffer bufferDef
	// TODO: Load buffer from file
	go checkForNewFilesToTransmit(&buffer)
	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
	// TODO: Save buffer to file
}
