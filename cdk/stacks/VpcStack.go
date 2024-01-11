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

func VpcStack(scope constructs.Construct, id string, props *VpcStackProps) awsec2.Vpc {
	var sprops awscdk.StackProps
	stack := awscdk.NewStack(scope, &id, &sprops)
	// crate the vpc for the postgres db
	vpc := awsec2.NewVpc(stack, jsii.String("VPC"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
		// For a VPC with only ISOLATED subnets, this value will be undefined.
		CreateInternetGateway: jsii.Bool(true),
		// Creating the VPC with only isolated subnets will not create a NAT gateway
		// NAT gateways are expensive
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("public-subnet"),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
		},
		NatGateways: jsii.Number(0),
	})

	return vpc
}
