package main

import (
	"flag"
	"net/http"
	"os"

	"golang.org/x/exp/slog"
)

var (
	buildString = "unknwown"
)

func initLogger(verbose bool) *slog.Logger {
	lvl := new(slog.LevelVar) // Info by default
	if verbose {
		lvl.Set(slog.LevelDebug)
	}

	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvl}))
}

func init() {
	flag.StringVar(&serverAddr, "addr", ":8080", "Address to run the HTTP server on")
	flag.StringVar(&mattermostURL, "webhook-url", "http://mattermost.internal", "Mattermost webhook URL")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	logger = initLogger(verbose)
}

func main() {
	flag.Parse()

	// Check if the webhook URL is provided via environment variable
	if envURL := os.Getenv("WEBHOOKURL"); envURL != "" {
		mattermostURL = envURL
	}

	if mattermostURL == "" {
		logger.Error("Mattermost webhook URL is not provided. Use the -webhook-url flag or set the WEBHOOKURL environment variable.")
		os.Exit(1)
	}

	// Define handlers.
	http.HandleFunc("/alert", handleAlert)
	http.HandleFunc("/liveness", handleLiveness)
	http.HandleFunc("/readiness", handleReadiness)

	logger.Info("Starting server", "addr", serverAddr, "version", buildString)
	if err := (http.ListenAndServe(serverAddr, nil)); err != nil {
		logger.Error("Error starting server", "err", err)
		os.Exit(1)
	}
}
