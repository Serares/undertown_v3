module github.com/Serares/ssr/admin

go 1.21.6

require (
	github.com/Serares/undertown_v3/ssr/includes v0.0.0-00010101000000-000000000000
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/a-h/templ v0.2.543
	github.com/akrylysov/algnhsa v1.1.0
	github.com/joho/godotenv v1.5.1
)

require github.com/aws/aws-lambda-go v1.43.0 // indirect

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
)

replace github.com/Serares/undertown_v3/ssr/includes => ../includes

replace github.com/Serares/undertown_v3/utils => ../../../utils

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository
