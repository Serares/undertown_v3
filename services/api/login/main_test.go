package main

import (
	"fmt"
	"log/slog"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/login/handlers"
	"github.com/Serares/undertown_v3/services/api/login/service"
	"github.com/joho/godotenv"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()
	err := godotenv.Load(".env.dev")
	if err != nil {
		t.Error("error loading the .env file")
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userRepo, err := repository.NewUsersRepository()
	if err != nil {
		log.Error("error on initializing the db")
	}
	service := service.NewLoginService(log, userRepo)
	h := handlers.NewLoginHandler(log, &service)

	ts := httptest.NewServer(h)

	return ts.URL, func() {
		log.Info("Shutting down the test server")
		ts.Close()
	}
}

// This is an e2e test
func TestPost(t *testing.T) {
	url, cleanup := setupAPI(t)
	fmt.Println("The server url", url)
	// var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	defer cleanup()
}
