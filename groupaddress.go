package knxproj

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const unspecified = -1

type GroupAddress struct {
	Address       string
	Name          string
	ID            string
	MainType      int
	SubType       int
	LinkedDevices int
}

func (g GroupAddress) String() string {
	return fmt.Sprintf("%s %s %s", g.Address, g.Name, g.DataPointType())
}
func (g GroupAddress) DataPointType() string {
	if g.MainType == unspecified {
		return ""
	}

	return fmt.Sprintf("%d.%03d", g.MainType, g.SubType)
}

func NewGroupAddress(ga GroupAddress2) *GroupAddress {
	g := new(GroupAddress)
	g.MainType = unspecified
	g.SubType = 0

	g.Name = ga.Name
	if g.Name == "" {
		log.Println("no name found, using address", g.Address)
		g.Name = g.Address
	}

	// Parse the DPT if it is mentioned on the address.
	if ga.DatapointType != "" {
		split := strings.Split(ga.DatapointType, "-")
		g.MainType, _ = strconv.Atoi(split[1])
		if split[0] == "DPST" {
			g.SubType, _ = strconv.Atoi(split[2])
		}
	}

	g.ID = ga.ID
	g.Address = addressToString(ga.Address)
	return g
}

func addressToString(addr int) string {
	if AddressStyle == 3 {
		// Three Style Address 5/3/8
		// 00000000 00011111 = 0x1F 5 bits
		// 00000000 00000111 = 0x07 3 bits
		// 00000000 11111111 = 0xFF 8 bits
		main := addr >> 11 & 0x1F
		middle := addr >> 8 & 0x07
		sub := addr & 0xFF
		return fmt.Sprintf("%d/%d/%d", main, middle, sub)
	} else {
		// Two Style Address 5/11 bits
		// 00000000 00011111 = 0x1F 5 bits
		// 00000111 11111111 = 0x7FF 11 bits
		main := addr >> 11 & 0x1F
		sub := addr & 0x7FF
		return fmt.Sprintf("%d/%d", main, sub)
	}
}
