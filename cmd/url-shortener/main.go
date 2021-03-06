package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emar-kar/urlshortener/pkg/database"
	"github.com/emar-kar/urlshortener/pkg/handler"
	"github.com/emar-kar/urlshortener/pkg/service"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("$REDIS_URL must be set")
	}

	rdb, err := database.NewDB(redisURL)
	if err != nil {
		log.Fatalf("cannot create redis client: %s", err)
	}
	if _, err := rdb.Client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("cannot logging to redis: %s", err)
	}
	defer rdb.Client.Close()

	services := service.NewService(rdb)
	handlers := handler.NewHandler(services)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handlers.InitRoutes("release"),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	go func() {
		log.Println("server starting")

		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Println("server closed")
				return
			}
			log.Printf("server start failure: %s", err)
			done <- os.Interrupt
		}
	}()

	<-done
	fmt.Println()

	maxGracefulShutdownTime := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), maxGracefulShutdownTime)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}

	log.Println("server exited properly")
}
