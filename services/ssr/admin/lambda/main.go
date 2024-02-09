package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Serares/ssr/admin/handlers"
	"github.com/Serares/ssr/admin/middleware"
	"github.com/Serares/ssr/admin/service"
	"github.com/akrylysov/algnhsa"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewAdminClient(log)

	loginService := service.NewLoginService(log, client)
	submitService := service.NewSubmitService(log, client)
	listingService := service.NewListingService(log, client)
	editService := service.NewEditService(log, client)
	deleteService := service.NewDeleteService(log, client)

	m := http.NewServeMux()

	loginHanlder := handlers.NewLoginHandler(log, loginService)
	submitHandler := handlers.NewSubmitHandler(log, submitService)
	listingsHandler := handlers.NewListingsHandler(log, listingService)
	editHandler := handlers.NewEditHandler(log, editService)
	deleteHandler := handlers.NewDeleteHandler(log, deleteService)

	// This is not advised to use in prod
	m.Handle("/login/", loginHanlder)
	m.Handle("/login", loginHanlder)
	m.Handle("/submit/", middleware.NewMiddleware(submitHandler, middleware.WithSecure(false)))
	m.Handle("/submit", middleware.NewMiddleware(submitHandler, middleware.WithSecure(false)))
	m.Handle("/edit", middleware.NewMiddleware(editHandler, middleware.WithSecure(false)))
	m.Handle("/edit/", middleware.NewMiddleware(editHandler, middleware.WithSecure(false)))
	m.Handle("/list", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))
	m.Handle("/list/", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))
	m.Handle("/delete", middleware.NewMiddleware(deleteHandler, middleware.WithSecure(false)))
	m.Handle("/delete/", middleware.NewMiddleware(deleteHandler, middleware.WithSecure(false)))
	algnhsa.ListenAndServe(m, nil)
}
