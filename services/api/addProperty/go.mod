module github.com/Serares/undertown_v3/services/api/addProperty

go 1.21.4

require github.com/google/uuid v1.5.0

require (
	github.com/Serares/undertown_v3/repositories/repository v0.0.0-00010101000000-000000000000
	github.com/akrylysov/algnhsa v1.1.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)

require (
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/aws/aws-lambda-go v1.43.0 // indirect
)

replace github.com/Serares/undertown_v3/repositories/repository => ../../../repositories/repository

replace github.com/Serares/undertown_v3/utils => ../../../utils
