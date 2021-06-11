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

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/emar-kar/urlshortener/cmd/server"
	"github.com/emar-kar/urlshortener/internal/redis"
	"github.com/emar-kar/urlshortener/pkg/handler"
	"github.com/emar-kar/urlshortener/pkg/service"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/report.log",
		MaxBackups: 2,
		MaxAge:     1, //days
	})

	srv := new(server.Server)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer close(done)

	rdb := redis.NewDB()
	if _, err := rdb.Client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("cannot logging to redis: %s", err)
	}
	defer rdb.Client.Close()

	services := service.NewService(rdb)
	handlers := handler.NewHandler(services)

	go func() {
		err := srv.Run(
			"",
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
