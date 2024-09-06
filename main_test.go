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
