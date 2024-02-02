package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/addProperty/handler"
	"github.com/Serares/undertown_v3/services/api/addProperty/service"
)

// 🚀
func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3031"
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := repository.NewPropertiesRepo()
	ss := service.NewSubmitService(log, db)

	if err != nil {
		log.Error("error on initializing the db")
	}

	hh := handler.New(log, ss)
	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      hh,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
