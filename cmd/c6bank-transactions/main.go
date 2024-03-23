package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const MAX_UPLOAD_SIZE = 10 * 1024 * 1024 // 10MB

func main() {
	log.SetFlags(0)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/upload", uploadHandler)

	port := getenv("PORT", "4500")

	httpServer := http.Server{
		Addr:              "127.0.0.1:" + port,
		Handler:           mux,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 1,
		WriteTimeout:      time.Second * 10,
		IdleTimeout:       time.Minute,
		MaxHeaderBytes:    1 << 22, // 4MB
	}

	go func() {
		log.Printf("server listening at :%s", port)

		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalln("ERROR", err)
		}
	}()

	<-ctx.Done()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("ERROR http server shutdown: %s", err)
	}
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
