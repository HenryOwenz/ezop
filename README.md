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

A terminal-based application that unifies multi-cloud operations across AWS, Azure, and GCP.

> *Where your clouds converge.*

[![Lint](https://github.com/HenryOwenz/cloudgate/actions/workflows/lint.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/lint.yml)
[![Build](https://github.com/HenryOwenz/cloudgate/actions/workflows/build.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/build.yml)
[![Test](https://github.com/HenryOwenz/cloudgate/actions/workflows/test.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/test.yml)

## Features

- **AWS Integration**
  - Multi-account/region management


  <details>
  <summary><b>📋 Available AWS Services & Operations</b></summary>
  
  | Service | Operation | Description |
  |---------|-----------|-------------|
  | **Lambda** | | |
  | | Function Status | View all Lambda functions with runtime and last update info<br><br>**Function Details View:**<br>Select any function to inspect detailed configuration including:<br>• Memory allocation<br>• Timeout settings<br>• Code size<br>• Package type<br>• Architecture<br>• Role ARN<br>• Log group |
  | **CodePipeline** | | |
  | | Pipeline Status | View status of all pipelines and their stages |
  | | Pipeline Approvals | List, approve, or reject pending manual approvals |
  | | Start Pipeline | Trigger pipeline execution with latest commit or specific revision |
  
  *Operations can be performed using any configured AWS profile and region (one active profile/region at a time)*  
  *Multi-account aggregation for services will be coming in the future*
  </details>

- **Terminal UI**
  - Fast, keyboard-driven interface
  - Context-aware navigation
  - Visual feedback and safety controls
  - Formatted display of timestamps and resource sizes

- **Coming Soon**
  - Azure integration
  - GCP support
  - Additional AWS services (S3, EC2, etc.)

## Installation

### Quick Install / Upgrade

**Linux/macOS:**
```bash
bash -c "$(curl -fsSL https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.sh)"
```

**Windows (PowerShell):**
```powershell
Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.ps1'))
```

These scripts will automatically download and install the latest version of cloudgate, handling upgrades cleanly if you already have it installed.

### From Source

```bash
git clone https://github.com/HenryOwenz/cloudgate.git
cd cloudgate
make build
make install  # Installs as 'cg' in your $GOPATH/bin
```

## Requirements

- Go 1.22+
- AWS credentials configured in `~/.aws/credentials` or `~/.aws/config`

## Usage

```bash
cg  # Launch the application
```

### Navigation

| Key       | Action                   |
|-----------|--------------------------|
| ↑/↓       | Navigate options         |
| Enter     | Select/Confirm           |
| Esc/-     | Go back/Cancel           |
| q         | Quit application         |
| Ctrl+c    | Force quit               |

## Development

### Testing

```bash
make test          # Run all tests
make test-unit     # Run unit tests only
make test-integration  # Run integration tests only
make test-coverage  # Generate coverage report
```

### CI/CD

This project uses GitHub Actions for continuous integration:
- Automated builds on each push and pull request
- Unit and integration tests
- Code linting with golangci-lint
- Test coverage reporting

## Architecture

cloudgate uses a dual-layer architecture:
- Provider layer: Abstracts cloud provider APIs
- UI layer: Handles user interaction and workflow

The application follows a modular design pattern that makes it easy to add new cloud services and operations. Each service is implemented as a separate module with clear interfaces, allowing for independent development and testing.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
