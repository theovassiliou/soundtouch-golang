package soundtouch

import (
	"encoding/xml"
)

// ContentItem describe a playable item
type ContentItem struct {
	Type         string `xml:"type,attr" json:",omitempty"`
	Source       Source `xml:"source,attr" json:",omitempty"`
	Location     string `xml:"location,attr" json:",omitempty"`
	Name         string `xml:"itemName" json:",omitempty"`
	IsPresetable bool   `xml:"isPresetable,attr"`
}

// Select sends the select command to the soundtouch system
func (s *Speaker) Select(item ContentItem) error {
	data, err := xml.Marshal(item)
	if err != nil {
		return err
	}

	_, err = s.SetData("select", data)
	if err != nil {
		return err
	}
	return nil
}
