package servers

import (
	"math"
	"sort"

	"github.com/lwf/go-speedtest/config"
)

type distanceSorter struct {
	config  config.Config
	servers []Server
}

func (d *distanceSorter) Less(i, j int) bool {
	return distance(d.servers[i], d.config) < distance(d.servers[j], d.config)
}

func (d *distanceSorter) Len() int {
	return len(d.servers)
}

func (d *distanceSorter) Swap(i, j int) {
	d.servers[i], d.servers[j] = d.servers[j], d.servers[i]
}

func distance(server Server, config config.Config) float64 {
	lat1 := degToRad(server.Latitude)
	lon1 := degToRad(server.Longitude)
	lat2 := degToRad(config.Client.Latitude)
	lon2 := degToRad(config.Client.Longitude)
	dlon := lon2 - lon1
	dlat := lat2 - lat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c * 6373
}

func degToRad(d float64) float64 {
	return d * math.Pi / 180
}

func SortByDistance(config config.Config, servers []Server) {
	sorter := &distanceSorter{
		config:  config,
		servers: servers,
	}
	sort.Sort(sorter)
}
