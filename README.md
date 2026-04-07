# Lexware Office MCP

A lightweight MCP server in Go for the [Lexware Office / Lexware API](https://developers.lexware.io/docs/).

The server uses `stdio`, authenticates with a private Lexware API token, and currently ships with a practical MVP toolset:

- `lexware_get_profile`
- `lexware_list_contacts`
- `lexware_get_contact`
- `lexware_create_contact`
- `lexware_update_contact`
- `lexware_list_articles`
- `lexware_get_article`
- `lexware_get_invoice`
- `lexware_list_vouchers`
- `lexware_api_request`

The generic `lexware_api_request` tool is intentionally included so the server remains useful before every Lexware endpoint has a dedicated typed wrapper.

## Requirements

- Go 1.22 or newer
- A private Lexware API token from `https://app.lexware.de/addons/public-api`

## Configuration

Example:

```bash
cp .env.example .env
```

Important environment variables:

- `LEXWARE_API_TOKEN`: required
- `LEXWARE_BASE_URL`: defaults to `https://api.lexware.io`
- `LEXWARE_USER_AGENT`: optional, defaults to `lexware-office-mcp/0.1.0`
- `LEXWARE_MIN_INTERVAL_MS`: optional, defaults to `550`

The default request interval is intentionally conservative because the official Lexware documentation currently states a limit of `2 requests per second`.

## Running

```bash
go mod tidy
go build -o bin/lexware-office-mcp ./cmd/lexware-office-mcp
LEXWARE_API_TOKEN=... ./bin/lexware-office-mcp
```

## MCP Client Example

For a local MCP client over `stdio`, the configuration is conceptually:

```json
{
  "mcpServers": {
    "lexware-office": {
      "command": "/absolute/path/to/bin/lexware-office-mcp",
      "env": {
        "LEXWARE_API_TOKEN": "your-private-api-key",
        "LEXWARE_BASE_URL": "https://api.lexware.io"
      }
    }
  }
}
```

## Examples

Fetch the current profile:

```json
{
  "name": "lexware_get_profile",
  "arguments": {}
}
```

Create a contact:

```json
{
  "name": "lexware_create_contact",
  "arguments": {
    "payload": {
      "version": 0,
      "roles": {
        "customer": {}
      },
      "person": {
        "firstName": "Inge",
        "lastName": "Musterfrau"
      }
    }
  }
}
```

Call any Lexware endpoint:

```json
{
  "name": "lexware_api_request",
  "arguments": {
    "method": "GET",
    "path": "/v1/countries"
  }
}
```

## Next Useful Extensions

- Dedicated tools for `invoices`, `vouchers`, `files`, and `event subscriptions`
- Better endpoint-specific input schemas instead of generic payload maps
- Optional OAuth or multi-tenant support if the project later needs to connect multiple Lexware accounts

## Sources

- [Lexware API Documentation](https://developers.lexware.io/docs/)
- [Model Context Protocol Go SDK](https://github.com/modelcontextprotocol/go-sdk)
