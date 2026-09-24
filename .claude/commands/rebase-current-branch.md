# Rebase Current Branch

Rebase the current git branch onto the latest origin/main, stashing any uncommitted changes and ensuring main is up to date first.

## Instructions

1. Save the current branch name
2. Stash any uncommitted changes if present
3. Fetch the latest changes from origin
4. Checkout main and pull the latest changes
5. Checkout back to the original branch
6. Rebase onto main
7. Pop stashed changes if any were stashed
8. If there are conflicts, inform the user and provide guidance on how to resolve them

## Commands to execute

```bash
# Save current branch name
CURRENT_BRANCH=$(git branch --show-current)

# Stash any uncommitted changes
git stash push -m "Auto-stash before rebase on $CURRENT_BRANCH"

# Fetch latest from origin
git fetch origin

# Update local main branch
git checkout main
git pull origin main
git push origin main

# Return to original branch and rebase
git checkout $CURRENT_BRANCH
git rebase main

# Pop stashed changes
git stash pop
```

## On conflicts

If there are merge conflicts during rebase:
- Show the conflicting files using `git status`
- Explain how to resolve conflicts
- Remind user to run `git rebase --continue` after resolving
- Or `git rebase --abort` to cancel the rebase

If there are conflicts when popping the stash:
- Resolve the conflicts manually
- Run `git stash drop` to remove the stash entry after resolving
