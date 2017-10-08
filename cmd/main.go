package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	*/
	for _, d := range knxproj.MainProject.GroupAddressList {
		log.Println(d.DataPointType(), d.Address, d.Name)
	}

	rankingsJSON, _ := json.MarshalIndent(knxproj.MainProject, "", "  ")
	ioutil.WriteFile("project.json", rankingsJSON, 0644)
	fmt.Println("Written devices to project.json")
}
