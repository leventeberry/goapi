# PowerShell Branch Protection Setup Script for GitHub
# This script sets up branch protection rules for the main branch

$ErrorActionPreference = "Stop"

# Repository information
$REPO_OWNER = "leventeberry"
$REPO_NAME = "goapi"
$BRANCH = "main"

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "  GitHub Branch Protection Setup" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Repository: $REPO_OWNER/$REPO_NAME"
Write-Host "Branch: $BRANCH"
Write-Host ""

# Check if GitHub CLI is installed
try {
    $null = Get-Command gh -ErrorAction Stop
} catch {
    Write-Host "Error: GitHub CLI (gh) is not installed" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install GitHub CLI:"
    Write-Host "  Windows: winget install GitHub.cli"
    Write-Host "  Or download from: https://cli.github.com/"
    Write-Host ""
    exit 1
}

# Check if authenticated
try {
    gh auth status 2>&1 | Out-Null
    if ($LASTEXITCODE -ne 0) {
        throw "Not authenticated"
    }
} catch {
    Write-Host "Not authenticated with GitHub CLI" -ForegroundColor Yellow
    Write-Host "Please run: gh auth login" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ GitHub CLI is installed and authenticated" -ForegroundColor Green
Write-Host ""

# Confirm before proceeding
$confirmation = Read-Host "This will set up branch protection for 'main'. Continue? (y/N)"
if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
    Write-Host "Aborted."
    exit 0
}

Write-Host ""
Write-Host "Setting up branch protection rules..." -ForegroundColor Yellow
Write-Host ""

# Set up branch protection using GitHub CLI
$protectionConfig = @{
    required_status_checks = @{
        strict = $true
        contexts = @()
    }
    enforce_admins = $true
    required_pull_request_reviews = @{
        required_approving_review_count = 1
        dismiss_stale_reviews = $true
        require_code_owner_reviews = $false
    }
    restrictions = $null
    required_linear_history = $true
    allow_force_pushes = $false
    allow_deletions = $false
} | ConvertTo-Json -Depth 10

try {
    gh api "repos/$REPO_OWNER/$REPO_NAME/branches/$BRANCH/protection" `
        --method PUT `
        --input - <<< $protectionConfig

    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✓ Branch protection rules successfully applied!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Protection rules enabled:"
        Write-Host "  ✓ Require pull request reviews (1 approval)"
        Write-Host "  ✓ Dismiss stale reviews on new commits"
        Write-Host "  ✓ Require branches to be up to date"
        Write-Host "  ✓ Require linear history"
        Write-Host "  ✓ Enforce rules for administrators"
        Write-Host "  ✓ Block force pushes"
        Write-Host "  ✓ Block branch deletion"
        Write-Host ""
        Write-Host "Next steps:"
        Write-Host "  1. Test by trying to push directly to main (should fail)"
        Write-Host "  2. Create a feature branch and open a PR"
        Write-Host "  3. Verify that PR requires approval before merging"
        Write-Host ""
        Write-Host "To view current protection rules:"
        Write-Host "  gh api repos/$REPO_OWNER/$REPO_NAME/branches/$BRANCH/protection"
        Write-Host ""
        Write-Host "Or visit: https://github.com/$REPO_OWNER/$REPO_NAME/settings/branches"
    } else {
        throw "Command failed"
    }
} catch {
    Write-Host ""
    Write-Host "✗ Failed to set up branch protection" -ForegroundColor Red
    Write-Host ""
    Write-Host "Possible reasons:"
    Write-Host "  - Insufficient permissions (need admin access)"
    Write-Host "  - Repository doesn't exist or is private"
    Write-Host "  - Network error"
    Write-Host ""
    Write-Host "You can set up branch protection manually:"
    Write-Host "  https://github.com/$REPO_OWNER/$REPO_NAME/settings/branches"
    exit 1
}

