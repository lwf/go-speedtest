package config

import (
	"encoding/xml"
	"log"

	"github.com/lwf/go-speedtest/util"
)

type Config struct {
	Client Client `xml:"client"`
}

type Client struct {
	Latitude  float64 `xml:"lat,attr"`
	Longitude float64 `xml:"lon,attr"`
	Isp       string  `xml:"isp,attr"`
}

func GetConfig() Config {
	r, err := util.HttpGet("http://www.speedtest.net/speedtest-config.php")
	assert(err)
	var c Config
	err = xml.Unmarshal(r, &c)
	assert(err)
	return c
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
