// Package knxproj parses the ETS file format and returns Groupaddresses with DPT.
package knxproj

import (
	"io/ioutil"
	"regexp"
)

const debug = false

var Language = "en-US"
var MainProject *Project

var Hardware map[string]Product
var ComObjects map[string]ComObject
var ComObjectRefs map[string]ComObjectRef
var Manufacturers map[string]Manufacturer
var AllGroupAddress map[string]GroupAddress
var AddressStyle = 3

func ReadProjects(tmpFolder string) {
	AllGroupAddress = make(map[string]GroupAddress)

	// form project folder, list files, and filter the important files
	files, _ := ioutil.ReadDir(tmpFolder)
	const projectPattern = "P-[0-9A-F]{4}"
	r := regexp.MustCompile(projectPattern)

	// TODO: can a project contain multiple project folders?
	for _, f := range files {
		if f.IsDir() && r.MatchString(f.Name()) {
			MainProject = NewProject(tmpFolder + f.Name() + "/")
			return
		}
	}
}

func ReadManufacturerData(tmpFolder string) {
	ComObjects = make(map[string]ComObject)
	ComObjectRefs = make(map[string]ComObjectRef)
	Hardware = make(map[string]Product)

	const manufacturerPattern = "M-[0-9A-F]{4}"
	r := regexp.MustCompile(manufacturerPattern)

	// M-0083_A-004D-12-E268
	// M-00C8_A-2820-40-090B-O00C5
	const manufacturerDevicePattern = "M-[0-9A-F]{4}_[0-9A-F]-[0-9A-F]{4}-[0-9A-F]{2}-[0-9A-F]{4}(-[0-9A-Z]{1}[0-9A-F]{4})?"
	mr := regexp.MustCompile(manufacturerDevicePattern)

	files, _ := ioutil.ReadDir(tmpFolder)
	for _, f := range files {
		if f.IsDir() && r.MatchString(f.Name()) {
			parseHardwareData(tmpFolder + f.Name() + "/Hardware.xml")
			mfiles, _ := ioutil.ReadDir(tmpFolder + f.Name())
			for _, mf := range mfiles {
				if mr.MatchString(mf.Name()) {
					path := tmpFolder + f.Name() + "/" + mf.Name()
					parseManufacturerData(path)
				}
			}
		}
	}
}
