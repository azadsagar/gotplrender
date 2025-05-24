# Go Template Renderer with AWS Secrets Manager Integration

A command-line utility that renders Go templates using variables stored in AWS Secrets Manager. This tool is particularly useful for generating configuration files, documents, or any text-based output using templates with secure variable management.

## Features

- Uses Go's powerful templating engine
- Integrates with AWS Secrets Manager for secure variable storage
- Strict validation of template variables (fails on missing keys)
- Simple command-line interface

## Prerequisites

- Go 1.21 or later
- AWS credentials configured (any of the following):
  - AWS CLI configuration (`~/.aws/credentials`)
  - Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
  - IAM role when running on AWS services

## Installation

```bash
# Clone the repository
git clone https://github.com/azadsagar/gotplrender
cd gotplrender

# Install dependencies
go mod download

# Build the binary
go build -o gotplrender

# Cross-compile for Windows x64
GOOS=windows GOARCH=amd64 go build -o gotplrender.exe

# Cross-compile for Linux x64
GOOS=linux GOARCH=amd64 go build -o gotplrender

# Cross-compile for macOS x64
GOOS=darwin GOARCH=amd64 go build -o gotplrender_mac
```

## Usage

```bash
./gotplrender -templateFile=<template-file> -outputFile=<output-file> -varSource=secretsmanager -secretArn=<secret-arn> -region=<aws-region>
```

### Command Line Arguments

- `-templateFile` (required): Path to the Go template file
- `-outputFile` (required): Path where the rendered output should be written
- `-varSource` (required): Source of template variables (currently only supports "secretsmanager")
- `-secretArn` (required): ARN or name of the AWS Secrets Manager secret
- `-region` (required): AWS Region where the secret is stored (e.g., us-east-1, eu-west-1)

### Example

1. Create a template file `config.tpl`:
```
app:
  name: {{.appName}}
  environment: {{.environment}}
  database:
    host: {{.dbHost}}
    port: {{.dbPort}}
    username: {{.dbUser}}
```

2. Store variables in AWS Secrets Manager as JSON:
```json
{
  "appName": "MyApp",
  "environment": "production",
  "dbHost": "db.example.com",
  "dbPort": "5432",
  "dbUser": "admin"
}
```

3. Run the utility:
```bash
./gotplrender \
  -templateFile=config.tpl \
  -outputFile=config.yaml \
  -varSource=secretsmanager \
  -secretArn=arn:aws:secretsmanager:region:account:secret:name \
  -region=us-east-1
```

## Template Syntax

The utility uses Go's [text/template](https://golang.org/pkg/text/template/) package. Some common template syntax:

- `{{.variableName}}` - Insert variable value
- `{{if .condition}}...{{end}}` - Conditional rendering
- `{{range .items}}...{{end}}` - Iterate over arrays/slices
- `{{.nested.value}}` - Access nested JSON values

## Error Handling

The utility implements strict error handling:

- Missing template variables will cause an error (no silent failures)
- Template parsing errors are reported immediately
- AWS credentials and permissions issues are clearly reported
- File access errors are handled with clear messages

## AWS Permissions Required

The AWS role/user needs the following permissions:
- `secretsmanager:GetSecretValue` for the specific secret(s) being accessed


Example IAM policy:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "secretsmanager:GetSecretValue",
            "Resource": "arn:aws:secretsmanager:us-east-1:123456789012:secret:my-app-config-Ab1Cd2"
        }
    ]
}
```

Note: Replace the ARNs in the policy with your actual Secret and KMS key ARNs. Using specific ARNs instead of wildcards (*) is a security best practice that follows the principle of least privilege.

## Security Considerations

- AWS credentials should never be hardcoded in the application
- Use IAM roles when running in AWS environments
- Restrict IAM permissions to specific secrets when possible
- Template variables are kept in memory only during rendering

## Future Enhancements

- Support for AWS Systems Manager Parameter Store
- Multiple template file support
- Additional variable sources
- Variable validation and type checking
- Template function extensions

## TODO

### AWS SSM Parameter Store Support
The next major feature will be support for AWS Systems Manager Parameter Store. This will include:

- Parameter tree parsing (e.g., `/myapp/prod/db/*` to get all database-related parameters)
- Support for different parameter types (String, StringList, SecureString)
- Automatic parameter hierarchy to JSON structure conversion
- Recursive parameter tree traversal
- Support for both plain and encrypted parameters

Example planned usage:
```bash
./gotplrender \
  -templateFile=config.tpl \
  -outputFile=config.yaml \
  -varSource=ssm \
  -ssmPath=/myapp/prod \
  -recursive=true
```

The SSM integration will convert parameter paths into nested JSON structures:
```
/myapp/prod/db/host → {"db": {"host": "value"}}
/myapp/prod/db/port → {"db": {"port": "value"}}
/myapp/prod/cache/url → {"cache": {"url": "value"}}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

