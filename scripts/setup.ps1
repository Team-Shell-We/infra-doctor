Write-Host "Configuring Git..."

git config core.hooksPath .githooks
git config commit.template .gitmessage

Write-Host ""
Write-Host "Done!"
Write-Host "[OK] Git Hooks configured"
Write-Host "[OK] Commit template configured"