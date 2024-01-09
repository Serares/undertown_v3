package stacks

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type ApiLambdaStackProps struct {
	awscdk.StackProps
	DbSecretArn string
}

type DbCfg struct {
	PSQL_USER     string
	PSQL_PASSWORD string
	PSQL_DB       string
	PSQL_HOST     string
	PSQL_PORT     int32
	PSQL_ENDPOINT string
}

func ApiLambdaStack(scope constructs.Construct, id string, props *ApiLambdaStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	var authorizationHeader = "Authorization"
	err := godotenv.Load(".env.dev")
	if err != nil {
		// handle error in the cdk stack
		panic(err)
	}
	// TODO create this secret in SSM
	var jwtSecret = os.Getenv("JWT_SECRET")
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Strip the binary, and remove the deprecated Lambda SDK RPC code for performance.
	// These options are not required, but make cold start faster.
	bundlingOptions := &awslambdago.BundlingOptions{
		GoBuildFlags: &[]*string{jsii.String(`-ldflags "-s -w" -tags lambda.norpc`)},
	}
	auroraHost := awscdk.Fn_ImportValue(jsii.String(DB_STACK_KEY_DB_HOST))
	// THIS IS THE
	auroraDbName := awscdk.Fn_ImportValue(jsii.String(DB_STACK_KEY_DB_NAME))
	auroraDbSecret := awscdk.Fn_ImportValue(jsii.String(DB_STACK_KEY_DB_SECRET))

	// Lambda role to access secrets manager
	apiLambdaRole := awsiam.NewRole(stack, jsii.String("SecretsManagerAccessRole"), &awsiam.RoleProps{
		// TODO think about if this needs to apply the least priviledge principle
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})
	apiLambdaRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")))
	apiLambdaRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("secretsmanager:GetSecretValue"),
		Resources: jsii.Strings(*auroraDbSecret),
	}))

	// API
	// Add a resource based policy to the API so that the homessr lambda can only call:
	registerLambda := awslambdago.NewGoFunction(stack, jsii.String("RegisterLambda"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/register/lambda"),
		Bundling:     bundlingOptions,
		Environment: &map[string]*string{
			"PSQL_HOST":     jsii.String(*auroraHost),
			"PSQL_USER":     jsii.String(""),
			"PSQL_PASSWORD": jsii.String(""),
			"PSQL_DB":       jsii.String(*auroraDbName),
			"PSQL_PORT":     jsii.String("5432"),
			// this is used to retreive the password and username
			"PSQL_SECRET_ARN": jsii.String(*auroraDbSecret),
		},
		Role:    apiLambdaRole,
		Timeout: awscdk.Duration_Seconds(jsii.Number(20)),
	})
	//Register Authorizer
	registerAuthorizerLambda := awslambdago.NewGoFunction(stack, jsii.String("RegisterAuthorizer"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/registerAuthorizer/lambda"),
		Bundling:     bundlingOptions,
		Environment: &map[string]*string{
			"JWT_SECRET": &jwtSecret,
		},
	})
	// API Authorizer lambda
	apiAuthorizerLambda := awslambdago.NewGoFunction(stack, jsii.String("APIAuthorizer"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/authorizer/lambda"),
		Bundling:     bundlingOptions,
		Environment: &map[string]*string{
			"JWT_SECRET": &jwtSecret,
		},
	})
	// Login Lambda
	loginLambda := awslambdago.NewGoFunction(stack, jsii.String("SPALogin"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/login/lambda"),
		Bundling:     bundlingOptions,
		Environment: &map[string]*string{
			"PSQL_HOST":     jsii.String(*auroraHost),
			"PSQL_USER":     jsii.String(""),
			"PSQL_PASSWORD": jsii.String(""),
			"PSQL_DB":       jsii.String(*auroraDbName),
			"PSQL_PORT":     jsii.String("5432"),
			// this is used to retreive the password and username
			"PSQL_SECRET_ARN": jsii.String(*auroraDbSecret),
		},
		Role:    apiLambdaRole,
		Timeout: awscdk.Duration_Seconds(jsii.Number(20)),
	})

	// Define the API Gateway
	spaApi := awsapigateway.NewRestApi(stack, jsii.String("SPAUndertownAPI"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("SPAUndertownAPI"),
		DeployOptions: &awsapigateway.StageOptions{
			TracingEnabled: jsii.Bool(true),
		},
		// CloudWatchRole: jsii.Bool(true),
	})

	// Define the SPA Authorizer
	spaAuth := awsapigateway.NewRequestAuthorizer(stack, jsii.String("GatewaySPAAuthorizer"), &awsapigateway.RequestAuthorizerProps{
		Handler: apiAuthorizerLambda,
		IdentitySources: &[]*string{
			awsapigateway.IdentitySource_Header(jsii.String(authorizationHeader)),
		},
	})
	// Define the Register Authorizer
	registerAuth := awsapigateway.NewRequestAuthorizer(stack, jsii.String("GatewayRegisterAuthorizer"), &awsapigateway.RequestAuthorizerProps{
		Handler: registerAuthorizerLambda,
		IdentitySources: &[]*string{
			awsapigateway.IdentitySource_Header(jsii.String(authorizationHeader)),
		},
	})

	// Create the /register route with the Lambda integration
	registerIntegration := awsapigateway.NewLambdaIntegration(registerLambda, &awsapigateway.LambdaIntegrationOptions{})
	registerResource := spaApi.Root().AddResource(jsii.String("register"), &awsapigateway.ResourceOptions{})
	registerResource.AddMethod(jsii.String("POST"), registerIntegration, &awsapigateway.MethodOptions{
		Authorizer: registerAuth,
	})
	// create the login route with lambda integration
	loginIntegration := awsapigateway.NewLambdaIntegration(loginLambda, &awsapigateway.LambdaIntegrationOptions{})
	loginResource := spaApi.Root().AddResource(jsii.String("login"), &awsapigateway.ResourceOptions{})
	loginResource.AddMethod(jsii.String("POST"), loginIntegration, &awsapigateway.MethodOptions{
		Authorizer: spaAuth,
	})

	deployment := awsapigateway.NewDeployment(stack, jsii.String("SPAAPIDeployment"), &awsapigateway.DeploymentProps{
		Api: spaApi,
	})

	// devLogGroup := awslogs.NewLogGroup(stack, jsii.String("devlogs"), &awslogs.LogGroupProps{})

	stage := awsapigateway.NewStage(stack, jsii.String("DevStage"), &awsapigateway.StageProps{
		Deployment: deployment,
		StageName:  jsii.String("dev"),
		// AccessLogDestination: awsapigateway.NewLogGroupLogDestination(devLogGroup),
		// AccessLogFormat: awsapigateway.AccessLogFormat_JsonWithStandardFields(&awsapigateway.JsonWithStandardFieldProps{
		// 	Caller:         jsii.Bool(false),
		// 	HttpMethod:     jsii.Bool(true),
		// 	Ip:             jsii.Bool(true),
		// 	Protocol:       jsii.Bool(true),
		// 	RequestTime:    jsii.Bool(true),
		// 	ResourcePath:   jsii.Bool(true),
		// 	ResponseLength: jsii.Bool(true),
		// 	Status:         jsii.Bool(true),
		// 	User:           jsii.Bool(true),
		// }),
	})
	spaApi.SetDeploymentStage(stage)
	// createPropertyLambda := awslambdago.NewGoFunction(stack, )

	// SSR
	// homeLambda := awslambdago.NewGoFunction(stack, jsii.String("ServerSideRender"), &awslambdago.GoFunctionProps{
	// 	Runtime:      awslambda.Runtime_PROVIDED_AL2(),
	// 	MemorySize:   jsii.Number(1024),
	// 	Architecture: awslambda.Architecture_ARM_64(),
	// 	Entry:        jsii.String("../services/ssr/homepage/lambda"),
	// 	Bundling:     bundlingOptions,
	// 	// Environment: &map[string]*string{
	// 	// 	"TABLE_NAME": db.TableName(),
	// 	// },
	// })
	// // Add a Function URL.
	// lambdaURL := homeLambda.AddFunctionUrl(&awslambda.FunctionUrlOptions{
	// 	AuthType: awslambda.FunctionUrlAuthType_NONE,
	// })
	// awscdk.NewCfnOutput(stack, jsii.String("homeLambdaURL"), &awscdk.CfnOutputProps{
	// 	ExportName: jsii.String("homeLambdaURL"),
	// 	Value:      lambdaURL.Url(),
	// })

	// assetsBucket := awss3.NewBucket(stack, jsii.String("assets"), &awss3.BucketProps{
	// 	BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
	// 	Encryption:        awss3.BucketEncryption_S3_MANAGED,
	// 	EnforceSSL:        jsii.Bool(true),
	// 	RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
	// 	Versioned:         jsii.Bool(false),
	// })

	// // Allow CloudFront to read from the bucket.
	// cfOAI := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("cfnOriginAccessIdentity"), &awscloudfront.OriginAccessIdentityProps{})
	// cfs := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{})
	// cfs.AddActions(jsii.String("s3:GetBucket*"))
	// cfs.AddActions(jsii.String("s3:GetObject*"))
	// cfs.AddActions(jsii.String("s3:List*"))
	// cfs.AddResources(assetsBucket.BucketArn())
	// cfs.AddResources(jsii.String(fmt.Sprintf("%v/*", *assetsBucket.BucketArn())))
	// cfs.AddCanonicalUserPrincipal(cfOAI.CloudFrontOriginAccessIdentityS3CanonicalUserId())
	// assetsBucket.AddToResourcePolicy(cfs)

	// // Add a CloudFront distribution to route between the public directory and the Lambda function URL.
	// lambdaURLDomain := awscdk.Fn_Select(jsii.Number(2), awscdk.Fn_Split(jsii.String("/"), lambdaURL.Url(), nil))
	// lambdaOrigin := awscloudfrontorigins.NewHttpOrigin(lambdaURLDomain, &awscloudfrontorigins.HttpOriginProps{
	// 	ProtocolPolicy: awscloudfront.OriginProtocolPolicy_HTTPS_ONLY,
	// })
	// cf := awscloudfront.NewDistribution(stack, jsii.String("customerFacing"), &awscloudfront.DistributionProps{
	// 	DefaultBehavior: &awscloudfront.BehaviorOptions{
	// 		AllowedMethods:       awscloudfront.AllowedMethods_ALLOW_ALL(),
	// 		Origin:               lambdaOrigin,
	// 		CachedMethods:        awscloudfront.CachedMethods_CACHE_GET_HEAD(),
	// 		OriginRequestPolicy:  awscloudfront.OriginRequestPolicy_ALL_VIEWER_EXCEPT_HOST_HEADER(),
	// 		CachePolicy:          awscloudfront.CachePolicy_CACHING_DISABLED(),
	// 		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
	// 	},
	// 	PriceClass: awscloudfront.PriceClass_PRICE_CLASS_100,
	// })

	// // Add /assets* to the distribution backed by S3.
	// assetsOrigin := awscloudfrontorigins.NewS3Origin(assetsBucket, &awscloudfrontorigins.S3OriginProps{
	// 	// Get content from the / directory in the bucket.
	// 	OriginPath:           jsii.String("/"),
	// 	OriginAccessIdentity: cfOAI,
	// })
	// cf.AddBehavior(jsii.String("/assets*"), assetsOrigin, nil)

	// // Export the domain.
	// awscdk.NewCfnOutput(stack, jsii.String("cloudFrontDomain"), &awscdk.CfnOutputProps{
	// 	ExportName: jsii.String("cloudfrontDomain"),
	// 	Value:      cf.DomainName(),
	// })

	// // Deploy the contents of the ./assets directory to the S3 bucket.
	// awss3deployment.NewBucketDeployment(stack, jsii.String("assetsDeployment"), &awss3deployment.BucketDeploymentProps{
	// 	DestinationBucket: assetsBucket,
	// 	Sources: &[]awss3deployment.ISource{
	// 		awss3deployment.Source_Asset(jsii.String("../services/ssr/homepage/assets"), nil),
	// 	},
	// 	DestinationKeyPrefix: jsii.String("assets"),
	// 	Distribution:         cf,
	// 	DistributionPaths:    jsii.Strings("/assets*"),
	// })

	return stack
}
