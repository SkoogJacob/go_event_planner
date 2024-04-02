package main

import (
	"context"
	"errors"
	"event_planner_api/event_db"
	"event_planner_api/routing"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func argparse() (db_file string, host string) {
	flag.StringVar(&db_file, "db-file", "./data/events.db",
		`The path to the sqlite database file. The path should contain the file name and
		the file name should end with '.db'`)
	flag.StringVar(&host, "host", ":9000",
		`The desired IP address + port for the server to listen on. If only port is given
		the server will listen on localhost:9000`)
	if !strings.HasSuffix(db_file, ".db") {
		log.Fatalf("db file should have .db suffix, got %s", db_file)
	}
	matcher, _ := regexp.Compile(`((\d{1-3}){4})?:\d{4,5}$`)
	if !matcher.MatchString(host) {
		log.Fatalf("the given host address does not look like a valid ipv4+port address")
	}
	return db_file, host
}

func main() {
	db_file, host := argparse()
	event_db.InitDb(db_file)
	defer event_db.CloseDb()

	server := routing.MakeServer(host)

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
