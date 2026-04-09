# Lexware Office MCP Server

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for the [Lexware Office API](https://developers.lexware.io/docs/), written in Go. Connect AI assistants like Claude, Cursor, Codex, or Windsurf directly to your Lexware Office account to manage contacts, create invoices, and more.

## Why?

Lexware Office is one of the most popular cloud accounting platforms in Germany. This MCP server lets AI assistants interact with your Lexware Office data through a safe, typed interface -- no manual API calls, no copy-pasting between tools.

**Key features:**

- **Typed tools only** -- every tool has a strict input/output schema, no raw API passthrough
- **Automatic rate limiting** -- built-in retry logic for Lexware API throttling (HTTP 429)
- **Stdio transport** -- works out of the box with Claude, Cursor, Codex, Windsurf, and other MCP clients
- **Minimal footprint** -- single binary or Docker image, no external services required
- **Multi-platform** -- pre-built binaries for Linux, macOS, and Windows (amd64 + arm64)

## Available Tools

| Tool | Description |
|------|-------------|
| `lexware_get_profile` | Fetch the Lexware Office profile for the configured API token |
| `lexware_create_simple_contact` | Create a customer contact with a name and optional reference |
| `lexware_create_invoice` | Create an invoice with line items, tax, and payment conditions |

## Prerequisites

You need a **Lexware API token**. Generate one at [app.lexware.de/addons/public-api](https://app.lexware.de/addons/public-api).

## Installation

Choose one of the following methods.

### Pre-built binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/Hermsi1337/lexware-office-mcp/releases):

```bash
# Linux (amd64)
curl -sL https://github.com/Hermsi1337/lexware-office-mcp/releases/latest/download/lexware-office-mcp_linux_amd64.tar.gz | tar xz

# macOS (Apple Silicon)
curl -sL https://github.com/Hermsi1337/lexware-office-mcp/releases/latest/download/lexware-office-mcp_darwin_arm64.tar.gz | tar xz

# Windows (amd64) -- download and extract the zip
```

### Docker

No installation needed. The MCP client starts the container automatically (see [Client Setup](#client-setup) below). To verify manually:

```bash
docker run --rm -e LEXWARE_API_TOKEN=your-token ghcr.io/hermsi1337/lexware-office-mcp
```

Multi-arch image supporting `linux/amd64` and `linux/arm64`.

### Build from source

Requires Go 1.26+:

```bash
git clone https://github.com/Hermsi1337/lexware-office-mcp.git
cd lexware-office-mcp
make build
```

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `LEXWARE_API_TOKEN` | Yes | -- | Your private Lexware API token |
| `LEXWARE_BASE_URL` | No | `https://api.lexware.io` | API base URL |
| `LEXWARE_USER_AGENT` | No | `lexware-office-mcp/0.1.0` | User-Agent header sent to Lexware |
| `LEXWARE_FINALIZE_INVOICES` | No | `false` | Automatically finalize invoices on creation |

All variables are passed as environment variables. When using Docker, the MCP client forwards them to the container automatically.

## Client Setup

Ready-to-use configuration files for all major MCP clients are provided in the [`example/`](example/) directory. Copy the relevant file, insert your API token, and adjust the binary path if needed.

> **Binary or Docker?** Each client below shows both options. The binary approach is simpler; the Docker approach requires no installation and keeps your system clean.

---

### Claude Desktop

<details>
<summary><strong>Binary</strong></summary>

Add to your `claude_desktop_config.json` ([full example](example/claude-desktop.json)):

- macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
- Windows: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "/absolute/path/to/lexware-office-mcp",
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```
</details>

<details>
<summary><strong>Docker</strong></summary>

Add to your `claude_desktop_config.json` ([full example](example/claude-desktop-docker.json)):

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "LEXWARE_API_TOKEN",
        "ghcr.io/hermsi1337/lexware-office-mcp:latest"
      ],
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```
</details>

---

### Claude Code

<details>
<summary><strong>Binary</strong></summary>

```bash
claude mcp add lexware-office \
  -e LEXWARE_API_TOKEN=your-private-api-key \
  -- /absolute/path/to/lexware-office-mcp
```
</details>

<details>
<summary><strong>Docker</strong></summary>

```bash
claude mcp add lexware-office \
  -e LEXWARE_API_TOKEN=your-private-api-key \
  -- docker run -i --rm -e LEXWARE_API_TOKEN ghcr.io/hermsi1337/lexware-office-mcp:latest
```
</details>

---

### Cursor

<details>
<summary><strong>Binary</strong></summary>

Add to `.cursor/mcp.json` in your project or `~/.cursor/mcp.json` globally ([full example](example/cursor.json)):

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "/absolute/path/to/lexware-office-mcp",
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```
</details>

<details>
<summary><strong>Docker</strong></summary>

Same location, using Docker ([full example](example/cursor-docker.json)):

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "LEXWARE_API_TOKEN",
        "ghcr.io/hermsi1337/lexware-office-mcp:latest"
      ],
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```
</details>

---

### OpenAI Codex

<details>
<summary><strong>Binary</strong></summary>

Add to `~/.codex/config.toml` or `.codex/config.toml` in your project ([full example](example/codex.toml)):

```toml
[mcp_servers.lexware-office]
command = "/absolute/path/to/lexware-office-mcp"

[mcp_servers.lexware-office.env]
LEXWARE_API_TOKEN = "your-private-api-key"
```
</details>

<details>
<summary><strong>Docker</strong></summary>

Same location, using Docker ([full example](example/codex-docker.toml)):

```toml
[mcp_servers.lexware-office]
command = "docker"
args = ["run", "-i", "--rm", "-e", "LEXWARE_API_TOKEN", "ghcr.io/hermsi1337/lexware-office-mcp:latest"]

[mcp_servers.lexware-office.env]
LEXWARE_API_TOKEN = "your-private-api-key"
```
</details>

---

### Windsurf

<details>
<summary><strong>Binary</strong></summary>

Add to `~/.codeium/windsurf/mcp_config.json` ([full example](example/windsurf.json)):

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "/absolute/path/to/lexware-office-mcp",
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```
</details>

<details>
<summary><strong>Docker</strong></summary>

Same location, using Docker ([full example](example/windsurf-docker.json)):

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "LEXWARE_API_TOKEN",
        "ghcr.io/hermsi1337/lexware-office-mcp:latest"
      ],
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key"
      }
    }
  }
}
```
</details>

---

## Tool Examples

Once connected, your AI assistant can call the following tools.

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
example/                            # Ready-to-use MCP client configs
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
