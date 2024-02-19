package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Serares/ssr/homepage/handlers"
	"github.com/Serares/ssr/homepage/service"
	"github.com/joho/godotenv"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()
	err := godotenv.Load(".env.dev")
	if err != nil {
		t.Error("error loading the .env file")
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewClient(log)
	homeService := service.NewHomeService(log, client)
	propertiesService := service.NewPropertiesService(log, client)
	singlePropertyService := service.NewPropertyService(log, client)

	m := http.NewServeMux()
	propertiesHandler := handlers.NewPropertiesHandler(log, *propertiesService, singlePropertyService)
	aboutHandler := handlers.NewAboutHandler(log)
	contactHandler := handlers.NewContactHandler(log)
	defaultHandler := handlers.NewDefaultHandler(log, homeService)

	// This is not advised to use in prod
	m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	m.Handle("/chirii/", propertiesHandler)
	m.Handle("/vanzari/", propertiesHandler)
	m.Handle("/about", aboutHandler)
	m.Handle("/contact", contactHandler)
	m.Handle("/", defaultHandler)

	ts := httptest.NewServer(m)

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
