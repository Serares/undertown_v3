package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type VpcStackProps struct {
	awscdk.StackProps
}

func VPC(scope constructs.Construct, id string, props *VpcStackProps) awsec2.Vpc {
	var sprops awscdk.StackProps
	stack := awscdk.NewStack(scope, &id, &sprops)
	vpc := awsec2.NewVpc(stack, jsii.String("VPC"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
		// For a VPC with only ISOLATED subnets, this value will be undefined.
		CreateInternetGateway: jsii.Bool(true),
		// Creating the VPC with only isolated subnets will not create a NAT gateway
		// NAT gateways are expensive
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("private-subnet"),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
		},
		NatGateways: jsii.Number(0),
	})

	// create the vpc endpoint for secretsmanager
	vpc.AddInterfaceEndpoint(jsii.String("SecretsManagerEndpoint"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_SECRETS_MANAGER(),
		Subnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})

	// create the vpc endpoint for S3
	// if lambda is not able to send to S3
	// you will have to use sqs messages and a proxy lambda to send the images to s3
	// TODO not sure if this works yet
	vpc.AddInterfaceEndpoint(jsii.String("S3Endpoint"), &awsec2.InterfaceVpcEndpointOptions{
		Service: awsec2.InterfaceVpcEndpointAwsService_S3(),
		Subnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
		},
	})
	return vpc
}
