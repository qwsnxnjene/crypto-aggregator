package main

import (
	"context"
	"crypto-aggregator/internal/aggregator"
	"crypto-aggregator/internal/handler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	agg := aggregator.NewAggregator()
	handle := handler.NewHandler(&agg)
	srv := &http.Server{Addr: ":8080"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() { agg.Run(ctx) }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		cancel()
		srv.Shutdown(context.Background())
	}()

	http.Handle("/prices", corsMiddleware(http.HandlerFunc(handle.Prices)))
	srv.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
