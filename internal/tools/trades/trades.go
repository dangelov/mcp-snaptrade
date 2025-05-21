package trades

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"dangelov.com/snaptrade-mcp/internal/snaptradeclient"
)

var Tool = mcp.NewTool("place_order",
	mcp.WithDescription("Place an order with your brokerage account"),
	mcp.WithString("brokerage",
		mcp.Required(),
		mcp.Description("The brokerage to place an order with"),
		mcp.Enum("Trading212", "Vanguard", "Schwab", "Alpaca", "Alpaca Paper", "Tradier", "Robinhood", "Fidelity", "ETrade"),
	),
	mcp.WithString("action",
		mcp.Required(),
		mcp.Description("The action to perform (BUY/SELL)"),
		mcp.Enum("BUY", "SELL"),
	),
	mcp.WithString("ticker",
		mcp.Required(),
		mcp.Description("The ticker symbol of the stock"),
	),
	mcp.WithNumber("quantity",
		mcp.Required(),
		mcp.Description("The quantity of shares to buy/sell"),
	),
)

func Handler(cl *snaptradeclient.SnapTradeClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		brokerage, ok := request.Params.Arguments["brokerage"].(string)
		if !ok {
			return mcp.NewToolResultError("Invalid brokerage"), nil
		}
		action, ok := request.Params.Arguments["action"].(string)
		if !ok {
			return mcp.NewToolResultError("Invalid action"), nil
		}
		ticker, ok := request.Params.Arguments["ticker"].(string)
		if !ok {
			return mcp.NewToolResultError("Invalid ticker"), nil
		}
		quantity, ok := request.Params.Arguments["quantity"].(float64)
		if !ok {
			return mcp.NewToolResultError("Invalid quantity"), nil
		}
		quantityFloat := float32(quantity)

		// Use the client method
		response, err := cl.ListUserAccounts()
		if err != nil {
			fmt.Println("Error retrieving accounts:", err)
			return mcp.NewToolResultError("Error retrieving accounts"), nil // Return error to MCP
		}
		if len(response) == 0 {
			return mcp.NewToolResultError("No brokerage accounts connected"), nil
		}

		accountId := ""
		for _, account := range response {
			if account.GetInstitutionName() == brokerage {
				accountId = account.Id
				break // Found the account
			}
		}
		if accountId == "" {
			return mcp.NewToolResultError(fmt.Sprintf("No matching account found for brokerage: %s", brokerage)), nil
		}

		// Use the client method
		orderRecord, err := cl.PlaceForceOrder(accountId, action, ticker, quantityFloat)
		if err != nil {
			fmt.Println("Error placing order:", err)
			return mcp.NewToolResultError(fmt.Sprintf("Error placing order: %s", err.Error())), nil // Return error to MCP
		}

		return mcp.NewToolResultText(fmt.Sprintf("Order placed successfully: %s. You can monitor the status of your order by asking me to show your recent orders.", *orderRecord.BrokerageOrderId)), nil
	}
}
