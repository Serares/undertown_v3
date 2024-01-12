**Each module has a counterpart `main.go` file at the root of the directory to simplify local testing**

- goose for db migrations
- SQLC for db transactions

- Infra revolves around using API GW with lambda integrations to SSR html templates and at the same time use the lambdas to proxy to lambda backends

**Try using ports starting with `40` for API and ports starting with `30` for SSR**

**To import local modules**

Utils

- go mod edit -require=github.com/Serares/undertown_v3/utils@v0.0.0
- go mod edit -replace=github.com/Serares/undertown_v3/utils=../../../utils

Repository

- go mod edit -require=github.com/Serares/undertown_v3/repositories/repository@v0.0.0
- go mod edit -replace=github.com/Serares/undertown_v3/repositories/repository=../../../repositories/repository

**Decided to use the cdk instead of cloudformation or aws sam**

- It was harder to locally run go lambdas using the sam cli (the only option that worked was to run the go lambdas in docker images and this increased the development complexity)
- The configuration yaml files proved to be a lot harder to create (lots of configuration had to be set in place) and the cdk also uploads the static files to a S3 bucket (when using the AWS SAM the uploading had to be done separatly by a .sh script)
- The cdk has a golang library which is a + because you can use one language throughout all the workenv

**Serverless V2 engine**

golang`
auroraCluster := awsrds.NewDatabaseCluster(stack, jsii.String("AuroraServerlessCluster"), &awsrds.DatabaseClusterProps{
Engine: awsrds.DatabaseClusterEngine_AuroraPostgres(&awsrds.AuroraPostgresClusterEngineProps{
Version: awsrds.AuroraPostgresEngineVersion_VER_14_4(),
}),
Readers: &[]awsrds.IClusterInstance{
awsrds.ClusterInstance_ServerlessV2(jsii.String("reader-instance"), &awsrds.ServerlessV2ClusterInstanceProps{
PubliclyAccessible: jsii.Bool(true),
AutoMinorVersionUpgrade: jsii.Bool(true),
}),
},
Writer: awsrds.ClusterInstance_ServerlessV2(jsii.String("writer-instance"), &awsrds.ServerlessV2ClusterInstanceProps{
PubliclyAccessible: jsii.Bool(true),
AutoMinorVersionUpgrade: jsii.Bool(true),
}),
Vpc: props.Vpc,
Credentials: awsrds.Credentials_FromSecret(secret, jsii.String(dbUsername)),
DefaultDatabaseName: jsii.String(DB_STACK_VALUE_DB_NAME),
ServerlessV2MaxCapacity: jsii.Number(4),
SecurityGroups: &[]awsec2.ISecurityGroup{databaseSecurityGroup},
VpcSubnets: &awsec2.SubnetSelection{
SubnetType: awsec2.SubnetType_PUBLIC,
},
},
)`

**This is a working configuration but it's not helping my case**

- Serverless v2 always has a capacity of `0.5` meaning that it's not trully serverless and it will charge even when it's not used

- Have to use serverless v1 and do some hacks to make it work with lambdas
