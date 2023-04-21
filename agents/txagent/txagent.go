package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const (
	packetNumberBytes = 8
	packetAmountBytes = 8
	fileNamelenBytes  = 255
)

func dataBytesInPacket(packetSize int64) int64 {
	// returns amount o bytes that fit in package.
	// Send package size minus tranmission overhead:
	// - packet number int64 -> 8 bytes
	// - total packets number int 64 -> 8 bytes
	// - file name 255 bytes
	dataSize := packetSize - packetNumberBytes - packetAmountBytes - fileNamelenBytes
	if dataSize <= 0 {
		log.Fatal(
			fmt.Sprintf(
				"Calculated amout of data to put in packet is %v. Please correct configuration with packet_size.bytes > 271.",
				dataSize,
			),
		)
	}
	return dataSize
}

func int64ToBytes(value int64) []byte {
	bt := make([]byte, 8)
	binary.LittleEndian.PutUint64(bt, uint64(value))
	return bt
}

func bytesToInt64(bytes []byte) int64 {
	return int64(binary.LittleEndian.Uint64(bytes))
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
			//TODO: It is potential risk that if file would be big and cycle short than file would be get to proccess again
			//TODO: Solution would be to keep in memory list of files being processed
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
	bytesRead := make(
		[]byte,
		dataInPacket,
	)
	var pckId packageId = 0
	fmt.Println("Processing file:", fileName)
	for {
		// amountRead, errRead := fileToRead.Read(bytesRead)
		_, errRead := fileToRead.Read(bytesRead)
		// fmt.Println(amountRead, "bytes read")
		fmt.Println(bytesRead)
		dtt.Packages[pckId] = bytesRead
		pckId++
		if errRead == io.EOF {
			dtt.PackagesNumber = pckId
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

		for _, pck := range dtt.Packages {
			// Sending data - thanks to to Jakob Borg advice for klaymen
			// https://forum.golangbridge.org/t/sending-out-udp-packets-fast-without-context-connection/7672/9
			//TODO Pack in byte array:
			//TODO - package number
			//TODO - number of packages
			//TODO - file name
			//TODO - data
			_, errSend := connUDP.WriteTo(pck, addrDestUDP)
			if errSend != nil {
				fmt.Println("error sendig udp message:", errSend)
			}
		}

		time.Sleep(cycleDuration)
	}
}
