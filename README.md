# ezop

A user-friendly interactive CLI tool for managing cloud operations across multiple providers.

## Installation

```bash
go install github.com/HenryOwenz/ezop@latest
```

## Prerequisites

- Cloud provider credentials configured appropriately
- For AWS: AWS credentials configured in your AWS config/credentials files
- For Azure: Coming soon
- For GCP: Coming soon

## Usage

Simply run:
```bash
ezop
```

The interactive interface will guide you through:
1. Selecting your cloud provider
2. Choosing the service and operation
3. Configuring provider-specific settings
4. Managing your cloud operations with a beautiful terminal UI

## Features

The interactive interface provides:
- Beautiful terminal UI with color-coded elements
- Multi-cloud provider support
  - AWS (Available)
  - Azure (Coming Soon)
  - GCP (Coming Soon)
- Provider-specific features:
  - AWS:
    - CodePipeline manual approval management
    - More services coming soon
  - Azure: Coming soon
  - GCP: Coming soon
- Interactive selection and management
- Guided workflow with clear steps
- Confirmation steps for safety

## Required Permissions

### AWS
For AWS CodePipeline operations, the following IAM permissions are required:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "codepipeline:ListActionExecutions",
                "codepipeline:PutApprovalResult",
                "codepipeline:ListPipelines",
                "codepipeline:GetPipelineState"
            ],
            "Resource": "arn:aws:codepipeline:*:*:*"
        }
    ]
}
```

### Azure
Coming soon

### GCP
Coming soon

## Development

To build from source:

```bash
git clone https://github.com/HenryOwenz/ezop.git
cd ezop
go build
```

## Safety Features

- **Multi-step Workflow**: Clear step-by-step process prevents accidental operations
- **Provider Selection**: Choose from available cloud providers
- **Service Selection**: Select from available services for each provider
- **Operation Selection**: Choose specific operations within each service
- **Confirmation Steps**: Verify your actions before they're executed
- **Clear Context**: Always shows which provider, service, and operation you're working with
- **Color-Coded UI**: Important information and actions are visually distinct
- **Error Handling**: Clear error messages when something goes wrong 