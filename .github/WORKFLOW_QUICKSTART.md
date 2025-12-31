# Auto-update Workflow Quick Start Guide

## 🎯 What This Does

Automatically keeps your PKL and pkl-go dependencies up-to-date by:
- Checking for new releases daily
- Creating PRs when updates are available
- Running tests to ensure compatibility
- Providing detailed changelogs

## 🚀 Getting Started

### 1. Verify Workflow is Active

After pushing/merging the workflow:

```bash
# Check if workflow file exists
ls -la .github/workflows/auto-update-dependencies.yml

# View workflow in GitHub
# Go to: https://github.com/kdeps/schema/actions
```

### 2. Test Locally First (Recommended)

Before the first scheduled run, test the logic locally:

```bash
# Run the test script
./scripts/test_workflow_logic.sh
```

**Current Status (as of test):**
```
Updates available:
  • PKL: 0.29.1 → 0.30.2
  • pkl-go: 0.11.1 → 0.12.1
```

### 3. Trigger First Run Manually

Go to GitHub Actions and trigger manually:

1. Visit: https://github.com/kdeps/schema/actions/workflows/auto-update-dependencies.yml
2. Click "Run workflow" (top right)
3. Select branch: `main`
4. Click "Run workflow" button

**Expected:** A PR will be created with title:
```
chore: auto-update PKL to 0.30.2 and pkl-go to 0.12.1
```

### 4. Review the PR

When the PR is created:

✅ **Check PR Details:**
- Title matches format: `chore: auto-update ...`
- Labels: `dependencies`, `autoupdate`
- Body contains changelog
- All files listed are relevant

✅ **Review Changed Files:**
- `versions.json` - versions updated correctly
- `go.mod` - pkl-go version updated
- `go.sum` - checksums updated
- `build.gradle.kts` - PKL plugin version updated
- `.pkl` files - minPklVersion updated
- Dependencies downloaded to `deps/pkl/external/`
- Assets copied to `assets/pkl/external/`

✅ **Verify Tests Pass:**
- All GitHub Actions checks pass
- Go tests successful
- No errors in workflow logs

### 5. Test PR Locally (Optional but Recommended)

```bash
# Checkout the PR
gh pr checkout <PR_NUMBER>

# Or manually
git fetch origin pull/<PR_NUMBER>/head:auto-update-test
git checkout auto-update-test

# Run tests
cd assets && go test -v ./...

# Test PKL evaluation
pkl eval deps/pkl/Tool.pkl --no-cache --format json

# Build
make generate
```

### 6. Merge the PR

Once verified:
1. Approve the PR
2. Merge to main
3. PR branch auto-deletes

### 7. Verify Automated Runs

The workflow will now run daily at 00:00 UTC.

**Monitor:**
- Check Actions tab: https://github.com/kdeps/schema/actions
- View workflow history
- Check for any failures

## 📅 Schedule

**Automatic Runs:**
- Daily at 00:00 UTC (midnight)
- Only creates PR if updates available

**Manual Runs:**
- Anytime via GitHub Actions UI
- Useful for immediate updates

## 🔧 Maintenance

### Monthly Tasks

```bash
# Check workflow status
gh workflow view auto-update-dependencies.yml

# List recent runs
gh run list --workflow=auto-update-dependencies.yml --limit 10

# Check for failed runs
gh run list --workflow=auto-update-dependencies.yml --status=failure
```

### When Updates Fail

1. **Check workflow logs:**
   ```bash
   gh run view <RUN_ID> --log
   ```

2. **Common issues:**
   - API rate limits → Wait and retry
   - Test failures → Review breaking changes
   - Network issues → Transient, retry later

3. **Manual fallback:**
   ```bash
   # Update manually if needed
   ./scripts/update_versions.sh
   ./scripts/update_all.sh
   ```

## 📊 Monitoring Dashboard

Track these on the Actions page:

- ✅ Success rate (should be >95%)
- ⏱️ Run duration (typical: 5-10 minutes)
- 📈 Update frequency (depends on upstream releases)
- 🔄 PRs created vs merged

## 🎓 Learning from PRs

Each auto-generated PR is a learning opportunity:

**What to review:**
1. **Release Notes** - Linked in PR body
2. **Breaking Changes** - Check before merging
3. **New Features** - PKL/pkl-go improvements
4. **Bug Fixes** - Fixes that benefit you

## 🆘 Troubleshooting

### Workflow Not Running

```bash
# Check cron schedule in workflow file
grep cron .github/workflows/auto-update-dependencies.yml

# Verify workflow is enabled
gh workflow list | grep auto-update
```

### No PR Created Despite Updates

Check workflow run logs:
```bash
gh run list --workflow=auto-update-dependencies.yml --limit 1
gh run view <RUN_ID> --log
```

Possible causes:
- No updates available
- Tests failed
- Permission issues

### Multiple PRs Created

If multiple PRs for same update:
1. Close duplicates
2. Keep newest PR
3. Consider adding existing PR check

### Tests Failing

Review the specific test failures:
1. Check PR for test logs
2. Review release notes for breaking changes
3. May need code updates to handle new version

## 🎯 Success Metrics

Your workflow is working well when:

- ✅ Runs complete without errors
- ✅ PRs are created within hours of upstream releases
- ✅ Tests pass consistently
- ✅ Updates merge smoothly
- ✅ No manual intervention needed

## 🔗 Useful Links

- **Workflow File:** `.github/workflows/auto-update-dependencies.yml`
- **Documentation:** `.github/workflows/README.md`
- **Testing Guide:** `.github/WORKFLOW_TESTING.md`
- **Test Script:** `./scripts/test_workflow_logic.sh`

- **Actions Page:** https://github.com/kdeps/schema/actions
- **PKL Releases:** https://github.com/apple/pkl/releases
- **pkl-go Releases:** https://github.com/apple/pkl-go/releases

## 💡 Tips

1. **Review PRs promptly** - Don't let them pile up
2. **Read release notes** - Understand what's changing
3. **Test locally** - For major version updates
4. **Monitor Actions** - Catch failures early
5. **Keep workflow updated** - Update actions versions periodically

## 🎉 You're All Set!

The workflow is now:
- ✅ Committed to main
- ✅ Ready to run
- ✅ Fully documented
- ✅ Tested and working

**Next scheduled run:** Tonight at 00:00 UTC

**Or trigger now:** [Run Workflow](https://github.com/kdeps/schema/actions/workflows/auto-update-dependencies.yml)

---

*For detailed testing procedures, see `.github/WORKFLOW_TESTING.md`*
