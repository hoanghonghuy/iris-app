# ============================================================================
# Seed script for Neon PostgreSQL
# ============================================================================

param(
    [string]$DatabaseUrl = $env:NEON_DATABASE_URL,
    [ValidateSet("demo", "master")]
    [string]$SeedType = "demo"
)

$ErrorActionPreference = "Stop"

# Validate database URL
if ([string]::IsNullOrWhiteSpace($DatabaseUrl)) {
    Write-Host "[ERROR] DATABASE_URL not provided" -ForegroundColor Red
    Write-Host ""
    Write-Host "Usage:" -ForegroundColor Yellow
    Write-Host "  .\seed-neon.ps1 -DatabaseUrl 'postgresql://user:pass@host/db?sslmode=require' -SeedType demo"
    Write-Host "  Or set NEON_DATABASE_URL environment variable"
    Write-Host ""
    Write-Host "Seed types:" -ForegroundColor Yellow
    Write-Host "  demo   - Quick demo data (10 students, 3 schools, 3 teachers, 3 parents)"
    Write-Host "  master - Full seed data (56 students, 8 schools, 24 teachers, 24 parents)"
    exit 1
}

# Get script directory and seed files path
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ProjectRoot = Split-Path -Parent (Split-Path -Parent $ScriptDir)

# Determine seed file
if ($SeedType -eq "demo") {
    $SeedFile = Join-Path $ScriptDir "seed_demo.sql"
} else {
    $SeedFile = Join-Path $ScriptDir "seed_master.sql"
}

# Validate seed file
if (-not (Test-Path $SeedFile)) {
    Write-Host "[ERROR] Seed file not found: $SeedFile" -ForegroundColor Red
    exit 1
}

Write-Host "Iris Seed Tool for Neon" -ForegroundColor Cyan
Write-Host "==========================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Seed file: $SeedFile" -ForegroundColor Gray
Write-Host "Database: Neon PostgreSQL" -ForegroundColor Gray
Write-Host "Seed type: $SeedType" -ForegroundColor Gray
Write-Host ""

# Check if psql is installed
try {
    $null = Get-Command psql -ErrorAction Stop
} catch {
    Write-Host "[ERROR] 'psql' CLI not found" -ForegroundColor Red
    Write-Host ""
    Write-Host "Install PostgreSQL client tools:" -ForegroundColor Yellow
    Write-Host "  - Windows: https://www.postgresql.org/download/windows/"
    Write-Host "  - Or use Docker: docker run --rm -i postgres:16-alpine psql ..." -ForegroundColor Gray
    exit 1
}

Write-Host "[RUNNING] Seeding database..." -ForegroundColor Green
Write-Host ""

# Run seed with psql
try {
    # Set UTF-8 encoding for PowerShell output
    $OutputEncoding = [System.Text.Encoding]::UTF8
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    
    if ($SeedType -eq "master") {
        # Master seed requires running 5 separate files in order
        Write-Host "[INFO] Master seed runs 5 files sequentially..." -ForegroundColor Gray
        Write-Host ""
        
        $SeedFiles = @(
            "seed_01_schools_classes.sql",
            "seed_02_users_profiles.sql",
            "seed_03_attendance_health.sql",
            "seed_04_posts_interactions.sql",
            "seed_05_appointments_chat_audit.sql"
        )
        
        $Step = 1
        foreach ($File in $SeedFiles) {
            $FilePath = Join-Path $ScriptDir $File
            Write-Host "[$Step/5] Running $File..." -ForegroundColor Cyan
            
            Get-Content $FilePath -Encoding UTF8 | & psql "$DatabaseUrl&client_encoding=UTF8" --set=client_encoding=UTF8
            
            if ($LASTEXITCODE -ne 0) {
                Write-Host ""
                Write-Host "[ERROR] Failed at step $Step ($File)" -ForegroundColor Red
                exit $LASTEXITCODE
            }
            
            $Step++
        }
    } else {
        # Demo seed is a single file
        Get-Content $SeedFile -Encoding UTF8 | & psql "$DatabaseUrl&client_encoding=UTF8" --set=client_encoding=UTF8
    }
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "[SUCCESS] Seed completed successfully!" -ForegroundColor Green
        Write-Host ""
        
        if ($SeedType -eq "demo") {
            Write-Host "Demo data summary:" -ForegroundColor Cyan
            Write-Host "  - 3 schools" -ForegroundColor Gray
            Write-Host "  - 5 classes" -ForegroundColor Gray
            Write-Host "  - 10 students" -ForegroundColor Gray
            Write-Host "  - 3 teachers" -ForegroundColor Gray
            Write-Host "  - 3 parents" -ForegroundColor Gray
            Write-Host "  - Attendance, health logs, posts, appointments, chat, audit logs" -ForegroundColor Gray
            Write-Host ""
            Write-Host "Demo accounts (password: 123456):" -ForegroundColor Cyan
            Write-Host "  - Super Admin:  admin@iris.local" -ForegroundColor Yellow
            Write-Host "  - School Admin: school-admin@iris.local" -ForegroundColor Yellow
            Write-Host "  - Teacher:      teacher1@iris.local, teacher2@iris.local, teacher3@iris.local" -ForegroundColor Yellow
            Write-Host "  - Parent:       parent1@iris.local, parent2@iris.local, parent3@iris.local" -ForegroundColor Yellow
        } else {
            Write-Host "Master data summary:" -ForegroundColor Cyan
            Write-Host "  - 8+ schools" -ForegroundColor Gray
            Write-Host "  - 56+ classes" -ForegroundColor Gray
            Write-Host "  - 56+ students" -ForegroundColor Gray
            Write-Host "  - 24+ teachers" -ForegroundColor Gray
            Write-Host "  - 24+ parents" -ForegroundColor Gray
            Write-Host "  - Full dataset with attendance, health, posts, appointments, chat, audit" -ForegroundColor Gray
            Write-Host ""
            Write-Host "Demo accounts (password: 123456):" -ForegroundColor Cyan
            Write-Host "  - Super Admin:  admin@iris.local" -ForegroundColor Yellow
            Write-Host "  - School Admins: school-admin-1@iris.local to school-admin-8@iris.local" -ForegroundColor Yellow
            Write-Host "  - Teachers:      teacher-1@iris.local to teacher-24@iris.local" -ForegroundColor Yellow
            Write-Host "  - Parents:       parent-1@iris.local to parent-24@iris.local" -ForegroundColor Yellow
        }
    } else {
        Write-Host ""
        Write-Host "[ERROR] Seed failed with exit code: $LASTEXITCODE" -ForegroundColor Red
        exit $LASTEXITCODE
    }
} catch {
    Write-Host ""
    Write-Host "[ERROR] Seed error: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Next steps:" -ForegroundColor Cyan
Write-Host "  - Update your .env with NEON_DATABASE_URL" -ForegroundColor Gray
Write-Host "  - Start API: cd apps/api/cmd/api && go run ." -ForegroundColor Gray
Write-Host "  - Start Web: cd apps/web && npm run dev" -ForegroundColor Gray
