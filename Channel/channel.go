package channel

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName       xml.Name      `xml:"channel"`
	Title         string        `xml:"title"`
	Link          string        `xml:"_ link"`
	Description   string        `xml:"description"`
	Language      string        `xml:"language"`
	LastBuildDate ChannelTime   `xml:"lastBuildDate"`
	Items         []ChannelItem `xml:"item"`
}

type ChannelTime struct {
	*time.Time
}

type ChannelItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

func (cht *ChannelTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var timexpr string

	d.DecodeElement(&timexpr, &start)
	parse, err := time.Parse(time.RFC1123, timexpr)
	if err != nil {
		return err
	}

	cht.Time = &parse

	return nil
}

type Trimmer struct {
	dec *xml.Decoder
}

func (tr Trimmer) Token() (xml.Token, error) {
	t, err := tr.dec.Token()
	if cd, ok := t.(xml.CharData); ok {

		t = xml.CharData(bytes.TrimSpace(cd))
	}
	return t, err
}

func newTokenReader(xmlText []byte) xml.TokenReader {
	baseDecoder := xml.NewDecoder(bytes.NewReader(xmlText))

	baseDecoder.Strict = false
	//This allow marking tags without namespaces as ex. xml:"_ link"
	//so it doesn't collision with ex. "atom:link"
	baseDecoder.DefaultSpace = "_"

	return Trimmer{baseDecoder}
}

func Parse(xmlText []byte) (*Channel, error) {
	var rss rss
	dec := xml.NewTokenDecoder(newTokenReader(xmlText))
	err := dec.Decode(&rss)

	if err != nil {
		return nil, err
	}

	return &rss.Channel, err
}

func NewSource(url string) *Source {
	return &Source{
		Url:            url,
		ContentFetcher: DefaultContentFetcher{},
	}
}

type Source struct {
	Url            string
	ContentFetcher ContentFetcher
}

func (src Source) Fetch() (*Channel, error) {
	xmlText, err := src.ContentFetcher.Get(src.Url)

	if err != nil {
		return nil, err
	}

	channel, err := Parse(xmlText)

	if err != nil {
		return nil, err
	}

	return channel, nil
}

type ContentFetcher interface {
	Get(url string) ([]byte, error)
}

type DefaultContentFetcher struct{}

func (f DefaultContentFetcher) Get(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
