package channel

import (
	"encoding/xml"
)

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

func Parse(xmlText []byte) (*Channel, error) {

	var rss rss
	err := xml.Unmarshal(xmlText, &rss)

	if err != nil {
		return nil, err
	}

	return &rss.Channel, nil
}
