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
	mu       sync.Mutex
	current  int
}

func NewBalancer(backendURLs []string) *Balancer {
	backends := make([]*url.URL, 0, len(backendURLs))
	for _, addr := range backendURLs {
		parsed, err := url.Parse(addr)
		if err != nil {
			log.Fatalf("Invalid backend URL %s: %v", addr, err)
		}
		backends = append(backends, parsed)
	}
	return &Balancer{
		backends: backends,
	}
}

func (b *Balancer) getNextBackend() *url.URL {
	b.mu.Lock()
	defer b.mu.Unlock()
	backend := b.backends[b.current]
	b.current = (b.current + 1) % len(b.backends)
	return backend
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := b.getNextBackend()

	proxy := httputil.NewSingleHostReverseProxy(backend)
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Backend %s unavailable: %v", backend, err)
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
	}

	log.Printf("Forwarding request %s to %s", r.URL.Path, backend)
	proxy.ServeHTTP(w, r)
}
