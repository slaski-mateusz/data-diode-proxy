package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

// few global variables

var (
	doneDir       = "done"
	configuration configDef
	pkgTtl        int64
	cycleDuration time.Duration
)

// Configuration file types

type filesDef struct {
	WorkDir         string `yaml:"workdir"`
	Pattern         string `yaml:"fpattern"`
	ProcessAfterSec int    `yaml:"process_after_sec"`
}

type ttlDef struct {
	Seconds int64 `yaml:"seconds"`
	Minutes int64 `yaml:"minutes"`
	Hours   int64 `yaml:"hours"`
}

type cycleDef struct {
	MiliSec int64 `yaml:"milisec"`
}

type packetSizeDef struct {
	Bytes int `yaml:"bytes"`
}

type configDef struct {
	Files      filesDef      `yaml:"files"`
	Cycle      cycleDef      `yaml:"cycle"`
	Ttl        ttlDef        `yaml:"ttl"`
	PacketSize packetSizeDef `yaml:"packet_size"`
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
		files, errRd := ioutil.ReadDir(configuration.Files.WorkDir)
		if errRd != nil {
			log.Fatalf(
				"Can not list working directory '%s'. Error: %v",
				configuration.Files.WorkDir,
				errRd,
			)
		}
		for _, f := range files {
			fileMatchPattern, _ := regexp.MatchString(configuration.Files.Pattern, f.Name())
			if !fileMatchPattern {
				continue
			}
			fileIsRegular := f.Mode().IsRegular()
			if !fileIsRegular {
				continue
			}
			filePath := filepath.Join(
				configuration.Files.WorkDir,
				f.Name(),
			)
			doneFilePath := filepath.Join(
				configuration.Files.WorkDir,
				doneDir,
				f.Name(),
			)
			fileToRead, errOpen := os.Open(filePath)
			if errOpen != nil {
				fmt.Printf("Can not read file %v %v\n", f.Name(), errOpen)
			} else {
				putFileContentToBuffer(f.Name(), fileToRead, buffer)
				errRen := os.Rename(filePath, doneFilePath)
				if errRen != nil {
					log.Fatalf(
						"Can not move file from '%v' to '%v'\n",
						filePath,
						doneFilePath,
					)
				}
			}

		}
		fmt.Println()
		time.Sleep(cycleDuration)
	}
}

func transmitBufferData() {

}

func renderBuffer(resWri http.ResponseWriter) {

}

func rencerConfig(resWri http.ResponseWriter) {

}

func serveBufferContent() {
	var router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("buffer", renderBuffer)
	router.HandleFunc("config", renderConfig)
	http.ListenAndServe("127.0.0.1:8686", router)
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
		configuration.Cycle.MiliSec * int64(time.Second),
	)
	var buffer bufferDef
	// TODO: Load buffer from file
	go checkForNewFilesToTransmit(&buffer)

	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
	// TODO: Save buffer to file
}
