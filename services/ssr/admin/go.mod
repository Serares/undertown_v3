module github.com/Serares/ssr/admin

go 1.21.6

require (
	github.com/Serares/undertown_v3/ssr/includes v0.0.0-00010101000000-000000000000
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/a-h/templ v0.2.543
	github.com/akrylysov/algnhsa v1.1.0
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.6
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.1
	github.com/aws/aws-sdk-go-v2/service/sqs v1.29.7
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/aws/aws-lambda-go v1.46.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.7 // indirect
	github.com/aws/smithy-go v1.19.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
)

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
)

replace github.com/Serares/undertown_v3/ssr/includes => ../includes

replace github.com/Serares/undertown_v3/utils => ../../../utils

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository
