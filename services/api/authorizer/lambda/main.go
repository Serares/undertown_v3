package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func Handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.Headers["Authorization"]
	err := godotenv.Load("../.env.local")
	secret := os.Getenv("JWT_SECRET")
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.Info("The event object", "event", event)
	if err != nil {
		log.Debug("env file is used only for local testing")
	}
	// Parse the JWT
	claims := utils.JWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		// Return a policy document that denies access
		return generatePolicy("user", "Deny", event.MethodArn), nil
	}
	return generatePolicy("user", "Deny", event.MethodArn), nil
}

func main() {
	lambda.Start(Handler)
}

// generatePolicy generates an IAM policy document
func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}
	authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   effect,
				Resource: []string{resource},
			},
		},
	}

	return authResponse
}
