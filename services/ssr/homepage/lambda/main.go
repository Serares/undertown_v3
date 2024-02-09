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
	client := service.NewClient(log)
	homeService := service.NewHomeService(log, client)
	propertiesService := service.NewPropertiesService(log, client)
	singlePropertyService := service.NewPropertyService(log, client)

	m := http.NewServeMux()
	propertiesHandler := handlers.NewPropertiesHandler(log, *propertiesService, singlePropertyService)
	defaultHandler := handlers.NewDefaultHandler(log, homeService)

	// This is not advised to use in prod
	// Lambda doesn't need to handle the route for assets
	// the route is handled by cloudfront
	// m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	m.Handle("/chirii/", propertiesHandler)
	m.Handle("/vanzari/", propertiesHandler)
	m.Handle("/", defaultHandler)

	algnhsa.ListenAndServe(m, nil)
}
