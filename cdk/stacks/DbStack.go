package stacks

import (
	"encoding/json"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

var (
	DB_STACK_KEY_DB_HOST   = "dbHost"
	DB_STACK_KEY_DB_SECRET = "dbSecret"
	DB_STACK_KEY_DB_NAME   = "dbName"
	DB_STACK_VALUE_DB_NAME = "undertown_v3"
)

type DbStackProps struct {
	awscdk.StackProps
}

// TODO it looks like there are too many problems deploying this aurora stack
func DbStack(scope constructs.Construct, id string, props *DbStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	stack := awscdk.NewStack(scope, &id, &sprops)
	var dbUsername = "postgres"
	secretContent := map[string]string{
		"username": dbUsername,
	}
	secretString, _ := json.Marshal(secretContent)

	// crate the vpc for the postgres db
	vpc := awsec2.NewVpc(stack, jsii.String("VPC"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})
	// Create a database secret
	secret := awssecretsmanager.NewSecret(stack, jsii.String("DBSecret"), &awssecretsmanager.SecretProps{
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			SecretStringTemplate: jsii.String(string(secretString)),
			GenerateStringKey:    jsii.String("password"),
			ExcludeCharacters:    jsii.String("/@\""),
		},
	})

	// create the aurora postgres database
	// store the database password and username in the secret
	//
	auroraCluster := awsrds.NewServerlessCluster(stack, jsii.String("AuroraCluster"), &awsrds.ServerlessClusterProps{
		Engine: awsrds.DatabaseClusterEngine_AuroraPostgres(&awsrds.AuroraPostgresClusterEngineProps{
			Version: awsrds.AuroraPostgresEngineVersion_VER_13_4(),
		}),
		Vpc:                 vpc,
		Credentials:         awsrds.Credentials_FromSecret(secret, jsii.String(dbUsername)),
		DefaultDatabaseName: jsii.String(DB_STACK_VALUE_DB_NAME),
		RemovalPolicy:       awscdk.RemovalPolicy_DESTROY,
		EnableDataApi:       jsii.Bool(true),
		Scaling:             &awsrds.ServerlessScalingOptions{MaxCapacity: awsrds.AuroraCapacityUnit_ACU_8},
	},
	)
	// TODO the secret is used for username and password
	// stack.ExportValue(secret.SecretArn(), &awscdk.ExportValueOptions{
	// 	Name: jsii.String(STACK_OUTPUT_DB_SECRET_KEY),
	// })
	awscdk.NewCfnOutput(stack, jsii.String(DB_STACK_KEY_DB_SECRET), &awscdk.CfnOutputProps{
		Value:      secret.SecretArn(),
		ExportName: jsii.String(DB_STACK_KEY_DB_SECRET),
	})

	awscdk.NewCfnOutput(stack, jsii.String(DB_STACK_KEY_DB_HOST), &awscdk.CfnOutputProps{
		Value:      auroraCluster.ClusterEndpoint().Hostname(),
		ExportName: jsii.String(DB_STACK_KEY_DB_HOST),
	})

	awscdk.NewCfnOutput(stack, jsii.String(DB_STACK_KEY_DB_NAME), &awscdk.CfnOutputProps{
		Value:      jsii.String(DB_STACK_VALUE_DB_NAME),
		ExportName: jsii.String(DB_STACK_KEY_DB_NAME),
	})

	return stack
}
