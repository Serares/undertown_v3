package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/ssr/admin/handlers"
	"github.com/Serares/ssr/admin/middleware"
	"github.com/Serares/ssr/admin/service"
	"github.com/Serares/ssr/admin/types"
	"github.com/Serares/ssr/admin/views"
	"github.com/Serares/undertown_v3/ssr/includes/components"
	"github.com/a-h/templ"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Error(
			"error trying to load the lambda context",
			"error", err,
		)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)
	s3PresignClient := s3.NewPresignClient(s3Client)

	loginService := service.NewLoginService(log, client)
	submitService := service.NewSubmitService(
		log,
		client,
		sqsClient,
		s3Client,
	)
	listingService := service.NewListingService(log, client)
	editService := service.NewEditService(
		log,
		client,
		sqsClient,
		s3Client,
	)
	deleteService := service.NewDeleteService(log, client)

	m := http.NewServeMux()

	loginHanlder := handlers.NewLoginHandler(log, loginService)
	submitHandler := handlers.NewSubmitHandler(log, submitService)
	listingsHandler := handlers.NewListingsHandler(log, listingService)
	editHandler := handlers.NewEditHandler(log, editService)
	deleteHandler := handlers.NewDeleteHandler(log, deleteService)
	presignHandler := handlers.NewPresignedS3Handler(log, s3PresignClient)
	// This is not advised to use in prod
	m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	m.Handle("/login/", loginHanlder)
	m.Handle("/login", loginHanlder)
	m.Handle("/submit/", middleware.NewMiddleware(submitHandler, middleware.WithSecure(false)))
	m.Handle("/submit", middleware.NewMiddleware(submitHandler, middleware.WithSecure(false)))
	m.Handle("/edit", middleware.NewMiddleware(editHandler, middleware.WithSecure(false)))
	m.Handle("/edit/", middleware.NewMiddleware(editHandler, middleware.WithSecure(false)))
	m.Handle("/delete", middleware.NewMiddleware(deleteHandler, middleware.WithSecure(false)))
	m.Handle("/delete/", middleware.NewMiddleware(deleteHandler, middleware.WithSecure(false)))
	m.Handle("/presign", middleware.NewMiddleware(presignHandler, middleware.WithSecure(false)))
	m.Handle("/", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))
	m.Handle("/test", templ.Handler(views.Dztest(types.BasicIncludes{Scripts: components.Scripts()}, types.SubmitProps{})))

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%s", port),
		Handler:      m,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	fmt.Printf("Listening on %v\n", server.Addr)
	server.ListenAndServe()
}
