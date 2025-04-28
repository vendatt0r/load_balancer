package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Balancer struct {
	backends []*url.URL
	index    int
	mu       sync.Mutex
}

func NewBalancer(backendUrls []string) *Balancer {
	var urls []*url.URL
	for _, addr := range backendUrls {
		parsed, err := url.Parse(addr)
		if err != nil {
			log.Fatalf("Invalid backend URL: %s", addr)
		}
		urls = append(urls, parsed)
	}
	return &Balancer{backends: urls}
}

func (b *Balancer) NextBackend() *url.URL {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.backends) == 0 {
		return nil
	}

	backend := b.backends[b.index]
	b.index = (b.index + 1) % len(b.backends)
	return backend
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := b.NextBackend()
	if backend == nil {
		http.Error(w, "No available backend", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backend)
	log.Printf("Forwarding request %s %s to backend %s", r.Method, r.URL.Path, backend.Host)
	proxy.ServeHTTP(w, r)
}
