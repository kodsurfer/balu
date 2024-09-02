package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
}

type LoadBalancer struct {
	backends []*Backend
	mutex    *sync.Mutex
	next     int
}

func (b *LoadBalancer) AddBackend(backend *Backend) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.backends = append(b.backends, backend)
}

func (b *LoadBalancer) NextBackend() *Backend {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if len(b.backends) == 0 {
		return nil
	}

	res := b.backends[b.next]
	b.next = (b.next + 1) % len(b.backends)

	return res
}

func main() {
	fmt.Println("Let's create load balancer")
}
