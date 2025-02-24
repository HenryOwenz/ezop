# ezop

A user-friendly interactive CLI tool for managing cloud operations across multiple providers, with a focus on AWS CodePipeline operations.

## Features

- 🎨 Beautiful terminal UI with:
  - Color-coded elements for better visibility
  - Dynamic context-aware navigation
  - Responsive table layouts
  - Interactive selection menus
  - Loading spinners for async operations
  - Clear error handling and display

- 🔄 AWS CodePipeline Operations:
  - View pipeline status and stages
  - Manage manual approval actions
  - Start pipeline executions
  - View detailed stage information
  - Real-time status updates

- 🛠️ AWS Configuration:
  - Automatic AWS profile detection
  - Region selection
  - Profile-based authentication
  - Support for multiple AWS profiles

- 🎯 Operation Categories:
  - Workflows
    - Pipeline Approvals
    - Pipeline Status
    - Start Pipeline Execution
  - Operations (Coming Soon)

- 🔒 Safety Features:
  - Multi-step confirmation process
  - Clear context display
  - Operation preview
  - Cancel options at every step
  - Error state recovery

## Installation

```bash
# Clone the repository
git clone https://github.com/HenryOwenz/ezop.git
cd ezop

# Build the project
go build

# Run the application
./ezop
```

## Prerequisites

- Go 1.21 or later
- AWS credentials configured in `~/.aws/credentials` or `~/.aws/config`
- Required AWS IAM permissions (see below)

## Required AWS Permissions

The following IAM permissions are required for AWS CodePipeline operations:

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
                "codepipeline:GetPipelineState",
                "codepipeline:StartPipelineExecution"
            ],
            "Resource": "arn:aws:codepipeline:*:*:*"
        }
    ]
}
```

## Usage Guide

1. Launch the application:
   ```bash
   ./ezop
   ```

2. Navigation:
   - Use ↑/↓ arrows to navigate
   - Press Enter to select
   - Press Esc or - to go back
   - Press q to quit
   - Press Tab to toggle manual input (where available)

3. AWS Configuration:
   - Select AWS profile from the list or enter manually
   - Choose AWS region from the list or enter manually

4. Operations:
   - Select AWS service (currently CodePipeline)
   - Choose operation category
   - Select specific operation
   - Follow the interactive prompts

## Key Bindings

- `↑/↓`: Navigate through options
- `Enter`: Select/Confirm
- `Esc/-`: Go back/Cancel
- `Tab`: Toggle manual input (where available)
- `q`: Quit application
- `Ctrl+c`: Force quit

## Future Enhancements

- Additional AWS Services support
- Azure integration
- GCP integration
- More CodePipeline operations
- Enhanced pipeline visualization
- Custom theme support
- Configuration file support
- Pipeline execution history
- Detailed stage information
- Cross-region operation support

## Development

The project structure follows Go best practices:

```
.
├── internal/
│   ├── aws/          # AWS-specific functionality
│   └── ui/           # Terminal UI components
│       ├── constants/  # UI constants and enums
│       ├── model.go    # Main UI model and logic
│       └── styles.go   # UI styling definitions
├── main.go           # Application entry point
└── README.md         # This file
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 
