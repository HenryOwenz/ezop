# ciselect

A CLI tool for managing AWS CodePipeline manual approvals efficiently.

## Installation

```bash
go install github.com/HenryOwenz/ciselect@latest
```

## Prerequisites

- AWS credentials configured in your AWS config/credentials files
- Appropriate IAM permissions for AWS CodePipeline operations

## Usage

All commands require both an AWS profile and region to be specified:
- `--profile` or `-p`: AWS profile to use
- `--region` or `-r`: AWS region to use

This ensures you're always aware of which AWS account and region you're operating in.

### List pending approvals
```bash
ciselect list --profile my-aws-profile --region us-west-2
```

### Approve an action
```bash
ciselect approve <pipeline-name> <stage-name> <action-name> --profile my-aws-profile --region us-west-2 -s "Approved by ciselect"
```

### Reject an action
```bash
ciselect reject <pipeline-name> <stage-name> <action-name> --profile my-aws-profile --region us-west-2 -s "Rejected by ciselect"
```

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

- **Required AWS Profile**: The tool requires explicit specification of the AWS profile to use, preventing accidental operations on the wrong AWS account.
- **Required AWS Region**: The tool requires explicit specification of the AWS region, ensuring clarity about where operations are being performed.
- **Explicit Permissions**: The IAM policy above specifies the exact permissions needed for operation.
- **Error Handling**: Clear error messages when profile, region, or other required parameters are missing or invalid. 