```bash             __                               __                           __               
            /  |                             /  |                         /  |              
   _______  $$ |   ______    __    __    ____$$ |   ______     ______    _$$ |_      ______  
  /       | $$ |  /      \  /  |  /  |  /    $$ |  /      \   /      \  / $$   |    /      \ 
 /$$$$$$$/  $$ | /$$$$$$  | $$ |  $$ | /$$$$$$$ | /$$$$$$  |  $$$$$$  | $$$$$$/    /$$$$$$  |
 $$ |       $$ | $$ |  $$ | $$ |  $$ | $$ |  $$ | $$ |  $$ |  /    $$ |   $$ | __  $$    $$ |
 $$ \_____  $$ | $$ \__$$ | $$ \__$$ | $$ \__$$ | $$ \__$$ | /$$$$$$$ |   $$ |/  | $$$$$$$$/ 
 $$       | $$ | $$    $$/  $$    $$/  $$    $$ | $$    $$ | $$    $$ |   $$  $$/  $$       |
   $$$$$$$/ $$/   $$$$$$/    $$$$$$/    $$$$$$$/   $$$$$$$ |  $$$$$$/     $$$$/     $$$$$$$/ 
                                                  /  \__$$ |                              
                                                  $$    $$/                               
                                                   $$$$$$/                                
```

# cloudgate
A seamless gateway to your cloud universe. cloudgate is a Terminal Application that unifies your multi-cloud operations across AWS, Azure, and GCP.
> *Where your clouds converge.*

## Features

### Cloud Provider Integration
- **AWS** - Full support for AWS services and operations
- **Azure** - Coming soon
- **GCP** - Coming soon

### Core Capabilities

#### AWS
- **Multi-Account Management** - Switch between AWS profiles and regions seamlessly
- **CodePipeline Integration** - View, approve, and trigger pipeline executions
- **Pipeline Status Monitoring** - Real-time updates on pipeline states and stages
- **Approval Workflows** - Manage manual approval actions in CodePipeline

#### Azure (Coming Soon)

#### GCP (Coming Soon)

### User Experience
- **Terminal-Based UI** - Fast, responsive, and keyboard-driven interface
- **Context-Aware Navigation** - Intuitive menus that adapt to your workflow
- **Safety Controls** - Confirmation steps and clear operation previews
- **Visual Feedback** - Loading indicators and status updates

## Installation

### Linux and MacOS

```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.sh)"
```

### Windows

Open PowerShell and run:

```powershell
Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.ps1'))
```

After installation, you can run cloudgate using the `cg` command from anywhere in your terminal.

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
                "<service>:<action>",
                
            ],
            "Resource": "arn:aws:<service>:*:*:*"
        }
    ]
}
```

## Usage

1. Launch cloudgate:
   ```bash
   cg
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

## Architecture

cloudgate is built with a dual-layer architecture that separates cloud provider implementations from the application's business logic. This design enables easy extension to support multiple cloud providers while maintaining a clean and maintainable codebase.

For detailed information about the architecture and design patterns used in cloudgate, please refer to the [documentation directory](documentation/README.md).

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

## Contributing

Contributions are welcome! Feel free to submit a Pull Request.

Please read our [documentation](documentation/README.md) to understand the architecture and design patterns used in cloudgate before contributing.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
