package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"time"

	"gopkg.in/yaml.v3"
)

type ttlDef struct {
	Seconds int `yaml:"seconds"`
	Minutes int `yaml:"minutes"`
	Hours   int `yaml:"hours"`
}

type cycleDef struct {
	Seconds int `yaml:"seconds"`
}

type configDef struct {
	Workdir     string   `yaml:"workdir"`
	FilePattern string   `yaml:"fpattern"`
	Cycle       cycleDef `yaml:"cycle"`
	Ttl         ttlDef   `yaml:"ttl"`
}

func (c *configDef) Load(confFileName string) {
	confFile, errRead := ioutil.ReadFile(confFileName)
	fmt.Println(errRead)
	if errRead != nil {
		log.Fatalf("Problem with reading config file '%s': %v", confFileName, errRead)
	}
	fmt.Print(string(confFile))
	errYaml := yaml.Unmarshal(confFile, c)
	fmt.Println(errYaml)
	fmt.Println(&c)
	if errYaml != nil {
		log.Fatalf("Problem with config yaml un-marshaling: %v", errYaml)
	}
}

func main() {
	fmt.Println(os.Executable())
	var cf configDef
	cf.Load("txagent.yaml")
	fmt.Printf("%+v\n", cf)
	pkgTtl := cf.Ttl.Seconds + 60*cf.Ttl.Minutes + 3600*cf.Ttl.Hours
	cycleDuration := time.Duration(
		cf.Cycle.Seconds * int(time.Second),
	)
	fmt.Println(cf.Workdir, pkgTtl)
	for {
		files, errRd := ioutil.ReadDir(cf.Workdir)
		if errRd != nil {
			log.Fatalf(
				"Can not list working directory '%s'. Error: %v",
				cf.Workdir,
				errRd,
			)
		}
		var timeToCheck time.Time
		timeToCheck = time.Now().Add(-cycleDuration)
		fmt.Println(time.Now())
		fmt.Println(timeToCheck)
		for _, f := range files {
			matchPattern, _ := regexp.MatchString(
				cf.FilePattern,
				f.Name(),
			)
			if f.Mode().IsRegular() {
				if matchPattern {
					if f.ModTime().After(timeToCheck) {
						fmt.Print(f.Name(), f.ModTime(), " ")
					}
				}
			}
		}
		fmt.Println()
		time.Sleep(cycleDuration)
	}
}
