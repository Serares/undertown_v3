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

	// This is not advised to use in prod
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	mux.Handle("/contact", http.StripPrefix("/contact", contactHandler))
	mux.Handle("/contact/", http.StripPrefix("/contact/", contactHandler))
	// m.Handle("/property/{ID}", http.StripPrefix("/property/", contactHandler))
	mux.Handle("/", homeHandler)

	algnhsa.ListenAndServe(mux, nil)
}
