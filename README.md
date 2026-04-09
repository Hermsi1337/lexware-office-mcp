# Lexware Office MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for the [Lexware Office API](https://developers.lexware.io/docs/), written in Go. Connect AI assistants like Claude, Cursor, or any MCP-compatible client directly to your Lexware Office account to manage contacts, create invoices, and more.

## Why?

Lexware Office is one of the most popular cloud accounting platforms in Germany. This MCP server lets AI assistants interact with your Lexware Office data through a safe, typed interface -- no manual API calls, no copy-pasting between tools.

**Key features:**

- **Typed tools only** -- every tool has a strict input/output schema, no raw API passthrough
- **Automatic rate limiting** -- built-in retry logic for Lexware API throttling (HTTP 429)
- **Stdio transport** -- works out of the box with Claude Code, Claude Desktop, Cursor, and other MCP clients
- **Minimal footprint** -- single binary, no external services required

## Available Tools

| Tool | Description |
|------|-------------|
| `lexware_get_profile` | Fetch the Lexware Office profile for the configured API token |
| `lexware_create_simple_contact` | Create a customer contact with a name and optional reference |
| `lexware_create_invoice` | Create an invoice with line items, tax, and payment conditions |

## Prerequisites

- **Lexware API token** -- generate one at [app.lexware.de/addons/public-api](https://app.lexware.de/addons/public-api)

## Installation

### Pre-built binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/Hermsi1337/lexware-office-mcp/releases):

```bash
# Example for Linux amd64
curl -sL https://github.com/Hermsi1337/lexware-office-mcp/releases/latest/download/lexware-office-mcp_linux_amd64.tar.gz | tar xz
```

Binaries are available for Linux, macOS, and Windows on both amd64 and arm64.

### Docker

```bash
docker run --rm -e LEXWARE_API_TOKEN=your-token ghcr.io/hermsi1337/lexware-office-mcp
```

### Build from source

Requires Go 1.22+:

```bash
git clone https://github.com/Hermsi1337/lexware-office-mcp.git
cd lexware-office-mcp
make build
```

## Configuration

Copy the example environment file and fill in your token:

```bash
cp .env.example .env
```

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LEXWARE_API_TOKEN` | Yes | -- | Your private Lexware API token |
| `LEXWARE_BASE_URL` | No | `https://api.lexware.io` | API base URL |
| `LEXWARE_USER_AGENT` | No | `lexware-office-mcp/0.1.0` | User-Agent header sent to Lexware |
| `LEXWARE_FINALIZE_INVOICES` | No | `false` | Automatically finalize invoices on creation |

## MCP Client Setup

Add the server to your MCP client configuration. The examples below show common setups.

### Claude Code

```bash
claude mcp add lexware-office -- /absolute/path/to/bin/lexware-office-mcp
```

Set the environment variable `LEXWARE_API_TOKEN` before launching Claude Code, or pass it inline.

### Claude Desktop / Cursor / Generic MCP Client

Add this to your MCP client configuration file (e.g. `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "/absolute/path/to/bin/lexware-office-mcp",
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```

## Usage Examples

### Fetch your profile

```json
{
  "name": "lexware_get_profile",
  "arguments": {}
}
```

### Create a contact

```json
{
  "name": "lexware_create_simple_contact",
  "arguments": {
    "name": "Max Mustermann",
    "sourceReference": "Order #12345"
  }
}
```

### Create an invoice

```json
{
  "name": "lexware_create_invoice",
  "arguments": {
    "invoice": {
      "voucherDate": "2026-04-07",
      "address": {
        "name": "Max Mustermann",
        "street": "Musterstrasse 1",
        "city": "Freiburg",
        "zip": "79098",
        "countryCode": "DE"
      },
      "lineItems": [
        {
          "type": "custom",
          "name": "Consulting",
          "quantity": 1,
          "unitName": "service",
          "unitPrice": {
            "currency": "EUR",
            "grossAmount": 119.0,
            "taxRatePercentage": 19
          }
        }
      ],
      "totalPrice": { "currency": "EUR" },
      "taxConditions": { "taxType": "gross" },
      "paymentConditions": {
        "paymentTermLabel": "Payable immediately",
        "paymentTermDuration": 0
      },
      "shippingConditions": { "shippingType": "none" }
    },
    "finalize": false
  }
}
```

## Project Structure

Follows the [golang-standards/project-layout](https://github.com/golang-standards/project-layout) convention:

```
cmd/lexware-office-mcp/              # Application entrypoint
internal/lexware/                    # API client, config, types, workflows
internal/server/                     # MCP server and tool registration
build/goreleaser/.goreleaser.yml     # GoReleaser configuration
build/package/docker/                # Dockerfiles
.github/workflows/                   # CI/CD pipelines
```

## Releasing

Releases are automated via [GoReleaser](https://goreleaser.com/) and GitHub Actions. Push a tag to trigger a release:

```bash
git tag v0.1.0
git push origin v0.1.0
```

This builds multi-platform binaries, creates a GitHub Release, and pushes multi-arch Docker images to `ghcr.io/hermsi1337/lexware-office-mcp`.

To test the release process locally (requires Docker):

```bash
make release-check     # Validate GoReleaser config
make release-snapshot  # Build without publishing
```

## Roadmap

- [ ] More contact operations (list, get, update)
- [ ] Article management tools
- [ ] Voucher and file workflows
- [ ] Event subscription support
- [ ] Unit and integration tests

## Links

- [Lexware Office API Documentation](https://developers.lexware.io/docs/)
- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)

## License

See [LICENSE](LICENSE) for details.
