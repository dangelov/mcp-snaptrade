package portfolio

import (
	"context"
	"fmt"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	st "github.com/passiv/snaptrade-sdks/sdks/go"

	"dangelov.com/snaptrade-mcp/internal/snaptradeclient"
)

var Tool = mcp.NewTool("portfolio",
	mcp.WithDescription("Check your portfolio and brokerage accounts for their positions and values."),
	mcp.WithString("brokerage",
		mcp.Required(),
		mcp.Description("The brokerage to connect to"),
		mcp.Enum("Trading212", "Vanguard", "Schwab", "Alpaca", "Alpaca Paper", "Tradier", "Robinhood", "Fidelity", "ETrade"),
	),
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
			fmt.Println("No brokerage accounts connected")
			return mcp.NewToolResultError("No brokerage accounts connected"), nil
		}

		reply := "Your connected brokerage accounts are:\n"

		// Group the accounts by institution
		accounts := make(map[string][]st.Account)
		for _, account := range response {
			institution := account.GetInstitutionName()
			if _, ok := accounts[institution]; !ok {
				accounts[institution] = []st.Account{}
			}
			accounts[institution] = append(accounts[institution], account)
		}

		// Find all accounts
		for institution, accountsList := range accounts {
			reply += fmt.Sprintf("%s:\n", institution)
			for _, account := range accountsList {
				marketValueDescription := ""
				if account.Balance.GetTotal().Amount != nil && account.Balance.GetTotal().Currency != nil {
					marketValue := *account.Balance.GetTotal().Amount
					marketValueCurrency := *account.Balance.GetTotal().Currency
					marketValueDescription = fmt.Sprintf("(value including cash: %.2f %s)", marketValue, marketValueCurrency)
				}
				reply += fmt.Sprintf("%s %s\n\nPlease show this markdown formatted table of all the positions under this account\n\n", *account.Name.Get(), marketValueDescription)

				// Now get all the positions for this account using the client method
				positions, err := cl.GetUserAccountPositions(account.Id)
				if err != nil {
					fmt.Println("Error retrieving positions for account", account.Name, ":", err)
					continue // Or add an error message to the reply
				}
				if len(positions) == 0 {
					reply += fmt.Sprintf("  - No positions found for account %s\n", *account.Name.Get())
					continue
				} else {
					tableData := make([][]string, len(positions))
					for i, position := range positions {
						tableData[i] = []string{
							position.Symbol.Symbol.Symbol,
							fmt.Sprintf("%.4f", *position.Units.Get()),
							fmt.Sprintf("%.2f %s", *position.Price.Get()**position.Units.Get(), *position.Symbol.Symbol.Currency.Code),
						}
					}

					basicTable, err := markdown.NewTableFormatterBuilder().
						Build("Instrument", "Units", "Value").
						Format(tableData)
					if err != nil {
						fmt.Println("Error creating table:", err)
						continue // Or add an error message to the reply
					}
					// Add the table to the reply
					reply += fmt.Sprintf("%s\n", basicTable)
				}
			}
		}

		fmt.Println(reply)
		return mcp.NewToolResultText(reply), nil
	}
}
