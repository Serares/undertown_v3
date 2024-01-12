package stacks

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type ApiLambdaStackProps struct {
	awscdk.StackProps
	Vpc awsec2.Vpc
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
	// // Security group for lambdas
	lambdaSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("LambdaSG"), &awsec2.SecurityGroupProps{
		Vpc: props.Vpc,
		// ... Lambda SG configurations ...
	})
	lambdaSecurityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_AllTcp(), jsii.String("APILambdaSecurityGroup"), nil)
	// Lambda role to access secrets manager
	apiLambdaRole := awsiam.NewRole(stack, jsii.String("APIServicesLambdaRole"), &awsiam.RoleProps{
		// TODO think about if this needs to apply the least priviledge principle
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})
	apiLambdaRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaBasicExecutionRole")))
	apiLambdaRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("service-role/AWSLambdaVPCAccessExecutionRole")))
	apiLambdaRole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonRDSDataFullAccess")))
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
		Timeout: awscdk.Duration_Seconds(jsii.Number(60 * 2)),
		Vpc:     props.Vpc,
		// TODO read this
		// if lambdas are part of a VPC they will need EIP's assigned to the subnets that lambdas are part of
		// https://stackoverflow.com/questions/52992085/why-cant-an-aws-lambda-function-inside-a-public-subnet-in-a-vpc-connect-to-the/52994841#52994841
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
		SecurityGroups: &[]awsec2.ISecurityGroup{
			lambdaSecurityGroup,
		},
	})

	// proxy lambda
	// used to delegate messages to sqs

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

	return stack
}
