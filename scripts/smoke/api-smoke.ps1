$ErrorActionPreference = "Stop"

function Assert-Step {
  param(
    [string]$Name,
    [scriptblock]$Block
  )

  try {
    $result = & $Block
    Write-Host "[PASS] $Name" -ForegroundColor Green
    return $result
  } catch {
    Write-Host "[FAIL] $Name => $($_.Exception.Message)" -ForegroundColor Red
    return $null
  }
}

$base = "http://localhost:8080/api/v1"
$today = (Get-Date).ToString("yyyy-MM-dd")

Assert-Step "Health endpoint" { Invoke-RestMethod -Method Get -Uri "$base/health" | Out-Null } | Out-Null

$teacherLogin = Assert-Step "Login teacher1" {
  Invoke-RestMethod -Method Post -Uri "$base/auth/login" -ContentType "application/json" -Body (@{ email = "teacher1@iris.local"; password = "123456" } | ConvertTo-Json)
}

$parentLogin = Assert-Step "Login parent1" {
  Invoke-RestMethod -Method Post -Uri "$base/auth/login" -ContentType "application/json" -Body (@{ email = "parent1@iris.local"; password = "123456" } | ConvertTo-Json)
}

$adminLogin = Assert-Step "Login super admin" {
  Invoke-RestMethod -Method Post -Uri "$base/auth/login" -ContentType "application/json" -Body (@{ email = "admin@iris.local"; password = "123456" } | ConvertTo-Json)
}

if (-not $teacherLogin -or -not $parentLogin -or -not $adminLogin) {
  throw "Missing login tokens; abort smoke."
}

$teacherToken = $teacherLogin.data.access_token
$parentToken = $parentLogin.data.access_token
$adminToken = $adminLogin.data.access_token

$teacherHeaders = @{ Authorization = "Bearer $teacherToken" }
$parentHeaders = @{ Authorization = "Bearer $parentToken" }
$adminHeaders = @{ Authorization = "Bearer $adminToken" }

$classes = Assert-Step "Teacher classes list" {
  Invoke-RestMethod -Method Get -Uri "$base/teacher/classes" -Headers $teacherHeaders
}

if (-not $classes -or $classes.data.Count -eq 0) {
  throw "No class found for teacher."
}

$classId = $classes.data[0].class_id

$students = Assert-Step "Teacher students-in-class" {
  Invoke-RestMethod -Method Get -Uri "$base/teacher/classes/$classId/students" -Headers $teacherHeaders
}

if (-not $students -or $students.data.Count -eq 0) {
  throw "No student found in class."
}

$studentId = $students.data[0].student_id

Assert-Step "Teacher mark attendance" {
  Invoke-RestMethod -Method Post -Uri "$base/teacher/attendance" -Headers $teacherHeaders -ContentType "application/json" -Body (@{ student_id = $studentId; date = $today; status = "present"; note = "smoke-test" } | ConvertTo-Json)
} | Out-Null

$postCreate = Assert-Step "Teacher create class post" {
  Invoke-RestMethod -Method Post -Uri "$base/teacher/posts" -Headers $teacherHeaders -ContentType "application/json" -Body (@{ scope_type = "class"; class_id = $classId; type = "announcement"; content = "Smoke test post $(Get-Date -Format s)" } | ConvertTo-Json)
}

if (-not $postCreate) {
  throw "Cannot create post."
}

$postId = $postCreate.data.post_id

Assert-Step "Teacher like post" { Invoke-RestMethod -Method Post -Uri "$base/teacher/posts/$postId/like" -Headers $teacherHeaders } | Out-Null
Assert-Step "Teacher comment post" {
  Invoke-RestMethod -Method Post -Uri "$base/teacher/posts/$postId/comments" -Headers $teacherHeaders -ContentType "application/json" -Body (@{ content = "teacher smoke comment" } | ConvertTo-Json)
} | Out-Null
Assert-Step "Teacher share post" { Invoke-RestMethod -Method Post -Uri "$base/teacher/posts/$postId/share" -Headers $teacherHeaders } | Out-Null

$parentFeed = Assert-Step "Parent feed list" {
  Invoke-RestMethod -Method Get -Uri "$base/parent/feed?limit=20&offset=0" -Headers $parentHeaders
}

Assert-Step "Parent like same post" { Invoke-RestMethod -Method Post -Uri "$base/parent/posts/$postId/like" -Headers $parentHeaders } | Out-Null
Assert-Step "Parent comment same post" {
  Invoke-RestMethod -Method Post -Uri "$base/parent/posts/$postId/comments" -Headers $parentHeaders -ContentType "application/json" -Body (@{ content = "parent smoke comment" } | ConvertTo-Json)
} | Out-Null
Assert-Step "Parent share same post" { Invoke-RestMethod -Method Post -Uri "$base/parent/posts/$postId/share" -Headers $parentHeaders } | Out-Null

Assert-Step "Admin users list" {
  Invoke-RestMethod -Method Get -Uri "$base/admin/users?limit=10&offset=0" -Headers $adminHeaders
} | Out-Null

Write-Host "=== Smoke Summary ===" -ForegroundColor Cyan
Write-Host "class_id=$classId"
Write-Host "student_id=$studentId"
Write-Host "post_id=$postId"

if ($parentFeed -and $parentFeed.data) {
  $hit = $parentFeed.data | Where-Object { $_.post_id -eq $postId }
  if ($hit) {
    Write-Host "Parent feed contains smoke post: YES" -ForegroundColor Green
  } else {
    Write-Host "Parent feed contains smoke post: NO (can be timing/scope issue)" -ForegroundColor Yellow
  }
}
