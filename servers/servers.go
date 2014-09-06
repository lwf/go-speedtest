package servers

import (
	"encoding/xml"
	"log"

	"github.com/lwf/go-speedtest/util"
)

type Servers struct {
	Servers []Server `xml:"servers>server"`
}

type Server struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
	Url       string  `xml:"url,attr"`
	Country   string  `xml:"country,attr"`
	Id        int     `xml:"id,attr"`
}

func GetServers() Servers {
	buf, err := util.HttpGet("http://c.speedtest.net/speedtest-servers-static.php")
	assert(err)
	var s Servers
	err = xml.Unmarshal(buf, &s)
	assert(err)
	return s
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
