# build.ps1
param (
    [string]$ProjectName = "fsearch",
    [string]$CliPath = "./cli"
)

# Directories
$OutputDir = "dist"
$BuildDir = "build"

# Get version from git, fallback if not available
try {
    $Version = git describe --tags --always
} catch {
    $Version = "v0.0.1"
}

# Platforms to build for
$Targets = @(
    @{ OS="windows"; ARCH="amd64" },
    @{ OS="windows"; ARCH="386" },
    @{ OS="linux"; ARCH="amd64" },
    @{ OS="linux"; ARCH="386" },
    @{ OS="darwin"; ARCH="amd64" },
    @{ OS="darwin"; ARCH="arm64" }
)

# Prepare directories
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

foreach ($target in $Targets) {
    $OS = $target.OS
    $ARCH = $target.ARCH

    Write-Host "Building $ProjectName for $OS/$ARCH..."
    
    # Name the binary just "fsearch" (add .exe on Windows)
    $ext = if ($OS -eq "windows") { ".exe" } else { "" }
    $BinaryName = "$ProjectName$ext"
    $BinaryPath = "$BuildDir\$BinaryName"

    # Set environment variables for cross-compilation
    $env:GOOS = $OS
    $env:GOARCH = $ARCH

    go build -o $BinaryPath $CliPath

    if (-Not (Test-Path $BinaryPath)) {
        Write-Error "Build failed for $OS/$ARCH"
        exit 1
    }

    # Package binary into ZIP
    $ArchiveName = "${ProjectName}_${Version}_${OS}_${ARCH}.zip"
    $ArchivePath = "$OutputDir\$ArchiveName"
    Compress-Archive -Path $BinaryPath -DestinationPath $ArchivePath -Force

    # Generate SHA256
    $FileHash = (Get-FileHash $ArchivePath -Algorithm SHA256).Hash
    $HashFile = "$ArchivePath.sha256"
    $FileHash | Out-File $HashFile

    Write-Host "Built and packaged $ArchiveName with SHA256: $FileHash"
}

Write-Host "All builds complete! Packages are in $OutputDir"
