package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/dennis/lexware-office-mcp/internal/lexware"
	"github.com/dennis/lexware-office-mcp/internal/server"
)

func main() {
	cfg, err := lexware.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	client := lexware.NewClient(cfg)
	srv := server.New(client)

	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
