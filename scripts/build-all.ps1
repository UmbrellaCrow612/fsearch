# root/scripts/build-all.ps1
# -------------------
# Assumes project structure:
# root/
#   cli/        <- main.go lives here
#   scripts/    <- this script
#   bin/        <- output folder (created automatically)

$ErrorActionPreference = "Stop"

# Targets: OS and ARCH
$targets = @(
    @{OS="windows"; ARCH="amd64"},
    @{OS="windows"; ARCH="arm64"},
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="darwin"; ARCH="amd64"},
    @{OS="darwin"; ARCH="arm64"}
)

# Resolve root relative to script location
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$root = Resolve-Path "$scriptDir/.."
$cliDir = Join-Path $root "cli"
$binDir = Join-Path $root "bin"

# Clean bin folder
if (Test-Path $binDir) { Remove-Item $binDir -Recurse -Force }
New-Item -ItemType Directory -Path $binDir | Out-Null

foreach ($t in $targets) {
    $os = $t.OS
    $arch = $t.ARCH
    Write-Host "Building for $os/$arch..."

    $outputName = if ($os -eq "windows") { "mycli.exe" } else { "mycli" }
    $outDir = Join-Path $binDir "$os-$arch"
    New-Item -ItemType Directory -Path $outDir | Out-Null

    # Build
    & go build -o (Join-Path $outDir $outputName) $cliDir

    # Zip the build
    $zipName = "mycli-$os-$arch.zip"
    $zipPath = Join-Path $binDir $zipName
    if (Test-Path $zipPath) { Remove-Item $zipPath }

    Add-Type -AssemblyName System.IO.Compression.FileSystem
    [System.IO.Compression.ZipFile]::CreateFromDirectory($outDir, $zipPath)

    # Compute SHA256
    $hash = (Get-FileHash $zipPath -Algorithm SHA256).Hash
    Write-Host "$zipName SHA256: $hash"

    # Create a .sha256 file next to zip
    "$hash  $zipName" | Out-File "$zipPath.sha256" -Encoding ASCII
}

Write-Host "All builds complete! Check the bin/ folder."
