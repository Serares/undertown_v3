package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	"github.com/Serares/undertown_v3/services/api/login/handlers"
	"github.com/Serares/undertown_v3/services/api/login/service"
	"github.com/joho/godotenv"
)

// 🔎
// logging in the user
func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3033"
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbUrl, err := utils.CreatePsqlUrl(context.Background(), log)
	if err != nil {
		log.Error("error on creating the connection string")
	}
	userRepo, err := repository.NewUsersRepository(dbUrl)
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