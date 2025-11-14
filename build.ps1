param (
    [string]$ProjectName    = "fsearch",
    [string]$CliPath        = "./cli",
    [string]$Owner          = "UmbrellaCrow612",            
    [string]$OutputDir      = "dist",
    [string]$BuildDir       = "build"
)

# Get version from git, fallback if not available
try {
    $Version = git describe --tags --always | ForEach-Object { $_.Trim() }
} catch {
    Write-Warning "Couldn't get version from git. Falling back to v0.0.1"
    $Version = "v0.0.1"
}

Write-Host "Version to release: $Version"

# Tag & push if not already tagged
$TagExists = git tag --list $Version
if (-not $TagExists) {
    Write-Host "Creating git tag $Version"
    git tag $Version
    git push origin $Version
} else {
    Write-Host "Tag $Version already exists. Skipping tag creation."
}

# Build targets
$Targets = @(
    @{ OS="windows"; ARCH="amd64" },
    @{ OS="windows"; ARCH="386" },
    @{ OS="linux";   ARCH="amd64" },
    @{ OS="linux";   ARCH="386" },
    @{ OS="darwin";  ARCH="amd64" },
    @{ OS="darwin";  ARCH="arm64" }
)

# Prepare directories
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null
New-Item -ItemType Directory -Force -Path $BuildDir  | Out-Null

$AssetList = @()  # To keep track of built ZIPs

foreach ($target in $Targets) {
    $OS   = $target.OS
    $ARCH = $target.ARCH

    Write-Host "Building $ProjectName for $OS/$ARCH..."
    $ext        = if ($OS -eq "windows") { ".exe" } else { "" }
    $BinaryName = "$ProjectName$ext"
    $BinaryPath = Join-Path $BuildDir $BinaryName

    $env:GOOS   = $OS
    $env:GOARCH = $ARCH

    go build -o $BinaryPath $CliPath
    if (-Not (Test-Path $BinaryPath)) {
        Write-Error "Build failed for $OS/$ARCH"
        exit 1
    }

    # Package binary
    $ArchiveName = "${ProjectName}_${Version}_${OS}_${ARCH}.zip"
    $ArchivePath = Join-Path $OutputDir $ArchiveName
    Compress-Archive -Path $BinaryPath -DestinationPath $ArchivePath -Force

    Write-Host "Built and packaged $ArchiveName"

    # Add archive to asset list
    $AssetList += $ArchivePath
}

Write-Host "All builds complete! Packages are in $OutputDir"

# Create GitHub release
$Repo = "$Owner/$ProjectName"
$ReleaseTag   = $Version
$ReleaseTitle = "$ProjectName $Version"
$ReleaseNotes = "Automatic build & release of version $Version"

Write-Host "Creating GitHub release $ReleaseTag in repo $Repo"
gh release create $ReleaseTag --repo $Repo --title "$ReleaseTitle" --notes "$ReleaseNotes"

# Upload assets
foreach ($file in $AssetList) {
    Write-Host "Uploading asset: $([System.IO.Path]::GetFileName($file))"
    gh release upload $ReleaseTag --repo $Repo $file --clobber
}

Write-Host "Release flow completed successfully for version $Version"
