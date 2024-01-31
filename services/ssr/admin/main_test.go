package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Serares/ssr/admin/handlers"
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
	// var requestBody bytes.Buffer
	// multipartWriter := multipart.NewWriter(&requestBody)
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	client := service.NewAdminClient(log)
	submitService := service.NewSubmitService(log, client)
	hh := handlers.NewSubmitHandler(log, submitService)
	loginService := service.NewLoginService(log, client)
	loginHandler := handlers.NewLoginHandler(log, loginService)
	mux := http.NewServeMux()
	mux.Handle("/submit", hh)
	mux.Handle("/login", loginHandler)

	ts := httptest.NewServer(mux)

	// propertyField, err := multipartWriter.CreateFormField("property")
	// if err != nil {
	// 	t.Errorf("error adding title field to the form")
	// }

	// _, err = propertyField.Write([]byte("mock_title"))
	// if err != nil {
	// 	t.Error("error adding mock value to title")
	// }

	// imageWriter, err := multipartWriter.CreateFormFile("images", mockImage.Name())
	// if err != nil {
	// 	t.Error("error trying to create the file for form")
	// }
	// _, err = io.Copy(imageWriter, mockImage)
	// if err != nil {
	// 	t.Error("error trying to copy the image to the form file writer")
	// }
	// image2Writer, _ := multipartWriter.CreateFormFile("images", mockImage2.Name())
	// _, err = io.Copy(image2Writer, mockImage2)

	// multipartWriter.Close()
	// request, err := http.NewRequest(http.MethodPost, ts.URL+"/submit", &requestBody)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// response, err := (&http.Client{}).Do(request)
	// if err != nil {
	// 	t.Errorf("error sending the request %v", err)
	// }
	// defer response.Body.Close()
	return ts.URL, func() {
		log.Info("Shutting down the test server")
		ts.Close()
	}
}

func TestPost(t *testing.T) {
	url, cleanup := setupAPI(t)
	fmt.Println("The server url", url)
	// var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	defer cleanup()
}
