package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Serares/ssr/admin/handlers"
	"github.com/Serares/ssr/admin/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewAdminClient(log)

	loginService := service.NewLoginService(log, client)
	submitService := service.NewSubmitService(log, client)
	m := http.NewServeMux()
	loginHanlder := handlers.NewLoginHandler(log, loginService)
	submitHandler := handlers.NewSubmitHandler(log, submitService)
	m.Handle("/login/", loginHanlder)
	m.Handle("/login", loginHanlder)
	m.Handle("/submit/", submitHandler)
	m.Handle("/submit", submitHandler)
	algnhsa.ListenAndServe(m, nil)
}
