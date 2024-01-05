package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	"github.com/Serares/undertown_v3/services/api/register/handlers"
	"github.com/Serares/undertown_v3/services/api/register/service"
	"github.com/joho/godotenv"
)

// 🪪
// Register users lambda
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3032"
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	dbUrl := utils.CreatePsqlUrl()
	userRepo, err := repository.NewUsersRepository(dbUrl)
	if err != nil {
		log.Error("error on initializing the db")
	}
	service := service.NewRegisterService(log, userRepo)
	h := handlers.NewRegisterHandler(log, service)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      h,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
