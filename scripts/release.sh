#!/bin/bash

# Release script for Todoy
# Usage: ./scripts/release.sh <version>
# Example: ./scripts/release.sh v1.0.0

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    echo "Example: $0 v1.0.0"
    exit 1
fi

VERSION=$1

# Validate version format (should start with 'v')
if [[ ! $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    echo "Error: Version should follow semantic versioning format (e.g., v1.0.0, v1.0.0-beta)"
    exit 1
fi

echo "Preparing release $VERSION..."

# Check if working directory is clean
if [[ -n $(git status --porcelain) ]]; then
    echo "Error: Working directory is not clean. Please commit or stash your changes."
    exit 1
fi

# Make sure we're on the main/master branch
CURRENT_BRANCH=$(git branch --show-current)
if [[ "$CURRENT_BRANCH" != "master" && "$CURRENT_BRANCH" != "main" ]]; then
    echo "Warning: You're not on the master/main branch. Current branch: $CURRENT_BRANCH"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Pull latest changes
echo "Pulling latest changes..."
git pull origin $CURRENT_BRANCH

# Update CHANGELOG.md (move unreleased to version)
echo "Updating CHANGELOG.md..."
TODAY=$(date +%Y-%m-%d)
sed -i.bak "s/## \[Unreleased\]/## [Unreleased]\n\n## [$VERSION] - $TODAY/" CHANGELOG.md
rm CHANGELOG.md.bak

# Commit changelog update
git add CHANGELOG.md
git commit -m "Update CHANGELOG for $VERSION release"

# Create and push tag
echo "Creating tag $VERSION..."
git tag -a $VERSION -m "Release $VERSION"
git push origin $CURRENT_BRANCH
git push origin $VERSION

echo "Release $VERSION has been tagged and pushed!"
echo "GitHub Actions will now build and create the release automatically."
echo "Check the Actions tab in your GitHub repository for progress."
