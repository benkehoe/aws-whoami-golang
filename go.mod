module github.com/benkehoe/aws-whoami-golang

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.17.5
	github.com/aws/aws-sdk-go-v2/config v1.18.14
	github.com/aws/aws-sdk-go-v2/service/iam v1.19.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.4
	github.com/aws/smithy-go v1.13.5
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.13.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.29 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.30 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.3 // indirect
)
