package main

// Configuration file types

type filesDef struct {
	WorkDir         string `yaml:"workdir"`
	DoneSubDir      string `yaml:"donesubdir"`
	Pattern         string `yaml:"fpattern"`
	ProcessAfterSec int64  `yaml:"process_after_sec"`
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
	Bytes int64 `yaml:"bytes"`
}

type networkDef struct {
	LocalPort int64  `yaml:"local_port"`
	DiodeIp   string `yaml:"diode_ip"`
	DiodePort string `yaml:"diode_port"`
}

type configDef struct {
	Files      filesDef      `yaml:"files"`
	Cycle      cycleDef      `yaml:"cycle"`
	Ttl        ttlDef        `yaml:"ttl"`
	PacketSize packetSizeDef `yaml:"packet_size"`
	Network    networkDef    `yaml:"network"`
}

// Buffer types

type validTill int64

type packageId int64

type transmittedFileName string

type packageData []byte

type packagesById map[packageId]packageData

type dataToTransmit struct {
	Vt             validTill
	Tfn            transmittedFileName
	PackagesNumber packageId
	Packages       packagesById
}

//transmit types

// type packageToTransmit struct {
// 	FileName       transmittedFileName
// 	Id             packageId
// 	ValidTill      validTill
// 	PackagesNumber int64
// 	Data           packageData
// }
