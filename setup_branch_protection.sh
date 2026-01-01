#!/bin/bash

# Branch Protection Setup Script for GitHub
# This script sets up branch protection rules for the main branch

set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Repository information
REPO_OWNER="leventeberry"
REPO_NAME="goapi"
BRANCH="main"

echo "=========================================="
echo "  GitHub Branch Protection Setup"
echo "=========================================="
echo ""
echo "Repository: $REPO_OWNER/$REPO_NAME"
echo "Branch: $BRANCH"
echo ""

# Check if GitHub CLI is installed
if ! command -v gh &> /dev/null; then
    echo -e "${RED}Error: GitHub CLI (gh) is not installed${NC}"
    echo ""
    echo "Please install GitHub CLI:"
    echo "  Windows: winget install GitHub.cli"
    echo "  Mac:     brew install gh"
    echo "  Linux:   See https://github.com/cli/cli/blob/trunk/docs/install_linux.md"
    echo ""
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    echo -e "${YELLOW}Not authenticated with GitHub CLI${NC}"
    echo "Please run: gh auth login"
    exit 1
fi

echo -e "${GREEN}✓ GitHub CLI is installed and authenticated${NC}"
echo ""

# Confirm before proceeding
read -p "This will set up branch protection for 'main'. Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted."
    exit 1
fi

echo ""
echo "Setting up branch protection rules..."
echo ""

# Set up branch protection
# Note: This uses the GitHub API directly via gh CLI
gh api \
  repos/$REPO_OWNER/$REPO_NAME/branches/$BRANCH/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":[]}' \
  --field enforce_admins=true \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true,"require_code_owner_reviews":false}' \
  --field restrictions=null \
  --field required_linear_history=true \
  --field allow_force_pushes=false \
  --field allow_deletions=false

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ Branch protection rules successfully applied!${NC}"
    echo ""
    echo "Protection rules enabled:"
    echo "  ✓ Require pull request reviews (1 approval)"
    echo "  ✓ Dismiss stale reviews on new commits"
    echo "  ✓ Require branches to be up to date"
    echo "  ✓ Require linear history"
    echo "  ✓ Enforce rules for administrators"
    echo "  ✓ Block force pushes"
    echo "  ✓ Block branch deletion"
    echo ""
    echo "Next steps:"
    echo "  1. Test by trying to push directly to main (should fail)"
    echo "  2. Create a feature branch and open a PR"
    echo "  3. Verify that PR requires approval before merging"
    echo ""
    echo "To view current protection rules:"
    echo "  gh api repos/$REPO_OWNER/$REPO_NAME/branches/$BRANCH/protection"
    echo ""
    echo "Or visit: https://github.com/$REPO_OWNER/$REPO_NAME/settings/branches"
else
    echo ""
    echo -e "${RED}✗ Failed to set up branch protection${NC}"
    echo ""
    echo "Possible reasons:"
    echo "  - Insufficient permissions (need admin access)"
    echo "  - Repository doesn't exist or is private"
    echo "  - Network error"
    echo ""
    echo "You can set up branch protection manually:"
    echo "  https://github.com/$REPO_OWNER/$REPO_NAME/settings/branches"
    exit 1
fi

