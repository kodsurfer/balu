package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
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

	server := &http.Server{
		Addr:    ":8080",
		Handler: b,
	}

	go func() {
		log.Println("Load balancer started at :8080")
		if err := http.ListenAndServe(":8080", b); err != nil {
			log.Fatal(err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
