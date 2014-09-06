package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/lwf/go-speedtest/transfer"
)

var sizes = []int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}

type Downloader struct {
	baseUrl string
}

func (d *Downloader) Init(u string) chan int {
	queue := make(chan int, len(sizes)*4)
	parsedUrl, err := url.Parse(u)
	assert(err)
	d.baseUrl = fmt.Sprintf("%s://%s%s", parsedUrl.Scheme, parsedUrl.Host, filepath.Dir(parsedUrl.Path))
	for _, s := range sizes {
		for i := 0; i < 4; i++ {
			queue <- s
		}
	}
	return queue
}

func (d *Downloader) Sink(size int) transfer.Sink {
	url := fmt.Sprintf("%s/random%dx%d.jpg", d.baseUrl, size, size)
	r, err := http.Get(url)
	assert(err)
	return &sink{
		r: r,
	}
}

func (*Downloader) String() string {
	return "Download"
}

type sink struct {
	r *http.Response
}

func (s *sink) Process() (int64, error) {
	return io.Copy(ioutil.Discard, s.r.Body)
}

func (s *sink) Close() error {
	return s.r.Body.Close()
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
