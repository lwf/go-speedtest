package transfer

import (
	"log"
	"sync"
	"sync/atomic"
)

type Result struct {
	n   int64
	err error
}

type Manager struct {
	url  string
	done chan struct{}
}

func (d *Manager) Run(worker Worker, n int, progress func()) int64 {
	queue := worker.Init(d.url)
	var count int64
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			d.workerLoop(worker, queue, progress, &count)
		}()
	}
	wg.Wait()
	return count
}

func (d *Manager) workerLoop(worker Worker, queue <-chan int, progress func(), c *int64) {
	results := make(chan Result)
	for {
		select {
		case <-d.done:
			return
		default:
		}

		select {
		case item := <-queue:
			sink := worker.Sink(item)
			go func() {
				defer sink.Close()
				n, err := sink.Process()
				progress()
				results <- Result{
					n:   n,
					err: err,
				}
			}()
			select {
			case <-d.done:
				sink.Close()
				select {
				case result := <-results:
					atomic.AddInt64(c, result.n)
				}
				return
			case result := <-results:
				assert(result.err)
				atomic.AddInt64(c, result.n)
			}
		default:
			return
		}
	}
}

func (d *Manager) Close() {
	close(d.done)
}

func NewManager(url string) *Manager {
	done := make(chan struct{})
	return &Manager{
		url:  url,
		done: done,
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
