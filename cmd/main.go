package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/yene/knxproj"
)

const tmpFolder = "temp/"

var projects []knxproj.Project

func main() {
	flag.Parse()

	knxprojFile := flag.Arg(0)

	if !strings.HasSuffix(knxprojFile, ".knxproj") {
		log.Println("Provided file is not a knxproj.")
		return
	}

	err := knxproj.Unzip(knxprojFile, tmpFolder)
	if err != nil {
		log.Fatal("unzip error", err)
	}

	knxproj.Language = "en-US"
	knxproj.ParseManufacturers(tmpFolder)
	knxproj.ReadManufacturerData(tmpFolder)
	knxproj.ReadProjects(tmpFolder)

	os.RemoveAll(tmpFolder)
	/*
			for _, d := range knxproj.MainProject.DeviceList {
				log.Println(d.String())
			}

		for _, d := range knxproj.MainProject.GroupAddressList {
			log.Println(d.DataPointType(), d.Address, d.Name, d.LinkedDevices)
		}*/

	for _, d := range knxproj.AllGroupAddress {
		log.Println(d.DataPointType(), d.Address, d.Name, d.LinkedDevices)
	}
	/*
		rankingsJSON, _ := json.MarshalIndent(knxproj.MainProject, "", "  ")
		ioutil.WriteFile("project.json", rankingsJSON, 0644)
		fmt.Println("Written devices to project.json")
	*/

	var all JSONAddresses

	for _, d := range knxproj.AllGroupAddress {
		all = append(all, JSONAddress{Address: d.Address, Name: d.Name, DPT: d.DataPointType(), LinkedDevices: d.LinkedDevices})
		if d.DataPointType() == "" && d.LinkedDevices > 0 {
			log.Println("This GA should have a DPT:", d)
		}
	}

	sort.Sort(all)

	addressesJSON, _ := json.MarshalIndent(all, "", "  ")
	ioutil.WriteFile("project-addresses.json", addressesJSON, 0644)
	fmt.Println("Written devices to project-addresses.json")

}

type JSONAddress struct {
	Address       string `json:"address"`
	Name          string `json:"name"`
	DPT           string `json:"dpt"`
	LinkedDevices int    `json:"linkeddevices"`
}

type JSONAddresses []JSONAddress

func (s JSONAddresses) Len() int {
	return len(s)
}

func (s JSONAddresses) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// TODO: test with 2 style address
func (s JSONAddresses) Less(i, j int) bool {
	si := strings.Split(s[i].Address, "/")
	sj := strings.Split(s[j].Address, "/")
	si0, _ := strconv.Atoi(si[0])
	sj0, _ := strconv.Atoi(sj[0])
	si1, _ := strconv.Atoi(si[1])
	sj1, _ := strconv.Atoi(sj[1])
	if si0 < sj0 {
		return true
	} else if si0 == sj0 && si1 < sj1 {
		return true
	} else if len(si) == 3 {
		si2, _ := strconv.Atoi(si[2])
		sj2, _ := strconv.Atoi(sj[2])
		if si0 == sj0 && si1 == sj1 && si2 < sj2 {
			return true
		}
	}
	return false
}
