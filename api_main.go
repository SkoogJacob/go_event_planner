package main

import (
	"context"
	"errors"
	"event_planner_api/event_db"
	"event_planner_api/routing"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	event_db.InitDb("./data/events.db")
	defer event_db.CloseDb()

	server := routing.MakeServer(":9000")

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	<-ctx.Done()
	log.Println("timeout of 5 seconds")
	log.Println("Server exiting")
}
