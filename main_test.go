package main

import (
	"net/http/httputil"
	"net/url"
	"reflect"
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
