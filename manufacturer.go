package knxproj

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

type ManufacturerData map[string]string

func parseManufacturerData(path string) ManufacturerData {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var q ManufacturerData2
	err = xml.Unmarshal(xmlFile, &q)

	if err != nil {
		log.Fatal("Error:", err)
	}

	// ComObjectRef -> DPT
	cache := make(ManufacturerData)

	// ComObject ID -> DPT
	comObjectDptCache := make(map[string]string)

	// Collect all ComObjects of the file, most important are DPT and ObjectSize
	for _, c := range q.ApplicationProgram.ComObjects {
		ComObjects[c.ID] = c
		if c.DatapointType == "" {
			if c.ObjectSize == "" {
				// Skip objects without type or size.
				log.Println("Skipping ComObject", c.ID)
				continue
			} else {
				c.DatapointType = dptForObjectSize(c.ObjectSize)
			}
		}
		comObjectDptCache[c.ID] = c.DatapointType
	}

	for _, c := range q.ApplicationProgram.ComObjectRefs {
		ComObjectRefs[c.ID] = c
		dpt := c.DatapointType
		if dpt == "" {
			// ask comobject cache
			dpt = comObjectDptCache[c.RefID]
		}

		if dpt == "" {
			//log.Printf("ComObjRef '%v' has no DPT??? file: %v", c.RefId, path)
			continue
		}

		cache[c.ID] = dpt
	}

	return cache
}

type ManufacturerData2 struct {
	ApplicationProgram ApplicationProgram `xml:"ManufacturerData>Manufacturer>ApplicationPrograms>ApplicationProgram"`
}

type ApplicationProgram struct {
	ComObjects    []ComObject    `xml:"Static>ComObjectTable>ComObject"`
	ComObjectRefs []ComObjectRef `xml:"Static>ComObjectRefs>ComObjectRef"`

	Name string `xml:"Name,attr"`
	ID   string `xml:"Id,attr"`
}

// Why do the ComObjectRef have the same data as the ComObject?
type ComObjectRef struct {
	Name              string `xml:"Name,attr"`
	RefID             string `xml:"RefId,attr"`
	Tag               string `xml:"Tag,attr"`
	Text              string `xml:"Text,attr"`
	TransmitFlag      string `xml:"TransmitFlag,attr"`
	ReadFlag          string `xml:"ReadFlag,attr"` // Enabled or Disabled
	WriteFlag         string `xml:"WriteFlag,attr"`
	CommunicationFlag string `xml:"CommunicationFlag,attr"`
	UpdateFlag        string `xml:"UpdateFlag,attr"`
	ReadOnInitFlag    string `xml:"ReadOnInitFlag,attr"`
	ObjectSize        string `xml:"ObjectSize,attr"`
	FunctionText      string `xml:"FunctionText,attr"`
	ID                string `xml:"Id,attr"`
	DatapointType     string `xml:"DatapointType,attr"`
}

type ComObject struct {
	ID                string `xml:"Id,attr"`
	Name              string `xml:"Name,attr"` // For what is name used? not visible in ETS
	Number            int    `xml:"Number,attr"`
	FunctionText      string `xml:"FunctionText,attr"`
	ObjectSize        string `xml:"ObjectSize,attr"`
	Text              string `xml:"Text,attr"`
	WriteFlag         string `xml:"WriteFlag,attr"`
	ReadFlag          string `xml:"ReadFlag,attr"`
	ReadOnInitFlag    string `xml:"ReadOnInitFlag,attr"`
	CommunicationFlag string `xml:"CommunicationFlag,attr"`
	TransmitFlag      string `xml:"TransmitFlag,attr"`
	UpdateFlag        string `xml:"UpdateFlag,attr"`
	DatapointType     string `xml:"DatapointType,attr"`
}

// Convert 1 Bit, 2 Bytes, 4 Bit etc to DPT
// EIS data types
// https://www.domotiga.nl/projects/selfbus-knx-eib/wiki/Datatypes
func dptForObjectSize(os string) string {
	switch os {
	case "1 Bit":
		// DPT 1001
		return "DPT-1"
	case "2 Bit":
		return "DPT-2"
	case "4 Bit":
		return "DPT-3"
	case "1 Byte":
		// It can be 4,5,6 - which is the most likely?
		return "DPT-6"
	case "2 Bytes":
		// It can be 7,8,9
		return "DPT-7"
	case "3 Bytes":
		return "DPT-10"
	case "4 Bytes":
		// It can be 12,13,14
		return "DPT-14"
	case "14 Bytes":
		return "DPT-16"
	}

	// missing: 3bit, 5bit, 6bit, 7bit, 6bytes, 8bytes, 10bytes

	//log.Println("Could not convert ObjectSize to DPT", os)
	return ""
}

func ParseManufacturers(path string) {
	Manufacturers = make(map[string]Manufacturer)

	xmlFile, err := ioutil.ReadFile(path + "knx_master.xml")
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var q XMLKNXMaster
	err = xml.Unmarshal(xmlFile, &q)

	if err != nil {
		log.Fatal("Error:", err)
	}

	for _, m := range q.Manufacturers {
		Manufacturers[m.ID] = m
	}

}

type XMLKNXMaster struct {
	Manufacturers []Manufacturer `xml:"MasterData>Manufacturers>Manufacturer"`
}

type Manufacturer struct {
	ImportGroup        string `xml:"ImportGroup,attr"`
	CompatibilityGroup string `xml:"CompatibilityGroup,attr"`
	ID                 string `xml:"Id,attr"`
	KnxManufacturerID  string `xml:"KnxManufacturerId,attr"`
	Name               string `xml:"Name,attr"`
	DefaultLanguage    string `xml:"DefaultLanguage,attr"`
	ImportRestriction  string `xml:"ImportRestriction,attr"`
}
