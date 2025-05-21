package connect

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"dangelov.com/snaptrade-mcp/internal/snaptradeclient"
)

var Tool = mcp.NewTool("connect_brokerage",
	mcp.WithDescription("Connect your brokerage account to see your portfolio and trades for that account."),
	mcp.WithString("brokerage",
		mcp.Required(),
		mcp.Description("The brokerage to connect to"),
		mcp.Enum("Trading212", "Vanguard", "Schwab", "Alpaca", "Alpaca Paper", "Tradier", "Robinhood", "Fidelity", "ETrade"),
	),
)

func Handler(cl *snaptradeclient.SnapTradeClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		brokerage, ok := request.Params.Arguments["brokerage"].(string)
		if !ok {
			return mcp.NewToolResultError("Invalid brokerage"), nil
		}

		nameToSlug := map[string]string{
			"Trading212":   "TRADING212",
			"Vanguard":     "VANGUARD",
			"Schwab":       "SCHWAB",
			"Alpaca":       "ALPACA",
			"Alpaca Paper": "ALPACA-PAPER",
			"Tradier":      "TRADIER",
			"Robinhood":    "ROBINHOOD",
			"Fidelity":     "FIDELITY",
			"ETrade":       "ETRADE",
		}
		slug, ok := nameToSlug[brokerage]
		if !ok {
			return mcp.NewToolResultError("Invalid brokerage"), nil
		}

		redirectURI, err := cl.LoginUserAndGetRedirectURI(slug)
		if err != nil {
			return mcp.NewToolResultError("Failed to generate connection link for SnapTrade"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("To connect to %s, you must present this link to the user. Please note that it expires, so even if you've shown it before, you need to show this new one:\n %s", brokerage, redirectURI)), nil
	}
}
