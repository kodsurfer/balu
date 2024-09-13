package main

import (
	"fmt"
	"log"
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
	mutex    sync.Mutex
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

	backends := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, backendURL := range backends {
		parsedURL, err := url.Parse(backendURL)
		if err != nil {
			log.Fatalf("Failed to parse backend URL %s: %v", backendURL, err)
		}
		b.AddBackend(&Backend{URL: parsedURL, ReverseProxy: httputil.NewSingleHostReverseProxy(parsedURL)})
	}

	fmt.Println("Load balancer started")
	err = http.ListenAndServe(":8080", b)
	if err != nil {
		log.Fatal(err)
	}
}
