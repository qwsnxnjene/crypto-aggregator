package main

import (
	"context"
	"crypto-aggregator/internal/aggregator"
	"crypto-aggregator/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	agg := aggregator.NewAggregator()
	handle := handler.NewHandler(&agg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { agg.Run(ctx) }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		cancel()
	}()

	http.Handle("/prices", corsMiddleware(http.HandlerFunc(handle.Prices)))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("[main]: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
