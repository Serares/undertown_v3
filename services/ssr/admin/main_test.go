package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Serares/ssr/admin/handlers"
	"github.com/Serares/ssr/admin/middleware"
	"github.com/Serares/ssr/admin/service"
	"github.com/joho/godotenv"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()
	err := godotenv.Load(".env.local")
	if err != nil {
		t.Error("error loading the .env file")
	}
	mockImage, err := os.Open("testdata/mockImage.jpg")
	if err != nil {
		t.Error("error reading the mock image file")
	}
	mockImage2, err := os.Open("testdata/mockImage2.png")
	if err != nil {
		t.Error("error reading the mock image file")
	}
	defer mockImage2.Close()
	defer mockImage.Close()
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewAdminClient(log)

	loginService := service.NewLoginService(log, client)
	submitService := service.NewSubmitService(log, client)
	listingService := service.NewListingService(log, client)
	editService := service.NewEditService(log, client)

	m := http.NewServeMux()

	loginHanlder := handlers.NewLoginHandler(log, loginService)
	submitHandler := handlers.NewSubmitHandler(log, submitService)
	listingsHandler := handlers.NewListingsHandler(log, listingService)
	editHandler := handlers.NewEditHandler(log, editService)
	// This is not advised to use in prod
	m.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))
	m.Handle("/login/", loginHanlder)
	m.Handle("/login", loginHanlder)
	m.Handle("/submit/", middleware.NewMiddleware(submitHandler, middleware.WithSecure(false)))
	m.Handle("/submit", middleware.NewMiddleware(submitHandler, middleware.WithSecure(false)))
	m.Handle("/edit", middleware.NewMiddleware(editHandler, middleware.WithSecure(false)))
	m.Handle("/edit/", middleware.NewMiddleware(editHandler, middleware.WithSecure(false)))
	m.Handle("/list", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))
	m.Handle("/list/", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))
	m.Handle("/delete", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))
	m.Handle("/delete/", middleware.NewMiddleware(listingsHandler, middleware.WithSecure(false)))

	ts := httptest.NewServer(m)

	titleField, err := multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("error adding title field to the form")
	}

	_, err = titleField.Write([]byte("mock_title"))
	if err != nil {
		t.Error("error adding mock value to title")
	}
	addressField, err := multipartWriter.CreateFormField("property_address")
	if err != nil {
		t.Errorf("error adding property_address field to the form")
	}

	_, err = addressField.Write([]byte("mock_address"))
	if err != nil {
		t.Error("error adding mock value to address")
	}
	isFeaturedField, err := multipartWriter.CreateFormField("is_featured")
	if err != nil {
		t.Errorf("error adding is_featured field to the form")
	}

	_, err = isFeaturedField.Write([]byte("on"))
	if err != nil {
		t.Error("error adding mock value to is_featured")
	}

	imageWriter, err := multipartWriter.CreateFormFile("images", mockImage.Name())
	if err != nil {
		t.Error("error trying to create the file for form")
	}
	_, err = io.Copy(imageWriter, mockImage)
	if err != nil {
		t.Error("error trying to copy the image to the form file writer")
	}
	image2Writer, _ := multipartWriter.CreateFormFile("images", mockImage2.Name())
	_, err = io.Copy(image2Writer, mockImage2)

	multipartWriter.Close()
	request, err := http.NewRequest(http.MethodPost, ts.URL+"/submit", &requestBody)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	response, err := (&http.Client{}).Do(request)
	if err != nil {
		t.Errorf("error sending the request %v", err)
	}
	defer response.Body.Close()
	return ts.URL, func() {
		log.Info("Shutting down the test server")
		ts.Close()
	}
}

// This is an e2e test
func TestPost(t *testing.T) {
	url, cleanup := setupAPI(t)
	fmt.Println("The server url", url)
	// var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	defer cleanup()
}
