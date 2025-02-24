# Cloudgate

A seamless gateway to your cloud universe. Cloudgate is an elegant Terminal UI that reflects and unifies your multi-cloud operations across AWS, Azure, and GCP.

> *Like its namesake sculpture that reflects Chicago's skyline in liquid mercury, Cloudgate provides a fluid, unified view into your cloud landscape.*

## Features

- 🪞 **Unified Reflection** - View and manage multiple cloud providers through a single, elegant interface
  - AWS integration
  - Azure integration (Coming Soon)
  - GCP integration (Coming Soon)

- 🌐 **Multi-Account Management**
  - Seamless switching between accounts
  - Profile-based authentication (AWS)
  - Cross-account operations (coming soon)

- 🎨 **Beautiful Terminal UI**
  - Fluid navigation
  - Context-aware menus
  - Real-time status updates
  - Responsive layouts
  - Interactive selections
  - Loading indicators

- 🔄 **Cloud Operations**
  - Pipeline management
  - Approval workflows
  - Status monitoring
  - Resource operations
  - Real-time updates

- 🛡️ **Safety Features**
  - Clear operation preview
  - Multi-step confirmations
  - Easy cancellation
  - Error recovery
  - Context awareness

## Installation

```bash
# Clone the repository
git clone https://github.com/HenryOwenz/cloudgate.git
cd cloudgate

# Build the project
go build

# Run Cloudgate
./cloudgate
```

## Prerequisites

- Go 1.21 or later
- Cloud provider credentials configured:
  - AWS: `~/.aws/credentials` or `~/.aws/config`
  - Azure: Coming soon
  - GCP: Coming soon

## Required AWS Permissions

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

## Usage

1. Launch Cloudgate:
   ```bash
   ./cloudgate
   ```

2. Navigation:
   - ↑/↓: Navigate options
   - Enter: Select
   - Esc/-: Go back
   - Tab: Toggle input (where available)
   - q: Quit

3. Provider Setup:
   - Select cloud provider
   - Choose account/profile
   - Select region/location
   - Access your services

4. Operations:
   - Choose service category
   - Select specific operation
   - Follow interactive prompts
   - Monitor progress

## Key Bindings

- `↑/↓`: Navigate through options
- `Enter`: Select/Confirm
- `Esc/-`: Go back/Cancel
- `Tab`: Toggle manual input
- `q`: Quit application
- `Ctrl+c`: Force quit

## Future Enhancements

- Azure integration
- GCP support
- Enhanced pipeline visualization
- Cross-provider operations
- Resource management
- Cost optimization
- Security scanning
- Custom themes
- Configuration profiles
- Operation history
- Detailed analytics

## Project Structure

```
.
├── internal/
│   ├── aws/          # AWS provider operations
│   ├── azure/        # Azure operations (coming soon)
│   ├── gcp/          # GCP operations (coming soon)
│   └── ui/           # Terminal UI components
│       ├── constants/  # UI constants and enums
│       ├── model.go    # Main UI model and logic
│       └── styles.go   # UI styling definitions
├── main.go           # Application entry point
└── README.md         # This file
```

## Contributing

Contributions are welcome! Feel free to submit a Pull Request.

## The Name

Cloudgate is inspired by the famous Cloud Gate sculpture in Chicago, known for its liquid mercury surface that reflects and transforms the city's skyline. Similarly, our tool provides a reflective interface that unifies and transforms how you interact with multiple cloud providers.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
