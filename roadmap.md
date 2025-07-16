# GOSANDBOX Roadmap

This roadmap outlines the items we need to fix, validate, and test to ensure the sandbox-credential CLI works reliably in all scenarios.

## 1. Environment Configuration
- [ ] Copy `example.env` to `.env` and populate all required variables:
  - `URL`, `USERNAME`, `PASSWORD`, `DOWNLOAD_KEY`
  - Optional overrides: `AWS_RELATIVE_PATH`, `ACLOUD_SELECTOR`
- [ ] Enhance error handling when `.env` is missing or malformed
  - Surface clear errors if required vars are unset

## 2. Browser / Rod Integration
- [ ] Verify Chrome/Chromium installation and `ROD_BROWSER` fallback
- [ ] Test headless operation (e.g., under `xvfb-run` or in Docker)
- [ ] Replace hard-coded `time.Sleep` calls with Rod explicit waits:
  - `WaitVisible`, `WaitEnabled`, or `ElementR` with timeouts
- [ ] Validate `ACLOUD_SELECTOR` override correctly targets the “Start AWS Sandbox” element

## 3. CLI Functionality
- [ ] Non-interactive mode (`--start-service`):
  - Boots, scrapes credentials, prints them, and exits
- [ ] Interactive menu flow — exercise each option:
  1. Scrape sandbox credentials
  2. Download text file of credentials
  3. Append credentials to AWS config
  4. Display credentials
  5. Set GitHub secrets
  6. Open AWS Console
  7. Write/read SQLite table
- [ ] Add flags for:
  - Custom `.env` file path
  - Custom selector (`--selector`)

## 4. Docker / Container Support
- [ ] Build and run with Docker Compose:
  ```sh
  docker-compose up --build
  ```
- [ ] Validate `.env` mounting and `./data` persistence
- [ ] Remove unused `ARG` directives in `Dockerfile` (USERNAME, PASSWORD, etc.)

## 5. Testing & CI
- [ ] Clean up `test/` directory:
  - Example mains marked `// +build ignore`
  - Ensure `go test ./...` succeeds
- [ ] Write unit tests for:
  - `acloud.Sandbox` override (with custom `ACLOUD_SELECTOR`)
  - Core login flows (`core.Login`, `core.SimpleLogin`)
  - AWS config append logic (`core.AppendAwsCredentials`)
  - SQLite repository CRUD operations
- [ ] Add integration tests (stub or real sandbox) in headless mode
- [ ] Configure CI (GitHub Actions) to:
  - Install Chromium
  - Run headless tests under `xvfb` or Docker
  - Run `go fmt`, `go vet`, and `golangci-lint`

## 6. Documentation & Cleanup
- [ ] Update `README.md` with:
  - Service mode usage (`--start-service`)
  - `ACLOUD_SELECTOR` override instructions
  - Docker Compose steps
  - Troubleshooting macOS permissions (`GOCACHE`, Full Disk Access)
- [ ] Remove debug logs (`cli.Success`) or add verbosity flags
- [ ] Add CONTRIBUTING guidelines and license file if open source

## 7. Optional Enhancements
- [ ] Publish to Homebrew with automated CI release
- [ ] Provide pre-built binaries or Docker images
- [ ] Improve retry logic around network/browser failures
  
---
_Last updated: $(date +"%Y-%m-%d")_