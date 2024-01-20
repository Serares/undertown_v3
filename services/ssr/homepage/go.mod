module github.com/Serares/ssr/homepage

go 1.21.4

require (
	github.com/a-h/templ v0.2.476
	github.com/akrylysov/algnhsa v1.0.0
	github.com/joho/godotenv v1.5.1
)

require github.com/golang-jwt/jwt/v5 v5.2.0 // indirect

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/aws/aws-lambda-go v1.37.0 // indirect
)

replace github.com/Serares/undertown_v3/utils => ../../../utils

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository
