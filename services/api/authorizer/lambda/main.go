package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/Serares/undertown_v3/utils"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt/v5"
)

func Handler(ctx context.Context, event events.APIGatewayCustomAuthorizerRequestTypeRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := event.Headers["Authorization"]
	// err := godotenv.Load("../.env.local")
	secret := os.Getenv("JWT_SECRET")
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil)).WithGroup("CRUD AUTHORIZER")
	log.Info("The event object", "event", event)
	// if err != nil {
	// 	log.Debug("env file is used only for local testing")
	// }
	// If the token exists it has to be base64 decoded
	// token, err := base64.RawStdEncoding.DecodeString(base64Token)
	// if err != nil {
	// 	log.Error("error trying to base64 decode the Authorization token", err)
	// }
	// Parse the JWT
	claims := utils.JWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(string(token), &claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}
		return []byte(secret), nil
	})

	if err != nil || !parsedToken.Valid {
		// Return a policy document that denies access
		log.Error("the token is invalid", "error", err)
		return generatePolicy("user", "Deny", event.MethodArn, ""), nil
	}
	generatedPolicy := generatePolicy("user", "Allow", event.MethodArn, claims.UserId)
	log.Info("Generated policy", "policy", generatedPolicy)
	return generatedPolicy, nil
}

func main() {
	lambda.Start(Handler)
}

// generatePolicy generates an IAM policy document
func generatePolicy(principalID, effect, resource, userId string) events.APIGatewayCustomAuthorizerResponse {
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
	authResponse.Context = map[string]interface{}{
		"userId": userId,
	}
	return authResponse
}
