package main

import (
	"log"
	"net/http"

	"load_balancer/balancer"
	"load_balancer/config"
	"load_balancer/limiter"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	b := balancer.NewBalancer(cfg.Backends)
	l := limiter.NewLimiter(10, 5) // 10 токенов, пополнение 5 токенов в секунду

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientID := limiter.GetClientIP(r)
		if !l.Allow(clientID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			log.Printf("Rate limit exceeded for %s", clientID)
			return
		}
		b.ServeHTTP(w, r)
	})

	log.Printf("Starting load balancer on port %s", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
