package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"load_balancer/balancer"
	"load_balancer/config"
	"load_balancer/ratelimiter"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	bal := balancer.NewBalancer(cfg.Backends)

	limiter := ratelimiter.NewRateLimiter(cfg.DefaultRateLimit.Capacity, cfg.DefaultRateLimit.RefillRate)
	for ip, limit := range cfg.RateLimits {
		limiter.SetClientLimit(ip, limit.Capacity, limit.RefillRate)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		log.Printf("Incoming request from %s: %s %s", host, r.Method, r.URL.Path)

		if !limiter.Allow(host) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for %s", host)
			return
		}

		bal.ServeHTTP(w, r)
	})

	log.Printf("Starting server on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
