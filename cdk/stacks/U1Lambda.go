package stacks

import (
	"cdk/utils"
	"net/http"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CrudEndpoints int

const (
	AddProperty CrudEndpoints = iota
	GetProperties
	GetProperty
)

func (ce CrudEndpoints) String() string {
	return [...]string{"addProperty", "getProperties", "getProperty"}[ce]
}

type U1LambdaProps struct {
	awscdk.StackProps
	Env          string
	AssetsBucket awss3.Bucket
}

type U1LambdaStack struct {
	Lambdas []IntegrationLambda
	Stack   awscdk.Stack
}

// This stack is used for CRUD operations
func U1Lambda(scope constructs.Construct, id string, props *U1LambdaProps) U1LambdaStack {
	var sprops awscdk.StackProps
	var lambdas []IntegrationLambda
	var assetsBucketArn string
	lambdasEnvVars := map[string]*string{
		// ❗
		// TODO try to obfuscate somehow the values
		// don't store them in plain text
		// store them as an encrypted string?
		// how to decrypt them
		"DB_HOST":        jsii.String(os.Getenv("DB_HOST")),
		"DB_NAME":        jsii.String(os.Getenv("DB_NAME")),
		"DB_PROTOCOL":    jsii.String(os.Getenv("DB_PROTOCOL")),
		"TURSO_DB_TOKEN": jsii.String(os.Getenv("TURSO_DB_TOKEN")),
	}
	if props.AssetsBucket.BucketArn() != nil {
		assetsBucketArn = *props.AssetsBucket.BucketName()
	}
	addPropertyEnv := map[string]*string{
		// ❗
		// TODO try to obfuscate somehow the values
		// don't store them in plain text
		// store them as an encrypted string?
		// how to decrypt them
		"DB_HOST":            jsii.String(os.Getenv("DB_HOST")),
		"DB_NAME":            jsii.String(os.Getenv("DB_NAME")),
		"DB_PROTOCOL":        jsii.String(os.Getenv("DB_PROTOCOL")),
		"TURSO_DB_TOKEN":     jsii.String(os.Getenv("TURSO_DB_TOKEN")),
		"ASSETS_BUCKET_NAME": jsii.String(assetsBucketArn),
	}
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	lambdaRole := utils.CreateLambdaBasicRole(stack, "lambdaBasicRoleU1", props.Env)
	s3BucketAccessRole := utils.CreateLambdaBasicRole(stack, "s3fullaccesslambdarole", props.Env)
	s3BucketAccessRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")))
	// Create the cwebp layer for image optimization
	cwebpLayer := awslambda.NewLayerVersion(stack, jsii.String("cwebp-layer"), &awslambda.LayerVersionProps{
		Code: awslambda.AssetCode_FromAsset(jsii.String("../layers/cwebp/cwebp-layer.zip"), nil),
		CompatibleRuntimes: &[]awslambda.Runtime{
			awslambda.Runtime_PROVIDED_AL2(),
			// Add other compatible runtimes if needed
		},
		LayerVersionName: jsii.String("cwebp-layer"),
		Description:      jsii.String("Layer with cwebp binary"),
	})

	addProperty := awslambdago.NewGoFunction(stack, jsii.Sprintf("AddProperty-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/addProperty/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &addPropertyEnv,
		Role:         s3BucketAccessRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(3 * 60)),
		Layers:       &[]awslambda.ILayerVersion{cwebpLayer},
	})

	getProperties := awslambdago.NewGoFunction(stack, jsii.Sprintf("GetProperties-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/getProperties/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &lambdasEnvVars,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	getProperty := awslambdago.NewGoFunction(stack, jsii.Sprintf("GetProperty-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/getProperty/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &lambdasEnvVars,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &addProperty,
		path:       AddProperty.String(),
		method:     http.MethodPost,
		authorizer: CRUDAuthorizer.String(),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &getProperties,
		path:       GetProperties.String(),
		method:     http.MethodGet,
		authorizer: CRUDAuthorizer.String(),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &getProperty,
		path:       GetProperty.String(),
		method:     http.MethodGet,
		authorizer: CRUDAuthorizer.String(),
	})

	return U1LambdaStack{
		Lambdas: lambdas,
		Stack:   stack,
	}
}
