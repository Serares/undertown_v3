package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/ssr/admin/handlers"
	"github.com/Serares/ssr/admin/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4031"
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewAdminClient(log)
	loginService := service.NewLoginService(log, client)
	submitService := service.NewSubmitService(log, client)
	m := http.NewServeMux()
	loginHanlder := handlers.NewLoginHandler(log, loginService)
	submitHandler := handlers.NewSubmitHandler(log, submitService)
	// This is not advised to use in prod
	m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	m.Handle("/login/", loginHanlder)
	m.Handle("/submit/", submitHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      m,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
