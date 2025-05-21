# Demo: MCP Server for SnapTrade Integration with Claude AI

This GitHub repository showcases a demonstration project for integrating an MCP (Model Context Protocol) server with SnapTrade, enabling a connection to Anthropic's Claude AI. It illustrates how financial account data, accessed via SnapTrade, can be made available to Claude through the MCP framework.

Project Overview:

This proof-of-concept demonstrates a basic pipeline:

- SnapTrade Integration: Connects to SnapTrade's API to access financial account information.
- MCP Server Implementation: Exposes SnapTrade functionalities and data through the Model Context Protocol.
- Claude Connectivity: Allows Claude, or other MCP-compatible AI models, to interact with the data exposed by the MCP server.

The goal is to provide a simple example of how an MCP server can be built to bridge AI models with third-party APIs like SnapTrade.

## Important Setup Instructions

To run this demo, you will need to:

- [Sign up for SnapTrade](https://dashboard.snaptrade.com/signup): You must create your own SnapTrade account to obtain a unique Client ID and Client Secret.
- Register a Demo User: Utilize the user registration process available on SnapTrade's API [documentation demo page](https://docs.snaptrade.com/demo/getting-started). This will provide you with the necessary user credentials.
- Configure Environment Variables: Once you have your SnapTrade Client ID, Client Secret, and registered user details, populate these values in the `.env` file located in the build output folder of this project. See the `.env.example` file for an example.

You will also need to configure your Claude Desktop app to load the new MCP server, like so:

```
{
  "mcpServers": {
    "SnapTrade": {
      "command": "/home/username/projects/mcp-snaptrade/bin/cli"
    }
  }
}
```

Please Note:

- **USE AT YOUR OWN RISK**
- Demo Purposes Only: This project is intended as a demonstration and a learning tool. It is not designed for production use and may not cover all edge cases or possess robust error handling.
- Personal Project & Likely Unmaintained: This repository is a personal project. As such, it is unlikely to receive regular updates, bug fixes, or active maintenance. Please use it as-is with this understanding.

I hope this example provides a helpful starting point for understanding how MCP servers can facilitate connections between AI models and external data services like SnapTrade!
