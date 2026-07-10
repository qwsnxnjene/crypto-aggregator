package main

import (
	"context"
	"crypto-aggregator/internal/aggregator"
	"crypto-aggregator/internal/handler"
	"crypto-aggregator/internal/storage"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	dbConn := dsn
	conn, err := storage.NewStorage(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	agg := aggregator.NewAggregator(conn)
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
	http.Handle("/ws", corsMiddleware(http.HandlerFunc(handle.PricesWS)))
	srv.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
