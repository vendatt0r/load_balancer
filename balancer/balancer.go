package balancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Balancer struct {
	backends []*url.URL
	current  uint32
}

func NewBalancer(backendStrings []string) *Balancer {
	var backends []*url.URL
	for _, b := range backendStrings {
		u, err := url.Parse(b)
		if err != nil {
			log.Fatalf("Invalid backend URL: %v", err)
		}
		backends = append(backends, u)
	}
	return &Balancer{backends: backends}
}

func (b *Balancer) nextBackend() *url.URL {
	index := atomic.AddUint32(&b.current, 1)
	return b.backends[int(index)%len(b.backends)]
}

func (b *Balancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tries := 0
	maxTries := len(b.backends)

	for tries < maxTries {
		backend := b.nextBackend()
		proxy := httputil.NewSingleHostReverseProxy(backend)

		proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
			log.Printf("Backend %s failed: %v. Retrying...", backend.String(), err)
			tries++
			if tries >= maxTries {
				http.Error(rw, "All backends are down", http.StatusServiceUnavailable)
				return
			}
			b.ServeHTTP(rw, req) // попробовать другой
		}

		log.Printf("Forwarding request %s %s to backend %s", r.Method, r.URL.Path, backend.String())
		proxy.ServeHTTP(w, r)
		return
	}
}
