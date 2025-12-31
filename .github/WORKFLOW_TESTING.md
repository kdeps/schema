# Workflow Testing Checklist

## Auto-update Dependencies Workflow

### Pre-merge Testing

Before merging the workflow, verify:

- [ ] YAML syntax is valid (already validated with yamllint)
- [ ] All required secrets are configured
  - [ ] `GITHUB_TOKEN` is available (automatically provided by GitHub)
- [ ] Workflow file is in correct location (`.github/workflows/`)
- [ ] Documentation is complete

### Post-merge Testing

After merging to main:

#### 1. Manual Trigger Test

**Steps:**
1. Go to GitHub repository → Actions tab
2. Click "Auto-update Dependencies" workflow
3. Click "Run workflow" button
4. Select `main` branch
5. Click "Run workflow"

**Expected Results:**
- [ ] Workflow starts successfully
- [ ] "Checkout repository" step completes
- [ ] Go, Java, and Gradle setup complete
- [ ] jq installation succeeds
- [ ] PKL CLI installation succeeds
- [ ] Latest versions fetched from GitHub API
- [ ] Current versions read from `versions.json`
- [ ] Version comparison completes

**If versions are up-to-date:**
- [ ] Workflow exits with "All dependencies are up to date"
- [ ] No PR created
- [ ] Workflow marked as successful

**If updates are available:**
- [ ] Files updated correctly
- [ ] Dependencies downloaded
- [ ] Tests run and pass
- [ ] PR created with correct format
- [ ] PR has proper labels
- [ ] PR body contains all sections

#### 2. Scheduled Run Test

**Wait for next scheduled run (00:00 UTC) or:**
- Change cron schedule to run sooner for testing
- Monitor the automated run

**Expected Results:**
- [ ] Workflow triggers automatically at scheduled time
- [ ] Same behavior as manual trigger
- [ ] Check GitHub Actions history for run

#### 3. PR Validation Test

When a PR is created:

**PR Format:**
- [ ] Title: `chore: auto-update pkl to X.X.X and/or pkl-go to X.X.X`
- [ ] Labels: `dependencies`, `autoupdate`
- [ ] Branch name: `auto-update-pkl-X.X.X-pkl-go-X.X.X`

**PR Body Contains:**
- [ ] "Auto-update Dependencies" header
- [ ] PKL Update section (if applicable)
  - [ ] Version change (old → new)
  - [ ] List of changed files
  - [ ] Release notes link
- [ ] pkl-go Update section (if applicable)
  - [ ] Version change (old → new)
  - [ ] List of changed files
  - [ ] Release notes link
- [ ] Test Results section
- [ ] Footer with workflow attribution

**Changed Files:**
- [ ] `versions.json` - correct versions
- [ ] `go.mod` - correct pkl-go version (if updated)
- [ ] `go.sum` - updated checksums (if pkl-go updated)
- [ ] `build.gradle.kts` - correct PKL version (if updated)
- [ ] `deps/pkl/*.pkl` - correct minPklVersion (if PKL updated)
- [ ] `assets/pkl/*.pkl` - correct minPklVersion (if PKL updated)
- [ ] Dependencies in `deps/pkl/external/`
- [ ] Dependencies in `assets/pkl/external/`

**CI/CD Checks:**
- [ ] All CI checks pass
- [ ] No merge conflicts
- [ ] Ready for review

#### 4. Test Update Application

When reviewing the PR:

**Local Testing:**
```bash
# Checkout the PR branch
gh pr checkout <PR_NUMBER>

# Verify versions
jq '.pkl.version, .dependencies."pkl-go".version' versions.json

# Verify go.mod
grep pkl-go go.mod

# Verify build.gradle.kts
grep 'id("org.pkl-lang")' build.gradle.kts

# Run tests
cd assets && go test -v ./...

# Try PKL evaluation
pkl eval deps/pkl/Tool.pkl --no-cache --format json

# Build
make generate
```

**Verification:**
- [ ] All tests pass locally
- [ ] PKL files evaluate correctly
- [ ] No import errors
- [ ] Build succeeds
- [ ] Assets embedded correctly

#### 5. Merge and Monitor

After merging:
- [ ] PR branch deleted automatically
- [ ] Main branch updated
- [ ] No workflow errors
- [ ] Badge on README shows passing status

### Troubleshooting Guide

#### Workflow Fails at Checkout
**Cause:** Permission issues
**Fix:** Verify repository permissions, check GITHUB_TOKEN

#### Workflow Fails at Version Fetch
**Cause:** GitHub API rate limiting or network issues
**Fix:** Wait and retry, check GitHub API status

#### Workflow Fails at Dependency Download
**Cause:** Network issues, wrong versions, missing dependencies
**Fix:**
- Check `versions.json` format
- Verify network connectivity
- Check if versions exist in upstream repos

#### Tests Fail After Update
**Cause:** Breaking changes in new version
**Fix:**
- Review release notes for breaking changes
- May need manual fixes to code
- Consider pinning to previous version temporarily

#### PR Not Created
**Cause:** No changes detected, permission issues
**Fix:**
- Check if versions actually changed
- Verify `contents: write` and `pull-requests: write` permissions
- Check workflow logs for errors

#### Multiple PRs Created
**Cause:** Workflow ran multiple times before first PR merged
**Fix:**
- Close duplicate PRs
- Merge the appropriate one
- May want to add check for existing PRs

### Performance Monitoring

Track these metrics:
- [ ] Workflow execution time (typical: 5-10 minutes)
- [ ] Frequency of updates (varies by upstream release cadence)
- [ ] Test success rate
- [ ] Time from PR creation to merge

### Security Considerations

- [ ] PRs are reviewed before merging
- [ ] Dependency sources are trusted (apple/pkl, apple/pkl-go)
- [ ] Tests catch potential issues
- [ ] No secrets exposed in workflow logs

### Maintenance

Monthly checks:
- [ ] Workflow still running on schedule
- [ ] No failed runs requiring attention
- [ ] Dependencies are up-to-date
- [ ] Documentation is current

Quarterly review:
- [ ] Update GitHub Actions versions if needed
- [ ] Review and update Go/Java versions
- [ ] Optimize workflow if needed
- [ ] Check for GitHub Actions best practice updates

## Quick Test Commands

```bash
# Check workflow syntax
yamllint .github/workflows/auto-update-dependencies.yml

# View current versions
jq '.' versions.json

# Check for latest pkl version
curl -s https://api.github.com/repos/apple/pkl/releases/latest | jq -r '.tag_name'

# Check for latest pkl-go version
curl -s https://api.github.com/repos/apple/pkl-go/releases/latest | jq -r '.tag_name'

# Test scripts locally
./scripts/update_versions.sh
./scripts/download_deps.sh
./scripts/fix_deps_imports.sh

# Verify offline functionality
pkl eval deps/pkl/Tool.pkl --no-cache

# Run Go tests
cd assets && go test -v ./...
```

## Success Criteria

The workflow is considered successful when:

1. ✅ Runs on schedule without errors
2. ✅ Detects new versions correctly
3. ✅ Creates well-formatted PRs
4. ✅ All tests pass in PRs
5. ✅ Updates are merged without issues
6. ✅ No manual intervention needed for routine updates
7. ✅ Repository stays current with upstream releases

## Rollback Plan

If the workflow causes issues:

1. **Disable Workflow:**
   ```bash
   # Edit .github/workflows/auto-update-dependencies.yml
   # Add to top of file:
   # on:
   #   workflow_dispatch:  # Remove schedule trigger
   ```

2. **Revert Changes:**
   ```bash
   git revert <commit-hash>
   git push origin main
   ```

3. **Manual Updates:**
   - Use existing scripts manually
   - Follow OFFLINE.md documentation
   - Fix issues before re-enabling workflow

4. **Fix and Re-enable:**
   - Address root cause
   - Test fix thoroughly
   - Re-enable scheduled runs
