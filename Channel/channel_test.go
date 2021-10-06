package channel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func parseTime(timeExpr string) *time.Time {
	ptime, _ := time.Parse(time.RFC1123Z, timeExpr)
	return &ptime
}

var rssParsingTests = []struct {
	xml     string
	channel Channel
}{

	{
		`<?xml version="1.0" encoding="UTF-8"?><rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:media="http://search.yahoo.com/mrss/" xmlns:nyt="http://www.nytimes.com/namespaces/rss/2.0" version="2.0"><channel><title>NYT &gt; Top Stories</title><link>https://www.nytimes.com</link><description>NYT channel description</description><language>en-us</language><lastBuildDate>Sun, 19 Sep 2021 06:27:36 +0000</lastBuildDate><image><title>NYT > Top Stories</title><url>https://static01.nyt.com/images/misc/NYT_logo_rss_250x40.png</url><link>https://www.nytimes.com</link></image></channel></rss>`,
		Channel{
			Title:       `NYT > Top Stories`,
			Link:        `https://www.nytimes.com`,
			Description: `NYT channel description`,
		},
	},
}

func TestRssParsing(t *testing.T) {
	for _, tt := range rssParsingTests {
		channel, _ := Parse([]byte(tt.xml))

		assert.Equal(t, tt.channel.Title, channel.Title)
		assert.Equal(t, tt.channel.Link, channel.Link)
		assert.Equal(t, tt.channel.Description, channel.Description)
	}
}

var rssParsingErrorTests = []string{
	`<?xml version="1.0"?><rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#" xmlns="http://purl.org/rss/1.0/"><channel rdf:about="http://www.xml.com/xml/news.rss"><title>XML.com</title><link>http://xml.com/pub</link><description>rss v1 description
	</description><image rdf:resource="http://xml.com/universal/images/xml_tiny.gif" /><items><rdf:Seq><rdf:li resource="http://xml.com/pub/2000/08/09/xslt/xslt.html" /><rdf:li resource="http://xml.com/pub/2000/08/09/rdfdb/index.html" /></rdf:Seq></items><textinput rdf:resource="http://search.xml.com" /></channel><image rdf:about="http://xml.com/universal/images/xml_tiny.gif"><title>XML.com</title><link>http://www.xml.com</link><url>http://xml.com/universal/images/xml_tiny.gif</url></image><item rdf:about="http://xml.com/pub/2000/08/09/xslt/xslt.html"><title>Processing Inclusions with XSLT</title><link>http://xml.com/pub/2000/08/09/xslt/xslt.html</link><description>rss description
	</description></item><item rdf:about="http://xml.com/pub/2000/08/09/rdfdb/index.html"><title>Putting RDF to Work</title><link>http://xml.com/pub/2000/08/09/rdfdb/index.html</link><description>item description</description></item><textinput rdf:about="http://search.xml.com"><title>Search XML.com</title><description>Search XML.com's XML collection</description><name>s</name><link>http://search.xml.com</link></textinput></rdf:RDF>`,
	`<?xml version="1.0"?><xml></xl>`,
	``,
}

func TestParsingRssThrowsErrorOnInvalidXml(t *testing.T) {
	for _, xmlText := range rssParsingErrorTests {
		_, err := Parse([]byte(xmlText))
		assert.Error(t, err)
	}
}
