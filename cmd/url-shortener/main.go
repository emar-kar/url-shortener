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

	"github.com/emar-kar/url-shortener/cmd/url-shortener/server"
	"github.com/emar-kar/url-shortener/pkg/handler"
)

// TODO: add config?
func main() {
	srv := new(server.Server)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	handlers := handler.NewHandler()

	go func() {
		err := srv.Run(
			"localhost",
			"8080",
			handlers.InitRoutes(""),
		)
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server closed")
			return
		}
		if err != nil {
			log.Printf("server start failure: %s", err)
			done <- os.Interrupt
		}
	}()

	<-done
	fmt.Println()

	maxGracefulShutdownTime := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), maxGracefulShutdownTime)
	defer cancel()

	if err := srv.ShutDown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
		return
	}

	log.Println("server exited properly")
}
