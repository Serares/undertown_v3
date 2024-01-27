package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestAuthorization(t *testing.T) {
	jwt := "ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxiV0ZwYkNJNkluUmxjM1JBWlcxaGFXd3VZMjl0SWl3aWRYTmxja2xrSWpvaVpEbGlORGhrTTJNdFpETm1aaTAwT0dNMkxUaGxOVGN0TjJGaFpXWTJNRGRoTkRZeUlpd2lhWE5oWkcxcGJpSTZkSEoxWlN3aWFYTnpjM0lpT21aaGJITmxmUS5iSW15eGlRd25rblRyV0Eyby0yLWMwU0N0S3liMVNIN1NrRlNvdXItRWdN"
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
