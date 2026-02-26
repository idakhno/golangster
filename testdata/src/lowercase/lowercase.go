package lowercase

import (
	"log"
	"log/slog"
)

func bad() {
	slog.Info("Starting server on port 8080")     // want `log message must start with a lowercase letter`
	slog.Error("Failed to connect to database")   // want `log message must start with a lowercase letter`
	slog.Debug("Request received")                // want `log message must start with a lowercase letter`
	slog.Warn("Cache miss detected")              // want `log message must start with a lowercase letter`
	log.Print("Application started")              // want `log message must start with a lowercase letter`
	log.Printf("Server listening on %s", ":8080") // want `log message must start with a lowercase letter`
}

func good() {
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Debug("request received")
	log.Print("application started")
	log.Printf("server listening on %s", ":8080")
}

func withLogger() {
	logger := slog.Default()
	logger.Info("Starting service") // want `log message must start with a lowercase letter`
	logger.Info("starting service") // OK
}
