package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sub2api/sub2api/handler"
)

const (
	defaultPort    = 8080
	defaultHost    = "0.0.0.0"
	appName        = "sub2api"
	appVersion     = "1.0.0"
)

func main() {
	var (
		port    int
		host    string
		version bool
	)

	flag.IntVar(&port, "port", getEnvInt("PORT", defaultPort), "Port to listen on")
	flag.StringVar(&host, "host", getEnv("HOST", defaultHost), "Host to bind to")
	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.Parse()

	if version {
		fmt.Printf("%s version %s\n", appName, appVersion)
		os.Exit(0)
	}

	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/", handler.Index)
	mux.HandleFunc("/sub", handler.Subscription)
	mux.HandleFunc("/health", handler.Health)

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("[%s] Starting server on %s", appName, addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[%s] Server failed: %v", appName, err)
	}
}

// getEnv returns the value of an environment variable or a default value.
func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

// getEnvInt returns the integer value of an environment variable or a default value.
func getEnvInt(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}
