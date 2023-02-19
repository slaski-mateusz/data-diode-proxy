package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func putFileContentToBuffer(fileName string, fileToRead *os.File) {
	bytesRead := make([]byte, configuration.PacketSize.Bytes)
	tfn := transmittedFileName(fileName)
	vt := validTill(time.Now().Unix() + int64(pkgTtl))
	var pkgId packageId = 1
	var fileContent []byte
	buffer.Lock()
	defer buffer.Unlock()
	for {
		var amountRead int
		amountRead, _ = fileToRead.Read(bytesRead)
		fmt.Print(amountRead, " ")
		if _, ok := buffer.content[vt]; !ok {
			buffer.content[vt] = make(packagesByFile)
		}
		if _, ok := buffer.content[vt][tfn]; !ok {
			buffer.content[vt][tfn] = make(packagesById)
		}
		buffer.content[vt][tfn][pkgId] = bytesRead
		fileContent = append(fileContent, bytesRead...)
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

func checkForNewFilesToTransmit() {
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
				putFileContentToBuffer(f.Name(), fileToRead)
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
	for {
		// https://www.golinuxcloud.com/golang-udp-server-client/
		// https://holwech.github.io/blog/Creating-a-simple-UDP-module/
		time.Sleep(cycleDuration)
	}
}
