package channel

import (
	"avilego.me/news_hub/news"
	"bytes"
	"encoding/xml"
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

type NewsFetchResult struct {
	previews []news.Preview
	error    error
}

func NewSource(url string) *Source {
	return &Source{
		Url:        url,
		HttpClient: DefaultHttpClient{},
	}
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

type rss struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName       xml.Name `xml:"channel"`
	Title         string   `xml:"title"`
	Link          string   `xml:"_ link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	LastBuildDate Time     `xml:"lastBuildDate"`
	Items         []Item   `xml:"item"`
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

type Time struct {
	*time.Time
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

func (cht *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var timeExpr string
	err := d.DecodeElement(&timeExpr, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(time.RFC1123, timeExpr)
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
