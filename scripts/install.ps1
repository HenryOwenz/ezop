# Installation script for Windows

$ErrorActionPreference = 'Stop'

# Colors for output
$ORANGE = [char]0x1b + "[38;5;208m"
$NC = [char]0x1b + "[0m"

# Display banner
Write-Host "$ORANGE
             __                               __                           __               
            /  |                             /  |                         /  |              
   _______  `$`$ |   ______    __    __    ____`$`$ |   ______     ______    _`$`$ |_      ______  
  /       | `$`$ |  /      \  /  |  /  |  /    `$`$ |  /      \   /      \  / `$`$   |    /      \ 
 /`$`$`$`$`$`$`$/  `$`$ | /`$`$`$`$`$`$`$  | `$`$ |  `$`$ | /`$`$`$`$`$`$`$ | /`$`$`$`$`$`$  |  `$`$`$`$`$`$  | `$`$`$`$`$/    /`$`$`$`$`$`$  |
 `$`$ |       `$`$ | `$`$ |  `$`$ | `$`$ |  `$`$ | `$`$ |  `$`$ | `$`$ |  `$`$ |  /    `$`$ |   `$`$ | __  `$`$    `$`$ |
 `$`$ \_____  `$`$ | `$`$ \__`$`$ | `$`$ \__`$`$ | `$`$ \__`$`$ | `$`$ \__`$`$ | /`$`$`$`$`$`$`$ |   `$`$ |/  | `$`$`$`$`$`$`$`$/ 
 `$`$       | `$`$ | `$`$    `$`$/  `$`$    `$`$/  `$`$    `$`$ | `$`$    `$`$ | `$`$    `$`$ |   `$`$  `$`$/  `$`$       |
  `$`$`$`$`$`$`$/  `$`$/   `$`$`$`$`$`$/    `$`$`$`$`$`$/    `$`$`$`$`$`$`$/   `$`$`$`$`$`$`$ |  `$`$`$`$`$`$/     `$`$`$`$/    `$`$`$`$`$`$`$/ 
                                                  /  \__`$`$ |                              
                                                  `$`$    `$`$/                               
                                                   `$`$`$`$`$`$/                                
$NC"

Write-Host "Installing cloudgate..." -ForegroundColor Blue

# Create installation directory
$InstallDir = "$env:LOCALAPPDATA\cloudgate"
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Detect architecture
$arch = (Get-WmiObject Win32_OperatingSystem).OSArchitecture
$arch = if ($arch -eq "64-bit") { "amd64" } else { "arm64" }
$binary = "cloudgate_windows_${arch}.exe"
$LatestRelease = "https://github.com/HenryOwenz/cloudgate/releases/latest/download/$binary"
$OutFile = "$InstallDir\cg.exe"

Write-Host "Downloading latest release..." -ForegroundColor Blue
Invoke-WebRequest -Uri $LatestRelease -OutFile $OutFile

# Add to PATH
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable(
        "Path",
        "$UserPath;$InstallDir",
        "User"
    )
}

Write-Host "Successfully installed cloudgate!" -ForegroundColor Green
Write-Host "Please restart your terminal to use the 'cg' command." 