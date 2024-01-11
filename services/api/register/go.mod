module github.com/Serares/undertown_v3/services/api/register

go 1.21.4

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0
	github.com/akrylysov/algnhsa v1.1.0
	github.com/google/uuid v1.5.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.17.0
)

require (
	github.com/aws/aws-lambda-go v1.43.0 // indirect
	github.com/aws/aws-sdk-go v1.49.17 // indirect
	github.com/aws/aws-sdk-go-v2 v1.24.1 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.26.3 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.7 // indirect
	github.com/aws/smithy-go v1.19.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

require (
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/lib/pq v1.10.9 // indirect
)

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository

replace github.com/Serares/undertown_v3/utils => ../../../utils
