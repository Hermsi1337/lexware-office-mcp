# AGENTS.md

## Project Purpose

This repository contains a Go-based MCP server for the Lexware Office public API. The project goal is to provide a maintainable, extensible bridge between MCP clients and Lexware resources such as contacts, articles, invoices, vouchers, files, and related accounting workflows.

Primary external reference:

- Lexware API documentation: https://developers.lexware.io/docs/

## Mandatory Repository Language

All repository-facing content must be written in English.

This rule applies to:

- source code comments
- commit messages
- pull request titles and descriptions
- README content
- AGENTS/CLAUDE guidance files
- issue text added by agents
- generated docs and examples
- test names when readable prose is involved

User prompts may be in German or another language, but the repository output must remain English unless the user explicitly requests a non-English artifact for product reasons.

## Working Agreement For Agents

- Preserve the existing Go codebase and extend it incrementally.
- Prefer small, composable changes over broad rewrites.
- Keep the server usable as an MVP even while dedicated endpoint wrappers are still incomplete.
- Favor typed MCP tools where the schema is stable and valuable.
- Do not add generic raw passthrough MCP tools. This repository should expose typed requests and typed responses only.
- Respect Lexware API limits and keep rate limiting or retry behavior explicit in code.
- Prefer official Lexware API behavior over assumptions; verify unclear details against the official documentation.

## Testing Conventions

- All tests use [testify](https://github.com/stretchr/testify).
- Use `require` exclusively -- never `assert`. Tests must fail immediately on the first violated expectation rather than accumulating soft failures.
- Every test file must be organized as one or more **testify suites** (`suite.Suite`). Do not use bare `Test*` functions.
- Use the suite lifecycle interfaces where applicable: `SetupSuite` / `TearDownSuite` for one-time setup, `SetupTest` / `TearDownTest` for per-test setup and teardown.
- Integration-style tests use `net/http/httptest` to mock the Lexware API. Create the mock server in `SetupTest` and close it in `TearDownTest`.

## Documentation Requirements

Before every commit, ensure both of these files are up to date with the current implementation:

- `README.md`
- `AGENTS.md`

If the code changes behavior, setup, scope, supported tools, constraints, or conventions, update the documentation in the same change before committing.

## Project Layout

The repository follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) convention:

- `cmd/lexware-office-mcp/main.go`: application entrypoint
- `internal/lexware/client.go`: authenticated Lexware HTTP client with 429 retry handling
- `internal/lexware/config.go`: environment-based configuration loading
- `internal/lexware/common_types.go`: shared types (Page, Address, LineItem, TotalPrice, etc.)
- `internal/lexware/{resource}_types.go`: per-resource type definitions
- `internal/lexware/{resource}.go`: per-resource Client methods, filters, result types
- `internal/lexware/{resource}_test.go`: per-resource test suites
- `internal/lexware/test_helpers_test.go`: shared baseSuite for httptest server lifecycle
- `internal/version/version.go`: build-time version injection via ldflags
- `internal/server/server.go`: MCP server shell (Server struct, constructor, result helper)
- `internal/server/{resource}.go`: per-resource input types, handlers, tool registration
- `build/goreleaser/.goreleaser.yml`: GoReleaser configuration for multi-platform releases
- `build/package/docker/`: Dockerfiles for container image builds
- `example/`: ready-to-use MCP client configuration files for Claude, Cursor, Codex, and Windsurf
- `.github/workflows/release.yml`: GitHub Actions workflow triggered by version tags

## Current Tool Surface

The repository exposes these MCP tools:

**Profile:**
- `lexware_get_profile`

**Contacts:**
- `lexware_create_simple_contact`
- `lexware_get_contact`
- `lexware_list_contacts`

**Invoices:**
- `lexware_create_invoice`
- `lexware_get_invoice`

**Articles:**
- `lexware_create_article`
- `lexware_get_article`
- `lexware_list_articles`

**Quotations:**
- `lexware_create_quotation`
- `lexware_get_quotation`

**Credit Notes:**
- `lexware_create_credit_note`
- `lexware_get_credit_note`

**Voucherlist:**
- `lexware_list_vouchers`

**Delivery Notes:**
- `lexware_create_delivery_note`
- `lexware_get_delivery_note`

**Order Confirmations:**
- `lexware_create_order_confirmation`
- `lexware_get_order_confirmation`

**Down Payment Invoices:**
- `lexware_get_down_payment_invoice`

**Recurring Templates:**
- `lexware_get_recurring_template`

**Reference Data:**
- `lexware_list_countries`
- `lexware_list_payment_conditions`
- `lexware_list_posting_categories`

When adding or removing tools, update `README.md` and this file before committing.

## Preferred Next Steps

- Add voucher file upload/download workflows
- Add dunning notice tools
- Add event subscription support
- Consider better error mapping and retry strategy for Lexware API failures

## Release Process

Releases are managed by [GoReleaser](https://goreleaser.com/) and triggered by pushing a git tag matching `v*`:

1. Tag a commit: `git tag v0.1.0 && git push origin v0.1.0`
2. GitHub Actions builds multi-platform binaries (linux/darwin/windows, amd64/arm64)
3. Multi-arch Docker images are pushed to `ghcr.io/hermsi1337/lexware-office-mcp`
4. A GitHub Release is created with binaries and checksums

Local testing via Docker (no GoReleaser install needed):

- `make release-check` -- validate config
- `make release-snapshot` -- build without publishing

## Commit Hygiene

- Keep commits focused and readable.
- Use English commit messages.
- Do not commit local workspace artifacts such as `.codex`, `.gocache`, `.gomodcache`, or IDE metadata.
- Keep `.gitignore` current when new local-only files appear.

## CLAUDE.md

`CLAUDE.md` should remain a symlink to `AGENTS.md` so tools and coding agents that look for either file receive the same repository guidance.
