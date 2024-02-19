package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/ssr/homepage/handlers"
	"github.com/Serares/ssr/homepage/service"
	"github.com/Serares/undertown_v3/utils/constants"
	"github.com/joho/godotenv"
)

// TODO handling routing in the same lambda might not be a good idea but you can try and see how it works
// if it's a bad idea then youll have to serve pages from different lambdas

// This is meant to be used on localhost only
// Entrypoint for deployment in inside lambda directory

func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4030"
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewClient(log)
	homeService := service.NewHomeService(log, client)
	propertiesService := service.NewPropertiesService(log, client)
	singlePropertyService := service.NewPropertyService(log, client)

	m := http.NewServeMux()
	propertiesHandler := handlers.NewPropertiesHandler(log, *propertiesService, singlePropertyService)
	defaultHandler := handlers.NewDefaultHandler(log, homeService)

	rentPath := fmt.Sprintf("/%s/", constants.TRANSACTION_TYPE_UI_RENT)
	sellPath := fmt.Sprintf("/%s/", constants.TRANSACTION_TYPE_UI_SELL)

	// This is not advised to use in prod
	m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	m.Handle(rentPath, propertiesHandler)
	m.Handle(sellPath, propertiesHandler)
	m.Handle("/", defaultHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      m,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
