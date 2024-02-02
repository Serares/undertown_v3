module github.com/Serares/ssr/homepage

go 1.21.6

require (
	github.com/a-h/templ v0.2.543
	github.com/akrylysov/algnhsa v1.1.0
	github.com/joho/godotenv v1.5.1
)

require github.com/golang-jwt/jwt/v5 v5.2.0

require github.com/stretchr/testify v1.8.4 // indirect

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0
	github.com/Serares/undertown_v3/ssr/includes v0.0.0
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/aws/aws-lambda-go v1.44.0 // indirect
)

replace github.com/Serares/undertown_v3/utils => ../../../utils

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository

replace github.com/Serares/undertown_v3/ssr/includes => ../includes
