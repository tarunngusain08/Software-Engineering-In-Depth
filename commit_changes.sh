#!/bin/zsh

# Check if a commit message is provided
if [ $# -eq 0 ]; then
    echo "Please provide a base commit message"
    exit 1
fi

BASE_MESSAGE="$1"

# Change to the root of the git repository
cd "$(git rev-parse --show-toplevel)"

# Get list of changed files, excluding those in .gitignore
CHANGED_FILES=($(git ls-files -m -o --exclude-standard))

# Commit each file separately
for file in "${CHANGED_FILES[@]}"; do
    # Stage only this file
    git add "$file"
    
    # Commit with a specific message
    git commit -m "$BASE_MESSAGE: ${file##*/}"
    
    echo "Committed changes for $file"
done

# Push all commits
git push origin HEAD

echo "All changes committed and pushed successfully!"
