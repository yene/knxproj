package knxproj

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

type Device struct {
	ID               string  `json:"-"`
	ProductRefID     string  `json:"-"`
	Address          string  `json:"address"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Productname      string  `json:"productname"`
	Manufacturer     string  `json:"manufacturer"`
	Application      string  `json:"-"`
	Type             string  `json:"-"`
	Fingerprint      string  `json:"-"`
	SerialNumber     string  `json:"serialnumber"`
	Objects          Objects `json:"objects"`
	DownloadRequired bool    `json:"downloadrequired"`
}

type Object struct {
	Number       int      `json:"number"`
	Name         string   `json:"name"`
	Function     string   `json:"function"`
	Description  string   `json:"description"`
	Groupaddress []string `json:"groupaddress"`
	Bit          int      `json:"bit"`
	Flags        string   `json:"flags"`
	DPT          string   `json:"dpt"`
	RefID        string   `json:"-"`
	Channel      string   `json:"channel"`
	ReadOnInit   bool     `json:"-"`
}

type Objects []Object

func (slice Objects) Len() int {
	return len(slice)
}

func (slice Objects) Less(i, j int) bool {
	return slice[i].Number < slice[j].Number
}

func (slice Objects) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func (d Device) String() string {
	return fmt.Sprintf("%s - %s %s - %s", d.Address, d.Description, d.Productname, d.Manufacturer)
}

type ParameterInstanceRef struct {
	RefID string `xml:"RefId,attr"`
	Value string `xml:"Value,attr"`
}

type ComObjectInstanceRef struct {
	RefID             string     `xml:"RefId,attr"`
	IsActive          bool       `xml:"IsActive,attr"`
	DatapointType     string     `xml:"DatapointType,attr"`
	Text              string     `xml:"Text,attr"`
	ChannelID         string     `xml:"ChannelId,attr"`
	Connectors        Connectors `xml:"Connectors"`
	Description       string     `xml:"Description,attr"`
	TransmitFlag      string     `xml:"TransmitFlag,attr"`
	UpdateFlag        string     `xml:"UpdateFlag,attr"`
	WriteFlag         string     `xml:"WriteFlag,attr"`
	ReadFlag          string     `xml:"ReadFlag,attr"`
	CommunicationFlag string     `xml:"CommunicationFlag,attr"`
}

type Connectors struct {
	Addresses []GroupAddressRef `xml:",any"`
}

type GroupAddressRef struct {
	GroupAddressRefID string `xml:"GroupAddressRefId,attr"`
	Type              string // "send" or "receive"
}

func NewDevice(area string, line string, xmlDevice DeviceInstance) *Device {
	d := new(Device)
	d.ID = xmlDevice.ID
	d.ProductRefID = xmlDevice.ProductRefID
	d.Name = xmlDevice.Name

	h := Hardware[d.ProductRefID]
	d.Productname = h.Text
	d.SerialNumber = h.SerialNumber
	d.Description = xmlDevice.Description

	if xmlDevice.LastModified.After(xmlDevice.LastDownload) {
		d.DownloadRequired = true
	}

	// Parse the manufacturer from ProductRefID
	mID := strings.Split(d.ProductRefID, "_")[0]
	d.Manufacturer = Manufacturers[mID].Name

	// Expand the physical address
	d.Address = area + "." + line + "." + xmlDevice.Address

	// iterating over the device objects
	d.Objects = make(Objects, 0)
	for _, c := range xmlDevice.ComObjectInstanceRefs {
		o := Object{}
		o.RefID = c.RefID
		o.Description = c.Description
		o.Channel = findChannelName(c.ChannelID, xmlDevice.ChannelInstances)
		ref := ComObjectRefs[c.RefID]
		comObject := ComObjects[ref.RefID]
		if comObject.ObjectSize == "LegacyVarData" {
			continue
		}

		// Use the values from comobject first
		o.Name = comObject.Text
		o.Function = comObject.FunctionText
		o.Number = comObject.Number
		o.Bit = sizeStringToInt(comObject.ObjectSize)
		o.DPT = comObject.DatapointType
		o.ReadOnInit = t2b(comObject.ReadOnInitFlag)
		o.DPT = stringToDPT(comObject.DatapointType)

		o.Flags = flags(o.Flags, comObject.CommunicationFlag, "C")
		o.Flags = flags(o.Flags, comObject.ReadFlag, "R")
		o.Flags = flags(o.Flags, comObject.WriteFlag, "W")
		o.Flags = flags(o.Flags, comObject.TransmitFlag, "T")
		o.Flags = flags(o.Flags, comObject.UpdateFlag, "U")

		// Overwrite values with the ones from ComObjectRef
		if ref.Text != "" {
			o.Name = ref.Text
		}
		if ref.FunctionText != "" {
			o.Function = ref.FunctionText
		}

		if ref.ObjectSize != "" {
			o.Bit = sizeStringToInt(ref.ObjectSize)
		}

		if ref.DatapointType != "" {
			o.DPT = stringToDPT(ref.DatapointType)
		}

		o.Flags = flags(o.Flags, ref.CommunicationFlag, "C")
		o.Flags = flags(o.Flags, ref.ReadFlag, "R")
		o.Flags = flags(o.Flags, ref.WriteFlag, "W")
		o.Flags = flags(o.Flags, ref.TransmitFlag, "T")
		o.Flags = flags(o.Flags, ref.UpdateFlag, "U")

		// Overwrite values with the ones from the ComObjectInstanceRef inside 0.xml
		if c.DatapointType != "" {
			o.DPT = stringToDPT(c.DatapointType)
		}

		if c.Text != "" {
			o.Name = c.Text
		}

		o.Flags = flags(o.Flags, c.CommunicationFlag, "C")
		o.Flags = flags(o.Flags, c.ReadFlag, "R")
		o.Flags = flags(o.Flags, c.WriteFlag, "W")
		o.Flags = flags(o.Flags, c.TransmitFlag, "T")
		o.Flags = flags(o.Flags, c.UpdateFlag, "U")

		o.Groupaddress = make([]string, 0)
		for _, a := range c.Connectors.Addresses {
			addr := GAddresses[a.GroupAddressRefID]
			// If no description is set use the first address as fallback.
			if o.Description == "" {
				o.Description = addr.Name
			}
			o.Groupaddress = append(o.Groupaddress, addr.Address)
		}

		d.Objects = append(d.Objects, o)
	}

	sort.Sort(d.Objects)

	return d
}

func findChannelName(channelID string, channels []ChannelInstance) string {
	for _, c := range channels {
		if c.RefID == channelID {
			return c.Name
		}
	}
	return ""
}

func sizeStringToInt(size string) int {
	switch size {
	case "1 Bit":
		return 1
	case "2 Bit":
		return 2
	case "3 Bit":
		return 3
	case "4 Bit":
		return 4
	case "5 Bit":
		return 5
	case "6 Bit":
		return 6
	case "7 Bit":
		return 7
	case "1 Byte":
		return 8
	case "2 Bytes":
		return 16
	case "3 Bytes":
		return 24
	case "4 Bytes":
		return 32
	case "6 Bytes":
		return 48
	case "8 Bytes":
		return 64
	case "10 Bytes":
		return 80
	case "14 Bytes":
		return 112
	}

	log.Fatal("Could not convert Byte text to int", size)

	return 0
}

func t2b(t string) bool {
	return t == "Enabled"
}

func stringToDPT(dpt string) string {
	if len(dpt) == 0 {
		return ""
	}
	split := strings.Split(dpt, "-")
	mainType, _ := strconv.Atoi(split[1])
	subType := 0
	if split[0] == "DPST" {
		subType, _ = strconv.Atoi(split[2])
	}
	return fmt.Sprintf("%d.%03d", mainType, subType)
}

func flags(flags, value, flag string) string {
	if value == "" {
		return flags
	}

	contains := strings.Contains(flags, flag)

	if value == "Enabled" && !contains {
		return flags + flag
	}

	if value == "Disabled" && contains {
		return strings.Replace(flags, flag, "", 1)
	}

	return flags
}
