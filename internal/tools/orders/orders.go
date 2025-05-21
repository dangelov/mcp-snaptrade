package orders

import (
	"context"
	"fmt"
	"sort"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	st "github.com/passiv/snaptrade-sdks/sdks/go"

	"dangelov.com/snaptrade-mcp/internal/snaptradeclient"
)

var Tool = mcp.NewTool("orders",
	mcp.WithDescription("Check your recent orders across all your brokerage accounts."),
)

func Handler(cl *snaptradeclient.SnapTradeClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Use the client method
		response, err := cl.ListUserAccounts()
		if err != nil {
			fmt.Println("Error retrieving accounts:", err)
			return mcp.NewToolResultError("Error retrieving accounts"), nil
		}
		if len(response) == 0 {
			return mcp.NewToolResultError("No brokerage accounts connected"), nil
		}

		all_orders := []st.AccountOrderRecord{}
		for _, account := range response {
			// Now get all the orders for this account using the client method
			orders, err := cl.GetUserAccountOrders(account.Id)
			if err != nil {
				fmt.Println("Error retrieving orders for account", account.Name, ":", err)
				// Optionally add a message to the user or just skip
				continue
			}
			all_orders = append(all_orders, orders...)
		}

		sort.Slice(all_orders, func(i, j int) bool {
			// Handle potential nil TimePlaced
			timeI := all_orders[i].TimePlaced
			timeJ := all_orders[j].TimePlaced
			if timeI == nil && timeJ == nil {
				return false // Consider them equal if both nil
			}
			if timeI == nil {
				return false // Nil times go last
			}
			if timeJ == nil {
				return true // Nil times go last
			}
			return timeI.After(*timeJ)
		})

		if len(all_orders) == 0 {
			return mcp.NewToolResultText("No recent orders found."), nil
		}
		tableData := make([][]string, len(all_orders))
		for i, order := range all_orders {
			timePlacedStr := "N/A"
			if order.TimePlaced != nil {
				timePlacedStr = order.TimePlaced.Format("2006-01-02 15:04:05")
			}
			priceStr := "N/A"
			if order.ExecutionPrice.IsSet() && order.UniversalSymbol != nil {
				priceStr = fmt.Sprintf("%.2f %s", order.GetExecutionPrice(), *order.UniversalSymbol.Currency.Code)
			} else if order.ExecutionPrice.IsSet() {
				priceStr = fmt.Sprintf("%.2f", order.GetExecutionPrice())
			}

			tableData[i] = []string{
				order.UniversalSymbol.GetSymbol(),
				order.GetAction(),
				string(order.GetStatus()),
				timePlacedStr,
				fmt.Sprintf("%.4f", order.GetTotalQuantity()),
				priceStr,
			}
		}
		table, err := markdown.NewTableFormatterBuilder().
			Build("Instrument", "Order Type", "Status", "Time Placed", "Quantity", "Price").
			Format(tableData)
		if err != nil {
			fmt.Println("Error creating table:", err)
			return mcp.NewToolResultError("Error formatting order results"), nil
		}
		return mcp.NewToolResultText(table), nil
	}
}
