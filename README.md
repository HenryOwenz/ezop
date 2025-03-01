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

[![Go CI](https://github.com/HenryOwenz/cloudgate/actions/workflows/go-ci.yml/badge.svg)](https://github.com/HenryOwenz/cloudgate/actions/workflows/go-ci.yml)

## Features

- **AWS Integration**
  - Multi-account/region management
  - CodePipeline operations (view, approve, trigger)
  - Pipeline status monitoring
  - Approval workflow management

- **Terminal UI**
  - Fast, keyboard-driven interface
  - Context-aware navigation
  - Visual feedback and safety controls

- **Coming Soon**
  - Azure integration
  - GCP support

## Installation

### Quick Install

```bash
# Linux/macOS
bash -c "$(curl -fsSL https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.sh)"

# Windows (PowerShell)
Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/HenryOwenz/cloudgate/main/scripts/install.ps1'))
```

### From Source

```bash
git clone https://github.com/HenryOwenz/cloudgate.git
cd cloudgate
make build
make install  # Installs as 'cg' in your $GOPATH/bin
```

## Requirements

- Go 1.21+
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
| Tab       | Toggle manual input      |
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

This design enables easy extension to support multiple cloud providers while maintaining a clean codebase.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 
