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

func main() {
	fmt.Println("Let's create load balancer")
}
