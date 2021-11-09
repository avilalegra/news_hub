package rss

import (
	"avilego.me/recent_news/news"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"time"
)

type Source struct {
	Url        string
	HttpClient HttpClient
}

func (s Source) Fetch() (*Channel, error) {
	xmlText, err := s.HttpClient.Get(s.Url)
	if err != nil {
		return nil, err
	}
	channel, err := Parse(xmlText)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"_ link"`
	Description string   `xml:"description"`
	Language    string   `xml:"language"`
	Items       []Item   `xml:"item"`
}

func (ch Channel) GetNews() []news.Preview {
	extSource := news.Source{
		Title:       ch.Title,
		Link:        ch.Link,
		Description: ch.Description,
		Language:    ch.Language,
	}

	previews := make([]news.Preview, len(ch.Items))

	for i, item := range ch.Items {
		ext := news.Preview{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Source:      &extSource,
		}
		previews[i] = ext
	}

	return previews
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

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubTime     PubTime  `xml:"pubDate"`
}

type HttpClient interface {
	Get(url string) ([]byte, error)
}

type DefaultHttpClient struct{}

func (f DefaultHttpClient) Get(url string) ([]byte, error) {
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

func NewSource(url string) Source {
	return Source{
		Url:        url,
		HttpClient: DefaultHttpClient{},
	}
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

type PubTime struct {
	UnixTime int64
}

func (pubTime *PubTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var timeExpr string
	err := d.DecodeElement(&timeExpr, &start)
	if err != nil {
		return err
	}
	unixTime, err := parsePubTime(timeExpr)
	if err != nil {
		return err
	}
	pubTime.UnixTime = unixTime
	return nil
}

func parsePubTime(timeExpr string) (int64, error) {
	timeFormats := []string{
		time.RFC1123Z,
		time.RFC1123,
	}

	for _, fmt := range timeFormats {
		if parsedTime, err := time.Parse(fmt, timeExpr); err == nil {
			return parsedTime.Unix(), nil
		}
	}

	return 0, errors.New("not supported time format")
}

func newTokenReader(xmlText []byte) xml.TokenReader {
	baseDecoder := xml.NewDecoder(bytes.NewReader(xmlText))
	baseDecoder.Strict = false
	//This allow marking tags without namespaces as ex. xml:"_ link"
	//so it doesn't collision with ex. "atom:link"
	baseDecoder.DefaultSpace = "_"
	return Trimmer{baseDecoder}
}
