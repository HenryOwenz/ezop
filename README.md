# ciselect

A user-friendly interactive CLI tool for managing AWS CodePipeline manual approvals.

## Installation

```bash
go install github.com/HenryOwenz/ciselect@latest
```

## Prerequisites

- AWS credentials configured in your AWS config/credentials files
- Appropriate IAM permissions for AWS CodePipeline operations

## Usage

Simply run:
```bash
ciselect
```

The interactive interface will guide you through:
1. Selecting an AWS profile (from your configured profiles or type a custom one)
2. Choosing an AWS region (from common regions or type a custom one)
3. Managing your pipeline approvals with a beautiful terminal UI

## Features

The interactive interface provides:
- Beautiful terminal UI with color-coded elements
- List of available AWS profiles with ability to type custom ones
- List of AWS regions with ability to type custom ones
- Clear, formatted table of pending approvals
- Interactive selection and management of approvals
- Guided approve/reject workflow
- Confirmation steps for safety

## Required AWS Permissions

The following IAM permissions are required:

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

## Development

To build from source:

```bash
git clone https://github.com/HenryOwenz/ciselect.git
cd ciselect
go build
```

## Safety Features

- **Interactive Workflow**: Clear step-by-step process prevents accidental operations
- **Profile Selection**: Choose from available AWS profiles or type custom ones
- **Region Selection**: Choose from common regions or type custom ones
- **Confirmation Steps**: Verify your actions before they're executed
- **Clear Context**: Always shows which profile, region, and pipeline you're working with
- **Color-Coded UI**: Important information and actions are visually distinct
- **Error Handling**: Clear error messages when something goes wrong 