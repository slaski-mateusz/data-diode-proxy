package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

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
				// Processing only files that match configured pattern
				continue
			}
			fileIsRegular := f.Mode().IsRegular()
			if !fileIsRegular {
				// Not regular file as directory wouldn't be processed
				continue
			}
			fileModTime := f.ModTime()
			nowTime := time.Now()
			treshholdTime := nowTime.Add(-procDelay)
			if fileModTime.After(treshholdTime) {
				// Not time for processing file yet. Maybe it would processed.
				continue
			}
			fmt.Printf("File %v would be processed\n", f.Name())
			filePath := filepath.Join(
				configuration.Files.WorkDir,
				f.Name(),
			)
			doneFilePath := filepath.Join(
				configuration.Files.WorkDir,
				configuration.Files.DoneSubDir,
				f.Name(),
			)
			processFile(f.Name(), filePath, doneFilePath)
		}
		time.Sleep(cycleDuration)
	}
}

func processFile(fileName string, filePath string, doneFilePath string) {
	fileToRead, errOpen := os.Open(filePath)
	defer fileToRead.Close()
	if errOpen != nil {
		fmt.Printf("Can not read file %v %v\n", filePath, errOpen)
		return
	}
	dtt := new(dataToTransmit)
	dtt.Packages = make(packagesById)
	dtt.Vt = validTill(time.Now().Unix() + int64(pkgTtl))
	dtt.Tfn = transmittedFileName(fileName)
	bytesRead := make([]byte, configuration.PacketSize.Bytes)
	var pckId packageId = 0
	var fileContent []byte
	fmt.Println("Processing file:", fileName)
	for {
		// amountRead, errRead := fileToRead.Read(bytesRead)
		_, errRead := fileToRead.Read(bytesRead)
		// fmt.Println(amountRead, "bytes read")
		fmt.Println(bytesRead)
		dtt.Packages[pckId] = bytesRead
		fileContent = append(fileContent, bytesRead...)
		pckId++
		if errRead == io.EOF {
			dtt.PackagesNumber = pckId
			// fmt.Println("File content:")
			// fmt.Println(string(fileContent))
			errRen := os.Rename(filePath, doneFilePath)
			if errRen != nil {
				log.Fatalf(
					"Can not move file from '%v' to '%v'\n",
					filePath,
					doneFilePath,
				)
			} else {
				// We send dat only if file was successfully moved to "done" directory
				go sendFileData(dtt)
				// fmt.Println("After sending")
			}
			return
		}
	}
}

func sendFileData(dtt *dataToTransmit) {
	// fmt.Println("    Sending activated")
	for {
		now := validTill(time.Now().Unix())
		if now > dtt.Vt {
			fmt.Println(
				"Data for file",
				dtt.Tfn,
				"outdated. Stopping sending it",
			)
			// Data outdated not sending it
			return
		}
		// sending data simulation
		for id, pck := range dtt.Packages {
			fmt.Printf("File %v Package %+v data: %+v and string content %v\n", dtt.Tfn, id, pck, string(pck))
		}
		// Sending data
		// https://www.golinuxcloud.com/golang-udp-server-client/ ??
		// https://holwech.github.io/blog/Creating-a-simple-UDP-module/ ??
		time.Sleep(cycleDuration)
	}
}
