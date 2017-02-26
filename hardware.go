package knxproj

import (
	"encoding/xml"
	"io/ioutil"
	"log"
)

func parseHardwareData(path string) {
	xmlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var q XMLHardware
	err = xml.Unmarshal(xmlFile, &q)

	if err != nil {
		log.Fatal("Error:", err)
	}

	for _, h := range q.Hardware {
		for _, p := range h.Products {
			p.SerialNumber = h.SerialNumber
			Hardware[p.ID] = p
		}
	}
}

type XMLHardware struct {
	//Manufacturer       Manufacturer         `xml:"ManufacturerData>Manufacturer"`
	Hardware []Hardware2 `xml:"ManufacturerData>Manufacturer>Hardware>Hardware"`
	//TranslationElements []TranslationElement `xml:"ManufacturerData>Manufacturer>Languages>Language>TranslationUnit>TranslationElement"`
	//Translation        []Translation        `xml:"ManufacturerData>Manufacturer>Languages>Language>TranslationUnit>TranslationElement>Translation"`
	//Language           Language             `xml:"ManufacturerData>Manufacturer>Languages>Language"`
	//Products []Product `xml:"ManufacturerData>Manufacturer>Hardware>Hardware>Products>Product"`
	//TranslationUnit    []TranslationUnit    `xml:"ManufacturerData>Manufacturer>Languages>Language>TranslationUnit"`
}

/*
type ApplicationProgramRef struct {
	RefID string `xml:"RefId,attr"`
}
type Manufacturer struct {
	RefID string `xml:"RefId,attr"`
}
type Language struct {
	Identifier string `xml:"Identifier,attr"`
}
type TranslationUnit struct {
	RefID string `xml:"RefId,attr"`
}
type TranslationElement struct {
	RefID string `xml:"RefId,attr"`
}
type Translation struct {
	AttributeName string `xml:"AttributeName,attr"`
	Text          string `xml:"Text,attr"`
}*/
type Hardware2 struct {
	//IsCoupler             string `xml:"IsCoupler,attr"`
	SerialNumber string `xml:"SerialNumber,attr"`
	//VersionNumber         string `xml:"VersionNumber,attr"`
	//HasApplicationProgram string `xml:"HasApplicationProgram,attr"`
	//BusCurrent            string `xml:"BusCurrent,attr"`
	//ID                    string `xml:"Id,attr"`
	//HasIndividualAddress  string `xml:"HasIndividualAddress,attr"`
	//OriginalManufacturer  string `xml:"OriginalManufacturer,attr"`
	//Name                  string `xml:"Name,attr"`
	Products []Product `xml:"Products>Product"`
}
type Product struct {
	DefaultLanguage    string `xml:"DefaultLanguage,attr"`
	Hash               string `xml:"Hash,attr"`
	ID                 string `xml:"Id,attr"`
	Text               string `xml:"Text,attr"`
	OrderNumber        string `xml:"OrderNumber,attr"`
	IsRailMounted      string `xml:"IsRailMounted,attr"`
	VisibleDescription string `xml:"VisibleDescription,attr"`
	SerialNumber       string
}
