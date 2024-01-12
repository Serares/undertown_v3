package stacks

import (
	"os"

	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

var JwtSecret = os.Getenv("JWT_SECRET")

// Strip the binary, and remove the deprecated Lambda SDK RPC code for performance.
// These options are not required, but make cold start faster.
var BundlingOptions = &awslambdago.BundlingOptions{
	GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w" -tags lambda.norpc`)},
}
