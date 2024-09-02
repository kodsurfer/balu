package main

import (
	"fmt"
	"net/http"
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

func (b *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := b.NextBackend()
	if backend == nil {
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}
	backend.ReverseProxy.ServeHTTP(w, r)
}

func main() {
	b := &LoadBalancer{}

	back1, err := url.Parse("http://localhost:8081")
	if err != nil {
		fmt.Printf("url parse error: %v", err)
	}
	back2, err := url.Parse("http://localhost:8082")
	if err != nil {
		fmt.Printf("url parse error: %v", err)
	}
	back3, err := url.Parse("http://localhost:8083")
	if err != nil {
		fmt.Printf("url parse error: %v", err)
	}

	b.AddBackend(&Backend{URL: back1, ReverseProxy: httputil.NewSingleHostReverseProxy(back1)})
	b.AddBackend(&Backend{URL: back2, ReverseProxy: httputil.NewSingleHostReverseProxy(back2)})
	b.AddBackend(&Backend{URL: back3, ReverseProxy: httputil.NewSingleHostReverseProxy(back3)})

	fmt.Println("Load balancer started")
	http.ListenAndServe(":8080", b)
}
