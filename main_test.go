package main

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"sync"
	"testing"
)

func TestLoadBalancer_AddBackend(t *testing.T) {
	lb := &LoadBalancer{}

	backendURL, _ := url.Parse("http://localhost:8081")
	backend := &Backend{
		URL:          backendURL,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL),
	}

	lb.AddBackend(backend)

	if len(lb.backends) != 1 {
		t.Errorf("Expected 1 backend, got %d", len(lb.backends))
	}

	if !reflect.DeepEqual(lb.backends[0], backend) {
		t.Errorf("Expected backend %v, got %v", backend, lb.backends[0])
	}
}

func TestLoadBalancer_NextBackend(t *testing.T) {
	lb := &LoadBalancer{}

	backendURL1, _ := url.Parse("http://localhost:8081")
	backend1 := &Backend{
		URL:          backendURL1,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL1),
	}

	backendURL2, _ := url.Parse("http://localhost:8082")
	backend2 := &Backend{
		URL:          backendURL2,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL2),
	}

	lb.AddBackend(backend1)
	lb.AddBackend(backend2)

	// Test round-robin selection
	next1 := lb.NextBackend()
	if next1 != backend1 {
		t.Errorf("Expected backend1, got %v", next1)
	}

	next2 := lb.NextBackend()
	if next2 != backend2 {
		t.Errorf("Expected backend2, got %v", next2)
	}

	next3 := lb.NextBackend()
	if next3 != backend1 {
		t.Errorf("Expected backend1, got %v", next3)
	}
}

func TestLoadBalancer_ServeHTTP_NoBackends(t *testing.T) {
	lb := &LoadBalancer{}

	req := httptest.NewRequest("GET", "http://example.com", nil)
	w := httptest.NewRecorder()

	lb.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status Service Unavailable, got %v", resp.StatusCode)
	}
}

func TestLoadBalancer_Concurrency(t *testing.T) {
	lb := &LoadBalancer{}

	backendURL1, _ := url.Parse("http://localhost:8081")
	backend1 := &Backend{
		URL:          backendURL1,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL1),
	}

	backendURL2, _ := url.Parse("http://localhost:8082")
	backend2 := &Backend{
		URL:          backendURL2,
		ReverseProxy: httputil.NewSingleHostReverseProxy(backendURL2),
	}

	lb.AddBackend(backend1)
	lb.AddBackend(backend2)

	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			lb.NextBackend()
		}()
	}

	wg.Wait()

	// Check if the load balancer still works correctly after concurrent access
	next1 := lb.NextBackend()
	if next1 != backend1 {
		t.Errorf("Expected backend1, got %v", next1)
	}

	next2 := lb.NextBackend()
	if next2 != backend2 {
		t.Errorf("Expected backend2, got %v", next2)
	}
}
