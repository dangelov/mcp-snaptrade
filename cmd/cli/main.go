package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mark3labs/mcp-go/server"

	"dangelov.com/snaptrade-mcp/internal/snaptradeclient"
	"dangelov.com/snaptrade-mcp/internal/tools/connect"
	"dangelov.com/snaptrade-mcp/internal/tools/help"
	"dangelov.com/snaptrade-mcp/internal/tools/orders"
	"dangelov.com/snaptrade-mcp/internal/tools/portfolio"
	"dangelov.com/snaptrade-mcp/internal/tools/trades"
)

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
	exPath := filepath.Dir(ex)
	envPath := filepath.Join(exPath, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Error loading .env file from %s: %v. Will try to use environment variables directly.", envPath, err)
	}

	// Setup a new SnapTrade client
	stClient := snaptradeclient.New(os.Getenv("SNAPTRADE_ID"), os.Getenv("SNAPTRADE_SECRET"), os.Getenv("USER_ID"), os.Getenv("USER_SECRET"))

	// Create a new MCP server
	s := server.NewMCPServer(
		"SnapTrade",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// Add the tool handlers, injecting the snaptrade client
	s.AddTool(help.Tool, help.Handler(stClient))
	s.AddTool(connect.Tool, connect.Handler(stClient))
	s.AddTool(portfolio.Tool, portfolio.Handler(stClient))
	s.AddTool(orders.Tool, orders.Handler(stClient))
	s.AddTool(trades.Tool, trades.Handler(stClient))

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
