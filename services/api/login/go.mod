module github.com/Serares/undertown_v3/services/api/login

go 1.21.4

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.17.0
)

require (
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.5.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
)

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository

replace github.com/Serares/undertown_v3/utils => ../../../utils
