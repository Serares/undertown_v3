module cdk

go 1.21.4

toolchain go1.21.6

require (
	github.com/aws/aws-cdk-go/awscdk/v2 v2.120.0
	github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2 v2.117.0-alpha.0
	github.com/aws/constructs-go/constructs/v10 v10.3.0
	github.com/aws/jsii-runtime-go v1.94.0
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/Serares/undertown_v3/utils v0.0.0
	github.com/cdklabs/awscdk-asset-awscli-go/awscliv1/v2 v2.2.201 // indirect
	github.com/cdklabs/awscdk-asset-kubectl-go/kubectlv20/v2 v2.1.2 // indirect
	github.com/cdklabs/awscdk-asset-node-proxy-agent-go/nodeproxyagentv6/v2 v2.0.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/yuin/goldmark v1.4.13 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
)

replace github.com/Serares/undertown_v3/utils => ../utils
