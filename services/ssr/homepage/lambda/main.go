package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Serares/ssr/homepage/handlers"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mux := http.NewServeMux()

	contactHandler := handlers.NewContactHandler(log)
	homeHandler := handlers.NewHomeHandler(log)

	// you might not need to create a route for the static assets
	// the route is created in the cdk
	mux.Handle("/contact", http.StripPrefix("/contact", contactHandler))
	mux.Handle("/contact/", http.StripPrefix("/contact/", contactHandler))
	mux.Handle("/property/", http.StripPrefix("/property/", contactHandler))
	mux.Handle("/", homeHandler)

	algnhsa.ListenAndServe(mux, nil)
}
