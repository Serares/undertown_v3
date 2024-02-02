package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/login/handlers"
	"github.com/Serares/undertown_v3/services/api/login/service"
	"github.com/joho/godotenv"
)

// 🔎
// logging in the user
func main() {
	err := godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3033"
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userRepo, err := repository.NewUsersRepository()
	if err != nil {
		log.Error("error on initializing the db")
	}
	service := service.NewLoginService(log, userRepo)
	h := handlers.NewLoginHandler(log, &service)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      h,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
