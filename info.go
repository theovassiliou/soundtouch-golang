package soundtouch

import (
	"encoding/xml"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Info defines the Info command for the soundtouch system
type Info struct {
	DeviceID    string      `xml:"deviceID,attr" json:",omitempty"`
	Name        string      `xml:"name" json:",omitempty"`
	Type        string      `xml:"type" json:",omitempty"`
	IPAddress   []string    `xml:"networkInfo>ipAddress"`
	Component   []Component `xml:"components>component" json:",omitempty"`
	MargeUrl    string      `xml:"margeURL" json:",omitempty"`
	ModuleType  string      `xml:"moduleType" json:",omitempty"`
	Variant     string      `xml:"variant" json:",omitempty"`
	VariantMode string      `xml:"variantMode" json:",omitempty"`
	CountryCode string      `xml:"countryCode" json:",omitempty"`
	RegionCode  string      `xml:"regionCode" json:",omitempty"`
	Raw         []byte      `json:"-"`
}

// Component contains some component information
type Component struct {
	ComponentCategory string `xml:"componentCategory" json:",omitempty"`
	SoftwareVersion   string `xml:"softwareVersion" json:",omitempty"`
	SerialNumber      string `xml:"serialNumber" json:",omitempty"`
}

// Info retrieves speaker information and updates the speakers info field
func (s *Speaker) Info() (Info, error) {
	body, err := s.GetData("info")
	if err != nil {
		return Info{}, err
	}

	info := Info{
		Raw: body,
	}
	err = xml.Unmarshal(body, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

// String creates, depending of the loglevel different string representations on info message
func (s Info) String() string {
	if log.GetLevel() >= log.TraceLevel {
		return fmt.Sprintf("%v (%v): %v\n%v", s.Name, s.DeviceID, s.Type, string(s.Raw))
	}
	return fmt.Sprintf("%v (%v): %v", s.Name, s.DeviceID, s.Type)
}
