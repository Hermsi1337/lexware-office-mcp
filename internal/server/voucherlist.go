package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
)

type listVouchersInput struct {
	VoucherType   string `json:"voucherType,omitempty" jsonschema:"Comma-separated voucher types: salesinvoice, salescreditnote, purchaseinvoice, purchasecreditnote"`
	VoucherStatus string `json:"voucherStatus,omitempty" jsonschema:"Comma-separated statuses: open, paid, paidoff, voided, transferred, sepadebit, unchecked"`
	Page          int    `json:"page,omitempty" jsonschema:"Page number for pagination (0-based)"`
	Size          int    `json:"size,omitempty" jsonschema:"Results per page (1-250, default 250)"`
}

func (s *Server) registerVoucherlistTools() {
	mcp.AddTool(s.Server, &mcp.Tool{
		Name:        "lexware_list_vouchers",
		Description: "List vouchers across all document types with filters for voucher type and status. Returns a unified view of invoices, credit notes, and purchase documents with contact references. Note: no search by recipient name is available; to find vouchers for a specific contact, first look up the contact via lexware_list_contacts.",
	}, s.listVouchers)
}

func (s *Server) listVouchers(ctx context.Context, _ *mcp.CallToolRequest, input listVouchersInput) (*mcp.CallToolResult, map[string]any, error) {
	filter := lexware.VoucherlistFilter{
		VoucherType:   input.VoucherType,
		VoucherStatus: input.VoucherStatus,
		Page:          input.Page,
		Size:          input.Size,
	}

	result, err := s.client.ListVouchers(ctx, filter)
	return s.workflowResult("list vouchers", result, err)
}
