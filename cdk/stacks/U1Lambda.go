package stacks

import (
	"net/http"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type U1LambdaProps struct {
	awscdk.StackProps
	Vpc         awsec2.Vpc
	authorizers APIAuthorizersMap
}

type DbCfg struct {
	PSQL_USER     string
	PSQL_PASSWORD string
	PSQL_DB       string
	PSQL_HOST     string
	PSQL_PORT     int32
	PSQL_ENDPOINT string
}

// This stack is used for CRUD operations
func U1Lambda(scope constructs.Construct, id string, props *U1LambdaProps) []IntegrationLambda {
	var sprops awscdk.StackProps
	var lambdas []IntegrationLambda

	// TODO create this secret in SSM
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	auroraHost := awscdk.Fn_ImportValue(jsii.String(DB_STACK_KEY_DB_HOST))
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
		Bundling:     BundlingOptions,
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

	// Login Lambda
	loginLambda := awslambdago.NewGoFunction(stack, jsii.String("SPALogin"), &awslambdago.GoFunctionProps{
		Runtime:      awslambda.Runtime_PROVIDED_AL2(),
		MemorySize:   jsii.Number(1024),
		Architecture: awslambda.Architecture_ARM_64(),
		Entry:        jsii.String("../services/api/login/lambda"),
		Bundling:     BundlingOptions,
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

	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &registerLambda,
		path:       "/register",
		method:     http.MethodPost,
		authorizer: props.authorizers[RegisterAuthorizerName].handler,
	})
	lambdas = append(lambdas, IntegrationLambda{
		goLambda:   &loginLambda,
		path:       "/login",
		method:     http.MethodPost,
		authorizer: nil,
	})

	return lambdas
}
