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
	Vpc awsec2.Vpc
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

	// Create a database secret
	secret := awssecretsmanager.NewSecret(stack, jsii.String("DBSecret"), &awssecretsmanager.SecretProps{
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			SecretStringTemplate: jsii.String(string(secretString)),
			GenerateStringKey:    jsii.String("password"),
			ExcludeCharacters:    jsii.String("/@\""),
		},
	})

	databaseSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("DatabaseSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc: props.Vpc,
	})

	databaseSecurityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(5432)), jsii.String("Allow Access To PG"), nil)
	// create the aurora postgres database
	// store the database password and username in the secret
	// TODO the database is provisioned in public subnets without NAT gateway
	// NAT gateway is expensive (and it's used mainly to run db in private subnets with connection through nat to the internet)
	auroraCluster := awsrds.NewDatabaseCluster(stack, jsii.String("AuroraServerlessCluster"), &awsrds.DatabaseClusterProps{
		Engine: awsrds.DatabaseClusterEngine_AuroraPostgres(&awsrds.AuroraPostgresClusterEngineProps{
			Version: awsrds.AuroraPostgresEngineVersion_VER_14_4(),
		}),
		Readers: &[]awsrds.IClusterInstance{
			awsrds.ClusterInstance_ServerlessV2(jsii.String("reader-instance"), &awsrds.ServerlessV2ClusterInstanceProps{
				PubliclyAccessible:      jsii.Bool(true),
				AutoMinorVersionUpgrade: jsii.Bool(true),
			}),
		},
		Writer: awsrds.ClusterInstance_ServerlessV2(jsii.String("writer-instance"), &awsrds.ServerlessV2ClusterInstanceProps{
			PubliclyAccessible: jsii.Bool(true),
		}),
		Vpc:                     props.Vpc,
		Credentials:             awsrds.Credentials_FromSecret(secret, jsii.String(dbUsername)),
		DefaultDatabaseName:     jsii.String(DB_STACK_VALUE_DB_NAME),
		ServerlessV2MaxCapacity: jsii.Number(4),
		SecurityGroups:          &[]awsec2.ISecurityGroup{databaseSecurityGroup},
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
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
