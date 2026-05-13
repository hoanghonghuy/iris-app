# ============================================================================
# Migration script for Neon PostgreSQL
# ============================================================================

param(
    [string]$DatabaseUrl = $env:NEON_DATABASE_URL,
    [switch]$Down,
    [int]$Steps = 0
)

$ErrorActionPreference = "Stop"

# Validate database URL
if ([string]::IsNullOrWhiteSpace($DatabaseUrl)) {
    Write-Host "[ERROR] DATABASE_URL not provided" -ForegroundColor Red
    Write-Host ""
    Write-Host "Usage:" -ForegroundColor Yellow
    Write-Host "  .\migrate-neon.ps1 -DatabaseUrl 'postgresql://user:pass@host/db?sslmode=require'"
    Write-Host "  Or set NEON_DATABASE_URL environment variable"
    exit 1
}

# Get script directory and migrations path
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $ScriptDir)
$MigrationsPath = Join-Path $ProjectRoot "apps\api\migrations"
# Convert to forward slashes for migrate CLI
$MigrationsPath = $MigrationsPath -replace '\\', '/'

# Validate migrations directory
if (-not (Test-Path $MigrationsPath)) {
    Write-Host "[ERROR] Migrations directory not found: $MigrationsPath" -ForegroundColor Red
    exit 1
}

Write-Host "Iris Migration Tool for Neon" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Migrations path: $MigrationsPath" -ForegroundColor Gray
Write-Host "Database: Neon PostgreSQL" -ForegroundColor Gray
Write-Host ""

# Check if migrate CLI is installed
try {
    $null = Get-Command migrate -ErrorAction Stop
} catch {
    Write-Host "[ERROR] 'migrate' CLI not found" -ForegroundColor Red
    Write-Host ""
    Write-Host "Install it with:" -ForegroundColor Yellow
    Write-Host "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
}

# Build migrate command
$MigrateArgs = @(
    "-path", $MigrationsPath,
    "-database", $DatabaseUrl
)

if ($Down) {
    if ($Steps -gt 0) {
        $MigrateArgs += "down", $Steps
        Write-Host "[DOWN] Running DOWN migration ($Steps steps)..." -ForegroundColor Yellow
    } else {
        Write-Host "[WARNING] -Down without -Steps will rollback ALL migrations" -ForegroundColor Yellow
        $Confirm = Read-Host "Are you sure? (yes/no)"
        if ($Confirm -ne "yes") {
            Write-Host "[CANCELLED]" -ForegroundColor Red
            exit 0
        }
        $MigrateArgs += "down"
        Write-Host "[DOWN] Running DOWN migration (all)..." -ForegroundColor Yellow
    }
} elseif ($Steps -gt 0) {
    $MigrateArgs += "up", $Steps
    Write-Host "[UP] Running UP migration ($Steps steps)..." -ForegroundColor Green
} else {
    $MigrateArgs += "up"
    Write-Host "[UP] Running UP migration (all pending)..." -ForegroundColor Green
}

Write-Host ""

# Run migration
try {
    & migrate $MigrateArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "[SUCCESS] Migration completed successfully!" -ForegroundColor Green
    } else {
        Write-Host ""
        Write-Host "[ERROR] Migration failed with exit code: $LASTEXITCODE" -ForegroundColor Red
        exit $LASTEXITCODE
    }
} catch {
    Write-Host ""
    Write-Host "[ERROR] Migration error: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  - Run seed: .\seed-neon.ps1 -DatabaseUrl 'your-neon-url'" -ForegroundColor Gray
Write-Host "  - Check version: migrate -path $MigrationsPath -database 'your-url' version" -ForegroundColor Gray
