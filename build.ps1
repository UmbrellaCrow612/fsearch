# build.ps1
param (
    [string]$ProjectName = "mycli",
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

# OS/Arch
$OS = "windows"
$ARCH = if ([System.Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Prepare directories
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

# Build the Go binary
Write-Host "Building $ProjectName..."
$BinaryPath = "$BuildDir\$ProjectName.exe"
go build -o $BinaryPath $CliPath

if (-Not (Test-Path $BinaryPath)) {
    Write-Error "Build failed, binary not found."
    exit 1
}

# Package binary into ZIP
$ArchiveName = "${ProjectName}_${Version}_${OS}_${ARCH}.zip"
$ArchivePath = "$OutputDir\$ArchiveName"
Write-Host "Packaging $ArchiveName..."
Compress-Archive -Path $BinaryPath -DestinationPath $ArchivePath -Force

# Generate SHA256 dynamically
Write-Host "Generating SHA256..."
$FileHash = (Get-FileHash $ArchivePath -Algorithm SHA256).Hash
$HashFile = "$ArchivePath.sha256"
$FileHash | Out-File $HashFile

# Generate README.md dynamically
$ReadmePath = "$OutputDir\README.md"
Write-Host "Generating README.md..."

# Start content
$readmeContent = "# $ProjectName $Version`n`n## Downloads`n`n| Package | SHA256 |`n|---------|--------|`n"

# Add current package dynamically
$readmeContent += "| [$ArchiveName](./$ArchiveName) | $FileHash |`n"

# Write README
$readmeContent | Out-File $ReadmePath -Encoding UTF8

Write-Host "Build complete!"
Write-Host "Binary + archive are in $OutputDir"
Write-Host "SHA256 checksum saved to $HashFile"
Write-Host "README generated at $ReadmePath"
