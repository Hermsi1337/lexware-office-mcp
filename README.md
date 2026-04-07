# Lexware Office MCP

A lightweight MCP server in Go for the [Lexware Office / Lexware API](https://developers.lexware.io/docs/).

The server uses `stdio`, authenticates with a private Lexware API token, and currently ships with a practical MVP toolset:

- `lexware_get_profile`
- `lexware_list_contacts`
- `lexware_get_contact`
- `lexware_create_contact`
- `lexware_create_simple_contact`
- `lexware_update_contact`
- `lexware_list_articles`
- `lexware_get_article`
- `lexware_get_invoice`
- `lexware_create_invoice`
- `lexware_list_vouchers`
- `lexware_api_request`

The generic `lexware_api_request` tool is intentionally included so the server remains useful before every Lexware endpoint has a dedicated typed wrapper.

The repository also includes an initial set of typed workflow coverage:

- profile retrieval
- simple contact creation
- invoice creation with the `finalize` query parameter

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
- `LEXWARE_FINALIZE_INVOICES`: optional, defaults to `false`

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

Create a simple contact using the legacy helper:

```json
{
  "name": "lexware_create_simple_contact",
  "arguments": {
    "name": "Max Mustermann",
    "sourceReference": "cardmarket orderID: 12345"
  }
}
```

Create an invoice using the legacy integration shape:

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
      "totalPrice": {
        "currency": "EUR"
      },
      "taxConditions": {
        "taxType": "gross"
      },
      "paymentConditions": {
        "paymentTermLabel": "Payable immediately",
        "paymentTermDuration": 0
      },
      "shippingConditions": {
        "shippingType": "none"
      }
    },
    "finalize": false
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
