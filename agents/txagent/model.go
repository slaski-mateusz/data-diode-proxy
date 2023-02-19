package main

import (
	"sync"
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

type packageData []byte

type packagesById map[packageId]packageData

type packagesByFile map[transmittedFileName]packagesById

type packagesByValidTill map[validTill]packagesByFile

type bufferDef struct {
	sync.Mutex
	content packagesByValidTill
}

//transmit types

type packageToTransmit struct {
	FileName       transmittedFileName
	Id             packageId
	ValidTill      validTill
	PackagesNumber int64
	Data           packageData
}

// TODO: follow https://stackoverflow.com/questions/50698689/managing-slices-with-mutex-for-performance-in-golang
