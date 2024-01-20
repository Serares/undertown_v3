package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Serares/ssr/homepage/handlers"
	"github.com/Serares/ssr/homepage/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewHomeClient(log)
	homeService := service.NewHomeService(log, client)
	propertiesService := service.NewPropertiesService(log, client)

	m := http.NewServeMux()
	propertiesHandler := handlers.NewPropertiesHandler(log, *propertiesService)
	homeHandler := handlers.NewHomeHandler(log, *homeService)

	// This is not advised to use in prod
	m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	m.Handle("/chirii/", propertiesHandler)
	m.Handle("/vanzari/", propertiesHandler)
	m.Handle("/", homeHandler)

	algnhsa.ListenAndServe(m, nil)
}
