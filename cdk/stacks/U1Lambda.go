package stacks

import (
	"cdk/utils"
	"net/http"
	"os"

	"github.com/Serares/undertown_v3/utils/env"
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CrudEndpoints int

const (
	AddProperty CrudEndpoints = iota
	GetProperties
	GetProperty
	DeleteProperty
)

func (ce CrudEndpoints) String() string {
	return [...]string{"addProperty", "getProperties", "getProperty", "deleteProperty"}[ce]
}

type U1LambdaProps struct {
	awscdk.StackProps
	Env                        string
	DeleteProcessedImagesQueue awssqs.Queue
}

type U1LambdaStack struct {
	Lambdas []IntegrationLambda
	Stack   awscdk.Stack
}

// This stack is used for CRUD operations
func U1Lambda(scope constructs.Construct, id string, props *U1LambdaProps) U1LambdaStack {
	var sprops awscdk.StackProps
	var lambdas []IntegrationLambda
	var deleteProcessedImagesQueueUrl string

	lambdasEnvVars := map[string]*string{
		// ❗
		// TODO try to obfuscate somehow the values
		// don't store them in plain text
		// store them as an encrypted string?
		// how to decrypt them
		env.DB_HOST:        jsii.String(os.Getenv(env.DB_HOST)),
		env.DB_NAME:        jsii.String(os.Getenv(env.DB_NAME)),
		env.DB_PROTOCOL:    jsii.String(os.Getenv(env.DB_PROTOCOL)),
		env.TURSO_DB_TOKEN: jsii.String(os.Getenv(env.TURSO_DB_TOKEN)),
	}

	if props.DeleteProcessedImagesQueue.QueueUrl() != nil {
		deleteProcessedImagesQueueUrl = *props.DeleteProcessedImagesQueue.QueueUrl()
	}

	deletePropertyLambdaEnv := map[string]*string{
		// ❗
		// TODO try to obfuscate somehow the values
		// don't store them in plain text
		// store them as an encrypted string?
		// how to decrypt them
		env.DB_HOST:        jsii.String(os.Getenv(env.DB_HOST)),
		env.DB_NAME:        jsii.String(os.Getenv(env.DB_NAME)),
		env.DB_PROTOCOL:    jsii.String(os.Getenv(env.DB_PROTOCOL)),
		env.TURSO_DB_TOKEN: jsii.String(os.Getenv(env.TURSO_DB_TOKEN)),
		env.SQS_DELETE_PROCESSED_IMAGES_QUEUE_URL: jsii.String(deleteProcessedImagesQueueUrl), // used to dispatch the names of images that needs to be deleted
	}

	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)
	lambdaRole := utils.CreateLambdaBasicRole(stack, "lambdaBasicRoleU1", props.Env)
	s3BucketAccessRole := utils.CreateLambdaBasicRole(stack, "s3fullaccesslambdarole", props.Env)
	s3BucketAccessRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3FullAccess")))

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

	deleteProperty := awslambdago.NewGoFunction(stack, jsii.Sprintf("DeleteProperty-%s", props.Env), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/deleteProperty/lambda"),
		Bundling:     BundlingOptions,
		Environment:  &deletePropertyLambdaEnv,
		Role:         lambdaRole,
		Timeout:      awscdk.Duration_Seconds(jsii.Number(30)),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &getProperties,
		path:       GetProperties.String(),
		method:     []string{http.MethodGet},
		authorizer: "",
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &getProperty,
		path:       GetProperty.String(),
		method:     []string{http.MethodGet},
		authorizer: CRUDAuthorizer.String(),
	})

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &deleteProperty,
		path:       DeleteProperty.String(),
		method:     []string{http.MethodDelete},
		authorizer: "",
	})

	// Delete property lambda will dispatch a delete images event
	props.DeleteProcessedImagesQueue.GrantSendMessages(deleteProperty)

	return U1LambdaStack{
		Lambdas: lambdas,
		Stack:   stack,
	}
}
