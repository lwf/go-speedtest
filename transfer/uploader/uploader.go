package uploader

import (
	"crypto/rand"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/lwf/go-speedtest/transfer"
)

var payloadSizes = []int{250000, 500000}

type Uploader struct {
	url string
	buf []byte
}

func (u *Uploader) Init(url string) chan int {
	u.buf = make([]byte, 500000)
	_, err := rand.Read(u.buf)
	assert(err)
	u.url = url
	queue := make(chan int, 25*len(payloadSizes))
	for _, p := range payloadSizes {
		for i := 0; i < 25; i++ {
			queue <- p
		}
	}
	return queue
}

func (u *Uploader) Sink(size int) transfer.Sink {
	p := url.Values{"content1": {string(u.buf[0:size])}}
	r, err := http.NewRequest("POST", u.url, strings.NewReader(p.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	assert(err)
	c := &http.Client{}
	s := &sink{
		request: r,
		client:  c,
		size:    int64(size),
	}
	return s
}

func (*Uploader) String() string {
	return "Upload"
}

type sink struct {
	request *http.Request
	client  *http.Client
	size    int64
}

func (s *sink) Process() (int64, error) {
	r, err := s.client.Do(s.request)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	return s.size, nil
}

type canceler interface {
	CancelRequest(*http.Request)
}

func (s *sink) Close() error { // TODO: implement this somehow
	return nil
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
