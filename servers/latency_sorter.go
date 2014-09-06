package servers

import (
	"sort"
	"time"

	"github.com/lwf/go-speedtest/config"
	"github.com/lwf/go-speedtest/util"
)

type latencySorter struct {
	config  config.Config
	servers []Server
	latency map[int]time.Duration
}

func (l *latencySorter) Swap(i, j int) {
	l.servers[i], l.servers[j] = l.servers[j], l.servers[i]
}

func (l *latencySorter) Len() int {
	return len(l.servers)
}

func (l *latencySorter) Less(i, j int) bool {
	return l.cacheLatency(l.servers[i]) < l.cacheLatency(l.servers[j])
}

func (l *latencySorter) cacheLatency(server Server) time.Duration {
	if l.latency[server.Id] == 0 {
		l.latency[server.Id] = latency(server)
	}
	return l.latency[server.Id]
}

func latency(server Server) time.Duration {
	t1 := time.Now()
	_, err := util.HttpGet(server.Url)
	assert(err)
	return time.Since(t1)
}

func SortByLatency(config config.Config, servers []Server) {
	latency := make(map[int]time.Duration)
	sorter := &latencySorter{
		config:  config,
		servers: servers,
		latency: latency,
	}
	sort.Sort(sorter)
}
