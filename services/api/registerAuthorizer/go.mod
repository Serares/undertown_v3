module github.com/Serares/undertown_v3/services/api/registerAuthorizer

go 1.21.4

require (
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/aws/aws-lambda-go v1.44.0
	github.com/golang-jwt/jwt/v5 v5.2.0
)

replace github.com/Serares/undertown_v3/utils => ../../../utils
