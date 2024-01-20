package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestAuthorization(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxMjMxMjMyMTMiLCJlbWFpbCI6InJlZ2lzdHJhdG9yQGVtYWlsLmNvbSIsImlzYWRtaW4iOnRydWV9.kwUih83NOKaBC9AA2GB3HjFXBqWAsdhMl5pfz7tNOQ0"
	// create a struct of [string]string
	headersStruct := make(map[string]string)
	headersStruct["Authorization"] = jwt
	os.Setenv("JWT_SECRET", "myjwtsecret")
	// request event
	event := events.APIGatewayCustomAuthorizerRequestTypeRequest{
		Headers: headersStruct,
	}
	t.Run("Testing the bearer token", func(t *testing.T) {
		_, err := Handler(context.Background(), event)
		if err != nil {
			t.Fatalf("error calling the handler %v", err)
		}

		// fmt.Printf("Event %w", event)
		// TODO check the returned event
	})

}
