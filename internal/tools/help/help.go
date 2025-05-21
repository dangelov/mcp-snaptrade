package help

import (
	"context"

	"dangelov.com/snaptrade-mcp/internal/snaptradeclient"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var Tool = mcp.NewTool("get_started_with_brokerage_connection",
	mcp.WithDescription("Provides information on how to connect your brokerage account and lists the supported brokerages."),
)

func Handler(cl *snaptradeclient.SnapTradeClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Note: This specific handler doesn't currently *use* the client,
		// but we adapt the signature for consistency and future use.
		return mcp.NewToolResultText("To get started with investing and portfolio management, please let us know which brokerage you have an account with. We can help you connect your account to any of the following brokerages: Trading212, Vanguard, Schwab, Alpaca, Alpaca Paper, Tradier, Robinhood, Fidelity, ETrade."), nil
	}
}
