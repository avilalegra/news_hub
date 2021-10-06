package channel

import (
	"encoding/xml"
	"time"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName       xml.Name    `xml:"channel"`
	Title         string      `xml:"title"`
	Link          string      `xml:"link"`
	Description   string      `xml:"description"`
	Language      string      `xml:"language"`
	LastBuildDate ChannelTime `xml:"lastBuildDate"`
}

type ChannelTime struct {
	*time.Time
}

func (cht *ChannelTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var timexpr string

	d.DecodeElement(&timexpr, &start)
	parse, err := time.Parse(time.RFC1123Z, timexpr)
	if err != nil {
		return err
	}

	cht.Time = &parse

	return nil
}

func Parse(xmlText []byte) (*Channel, error) {

	var rss rss
	err := xml.Unmarshal(xmlText, &rss)

	if err != nil {
		return nil, err
	}

	return &rss.Channel, nil
}
