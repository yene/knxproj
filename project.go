package knxproj

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var groupaddressList []GroupAddress

type Project struct {
	Name               string              `json:"name"`
	ID                 string              `json:"-"`
	Filename           string              `json:"filename"`
	GroupAddressStyle  int                 `json:"gaformat"`
	LastModified       time.Time           `json:"lastmodified"`
	DeviceList         Devices             `json:"devices"`
	GroupAddressList   GroupAddresses      `json:"-"`
	GroupAddressGroups []GroupAddressGroup `json:"-"`
}

type Devices []Device

func (s Devices) Len() int {
	return len(s)
}
func (s Devices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Devices) Less(i, j int) bool {
	si := strings.Split(s[i].Address, ".")
	sj := strings.Split(s[j].Address, ".")
	si0, _ := strconv.Atoi(si[0])
	sj0, _ := strconv.Atoi(sj[0])
	si1, _ := strconv.Atoi(si[1])
	sj1, _ := strconv.Atoi(sj[1])
	si2, _ := strconv.Atoi(si[2])
	sj2, _ := strconv.Atoi(sj[2])

	if si0 < sj0 {
		return true
	} else if si0 == sj0 && si1 < sj1 {
		return true
	} else if si0 == sj0 && si1 == sj1 && si2 < sj2 {
		return true
	}
	return false
}

type GroupAddresses []GroupAddress

func (s GroupAddresses) Len() int {
	return len(s)
}
func (s GroupAddresses) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// TODO: test with 2 style address
func (s GroupAddresses) Less(i, j int) bool {
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

func (p Project) String() string {
	return fmt.Sprintf("%s - %d - %v - %s", p.Name, p.GroupAddressStyle, p.LastModified, p.ID)
}

func NewProject(projFolder string) *Project {
	// TODO: what is the best way to attach data here?
	// callback, global state, channels, passing pointer, struct with methods
	p := new(Project)
	readProjectInformation(projFolder, p)
	readProjectData(projFolder, p)
	return p
}

/*
 * Read project information from project.xml/Project.xml
 */

type XMLProjectInformation struct {
	CreatedBy   string `xml:"CreatedBy,attr"`
	ToolVersion string `xml:"ToolVersion,attr"`
	Project     struct {
		ID                 string             `xml:"Id,attr"`
		ProjectInformation ProjectInformation `xml:"ProjectInformation"`
	}
}

type ProjectInformation struct {
	Comment           string `xml:"Comment,attr"`
	CodePage          string `xml:"CodePage,attr"`
	DeviceCount       string `xml:"DeviceCount,attr"`
	LastUsedPuid      string `xml:"LastUsedPuid,attr"`
	GUID              string `xml:"Guid,attr"`
	Name              string `xml:"Name,attr"`
	GroupAddressStyle string `xml:"GroupAddressStyle,attr"`
	LastModified      string `xml:"LastModified,attr"`
}

func readProjectInformation(projFolder string, p *Project) {
	projectfile := "Project.xml"
	if _, err := os.Stat(projFolder + projectfile); err != nil {
		// projectfile can be upper or lowercase
		projectfile = strings.ToLower(projectfile)
	}

	xmlFile, err := ioutil.ReadFile(projFolder + projectfile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var q XMLProjectInformation
	err = xml.Unmarshal(xmlFile, &q)

	if err != nil {
		log.Fatal("Error:", err)
	}

	p.ID = q.Project.ID
	p.Name = q.Project.ProjectInformation.Name

	if q.CreatedBy != "ETS5" {
		log.Fatal("Project was not created with ETS 5. Please upgrade to ETS 5.")
	}

	if q.Project.ProjectInformation.GroupAddressStyle == "ThreeLevel" {
		p.GroupAddressStyle = 3
	} else if q.Project.ProjectInformation.GroupAddressStyle == "TwoLevel" {
		p.GroupAddressStyle = 2
	} else {
		log.Fatalln("The GroupAddress Style of the project is not supported.", q.Project.ProjectInformation.GroupAddressStyle)
	}

	AddressStyle = p.GroupAddressStyle

	t, err := time.Parse("2006-01-02T15:04:05.0000000Z", q.Project.ProjectInformation.LastModified)
	p.LastModified = t
	if err != nil {
		log.Println(err)
	}
}

/*
 * Read project data from 0.xml
 */

type XMLProjectInformationData struct {
	Installations []Installation `xml:"Project>Installations>Installation"`
}

type Installation struct {
	Name           string       `xml:"Name,attr"`
	BCUKey         string       `xml:"BCUKey,attr"`
	DefaultLine    string       `xml:"DefaultLine,attr"`
	Topology       []Area       `xml:"Topology>Area"`
	GroupAddresses []GroupRange `xml:"GroupAddresses>GroupRanges>GroupRange"`
}

func readProjectData(projFolder string, p *Project) {
	const datafile = "0.xml"
	xmlFile, err := ioutil.ReadFile(projFolder + datafile)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}

	var q XMLProjectInformationData
	err = xml.Unmarshal(xmlFile, &q)

	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	if len(q.Installations) > 1 {
		log.Println("Notice: Found more than one Installation.")
	}

	for _, e := range q.Installations {
		addrs, groups := readGroupAdresses(e.GroupAddresses, p.GroupAddressStyle)
		p.GroupAddressList = append(p.GroupAddressList, addrs...)
		p.GroupAddressGroups = append(p.GroupAddressGroups, groups...)
		p.DeviceList = append(p.DeviceList, readDevices(e.Topology)...)
	}

	sort.Sort(p.DeviceList)
	sort.Sort(p.GroupAddressList)

}

type Area struct {
	ID      string `xml:"Id,attr"`
	Address string `xml:"Address,attr"`
	Name    string `xml:"Name,attr"`
	Puid    string `xml:"Puid,attr"`
	Lines   []Line `xml:"Line"`
}

type Line struct {
	ID              string           `xml:"Id,attr"`
	Address         string           `xml:"Address,attr"`
	Name            string           `xml:"Name,attr"`
	MediumTypeRefID string           `xml:"MediumTypeRefId,attr"`
	Puid            string           `xml:"Puid,attr"`
	Devices         []DeviceInstance `xml:"DeviceInstance"`
}

type ChannelInstance struct {
	IsActive string `xml:"IsActive,attr"`
	RefID    string `xml:"RefId,attr"`
	ID       string `xml:"Id,attr"`
	Name     string `xml:"Name,attr"`
}

type DeviceInstance struct {
	Description                               string                 `xml:"Description,attr"`
	MediumConfigLoaded                        string                 `xml:"MediumConfigLoaded,attr"`
	Comment                                   string                 `xml:"Comment,attr"`
	LastModified                              time.Time              `xml:"LastModified,attr"`
	Hardware2ProgramRefID                     string                 `xml:"Hardware2ProgramRefId,attr"`
	ID                                        string                 `xml:"Id,attr"`
	SerialNumber                              string                 `xml:"SerialNumber,attr"`
	IndividualAddressLoaded                   string                 `xml:"IndividualAddressLoaded,attr"`
	ProductRefID                              string                 `xml:"ProductRefId,attr"`
	ApplicationProgramLoaded                  string                 `xml:"ApplicationProgramLoaded,attr"`
	ParametersLoaded                          string                 `xml:"ParametersLoaded,attr"`
	IsCommunicationObjectVisibilityCalculated string                 `xml:"IsCommunicationObjectVisibilityCalculated,attr"`
	LastDownload                              time.Time              `xml:"LastDownload,attr"`
	CommunicationPartLoaded                   string                 `xml:"CommunicationPartLoaded,attr"`
	Name                                      string                 `xml:"Name,attr"`
	Address                                   string                 `xml:"Address,attr"`
	Puid                                      string                 `xml:"Puid,attr"`
	ComObjectInstanceRefs                     []ComObjectInstanceRef `xml:"ComObjectInstanceRefs>ComObjectInstanceRef"`
	ChannelInstances                          []ChannelInstance      `xml:"ChannelInstances>ChannelInstance"`
}

/*
 * Topology
 *   Area
 *     Line
 *       DeviceInstance
 */
func readDevices(topology []Area) []Device {
	var devices []Device
	for _, area := range topology {
		for _, line := range area.Lines {
			for _, device := range line.Devices {
				// Skip not configured devices.
				if device.Address == "" {
					continue
				}

				d := NewDevice(area.Address, line.Address, device)
				devices = append(devices, *d)
			}
		}
	}
	return devices
}

type GroupRange struct {
	Name         string          `xml:"Name,attr"`
	Puid         int             `xml:"Puid,attr"` // TODO: what does it mean?
	ID           string          `xml:"Id,attr"`
	RangeStart   int             `xml:"RangeStart,attr"` // TODO: what does it mean?
	RangeEnd     int             `xml:"RangeEnd,attr"`
	GroupRange   []GroupRange    `xml:"GroupRange"`
	GroupAddress []GroupAddress2 `xml:"GroupAddress"`
}

type GroupAddress2 struct {
	Puid          int    `xml:"Puid,attr"`
	ID            string `xml:"Id,attr"`
	Address       int    `xml:"Address,attr"`
	Name          string `xml:"Name,attr"`
	DatapointType string `xml:"DatapointType,attr"`
}

type GroupAddressGroup struct {
	Name    string
	Address string
}

/*
 * GroupAddresses
 *   GroupRanges
 *     GroupRange       <-- Main
 *       GroupRange         <-- Sub
 *         GroupAddress         <-- GA
 */
// TODO: maybe this should use a stack too, until it finds the address
func readGroupAdresses(ranges []GroupRange, style int) ([]GroupAddress, []GroupAddressGroup) {
	var gaddrs []GroupAddress
	var groups []GroupAddressGroup

	if style == 3 {
		for _, main := range ranges {
			if debug {
				log.Println(main.Name, main.RangeStart>>11&0x1F)
			}
			address := fmt.Sprintf("%d", main.RangeStart>>11&0x1F)
			groups = append(groups, GroupAddressGroup{main.Name, address})
			for _, sub := range main.GroupRange {
				if debug {
					log.Println(" ", sub.Name, sub.RangeStart>>11&0x1F, sub.RangeStart>>8&0x07)
				}
				address := fmt.Sprintf("%d/%d", sub.RangeStart>>11&0x1F, sub.RangeStart>>8&0x07)
				groups = append(groups, GroupAddressGroup{sub.Name, address})
				for _, addr := range sub.GroupAddress {
					ga := NewGroupAddress(addr)
					AllGroupAddress[ga.ID] = *ga
					if debug {
						log.Println("    ", ga.Name, ga.Address)
					}
					gaddrs = append(gaddrs, *ga)
				}
			}
		}
	} else {
		for _, main := range ranges {
			if debug {
				log.Println(main.Name, main.RangeStart>>11&0x1F)
			}
			address := fmt.Sprintf("%d", main.RangeStart>>11&0x1F)
			groups = append(groups, GroupAddressGroup{main.Name, address})
			for _, addr := range main.GroupAddress {
				ga := NewGroupAddress(addr)
				AllGroupAddress[ga.ID] = *ga
				if debug {
					log.Println(" ", ga.Name, ga.Address)
				}
				gaddrs = append(gaddrs, *ga)
			}

		}
	}
	return gaddrs, groups
}
