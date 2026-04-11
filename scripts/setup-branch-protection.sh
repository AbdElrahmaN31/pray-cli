#!/bin/sh
# Configure GitHub branch protection rules for the main branch.
# Requires: gh CLI authenticated with admin access to the repository.
#
# Usage: ./scripts/setup-branch-protection.sh

set -e

OWNER="AbdElrahmaN31"
REPO="pray-cli"
BRANCH="main"

echo "Configuring branch protection for $OWNER/$REPO ($BRANCH)..."

gh api \
  --method PUT \
  "repos/$OWNER/$REPO/branches/$BRANCH/protection" \
  --input - <<'EOF'
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["test (1.23)", "test (1.24)", "lint", "build"]
  },
  "enforce_admins": false,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false
}
EOF

echo "Branch protection configured successfully."
echo ""
echo "Rules applied to '$BRANCH':"
echo "  - Require PR with 1 approving review"
echo "  - Dismiss stale reviews on new pushes"
echo "  - Require CI status checks (test, lint, build) to pass"
echo "  - Require branch to be up-to-date before merging"
echo "  - Disallow force pushes and branch deletion"
