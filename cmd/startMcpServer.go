/*
Copyright © 2025 SuprSend
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"suprsend-cli/util"
)

var transport string

func search_docs_handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok {
		return nil, errors.New("query must be a string")
	}

	encodedQuery := url.QueryEscape(query)
	response, err := http.Get(fmt.Sprintf("https://rag.suprsend.com/?query=%s", encodedQuery))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(body)), nil
}

func fetch_docs_handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uri, ok := request.Params.Arguments["uri"].(string)
	if !ok {
		return nil, errors.New("uri must be a string")
	}
	// if uri doesn't end with .md, add it
	if !strings.HasSuffix(uri, ".md") {
		uri = uri + ".md"
	}
	url := fmt.Sprintf("https://docs.suprsend.com/%s", uri)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(body)), nil
}

// startMcpServerCmd represents the startMcpServer command
var startMcpServerCmd = &cobra.Command{
	Use:   "start-mcp-server",
	Short: "Starts MCP server for SuprSend",
	Long: `Starts the MCP server for SuprSend.
This server will handle all the requests from user about SuprSend capabilities and data.`,
	Run: func(cmd *cobra.Command, args []string) {
		mcpServer := server.NewMCPServer(
			"SuprSend",
			"0.1",
			server.WithResourceCapabilities(true, true),
			server.WithPromptCapabilities(true),
			server.WithToolCapabilities(true),
			server.WithLogging(),
		)

		doc_search := mcp.NewTool("search_suprsend_documentation",
			mcp.WithDescription(`Use this tool whenever you need technical guidance or answers related to SuprSend. It is particularly useful when:
			- A user asks about SuprSend’s capabilities, features, or integrations (e.g., Workflows, Templates, Tenants, Lists, Vendors, Connectors, etc.)
			- You’re unsure how a specific SuprSend functionality or API works
			- You’re writing or debugging an integration with SuprSend
			How to use:
				- Frame your queries using precise technical terms (avoid vague phrasing)
				- This tool returns a JSON array in text format — each item contains:
					- uri: the documentation path
					- snippet: a short excerpt relevant to the query
			Answering process:
				1.	Review all snippets in order to try to answer the question.
				2.	If by using all the snippets, you still don't have enough information to confidently answer:
					- Use the corresponding uri to fetch full documentation content vy using the fetch_suprsend_documentation tool.
				3.	Process each resource one-by-one in the order returned.
				4.	Use information from both snippets and full documentation to construct the final answer.`),
			mcp.WithString("query",
				mcp.Description(`Search query. The query should: 
				1. Identify the core concepts and intent 
				2. Add relevant synonyms and related terms 
				3. Structure the query to emphasize key terms 
				4. Include technical or domain-specific terminology if applicable`),
				mcp.Required(),
			),
		)

		fetch_doc := mcp.NewTool("fetch_suprsend_documentation",
			mcp.WithDescription(`Use this tool to fetch the full documentation content for the given uri.`),
			mcp.WithString("uri",
				mcp.Description(`The uri of the documentation to fetch.`),
				mcp.Required(),
			))

		mcpServer.AddTool(doc_search, search_docs_handler)
		mcpServer.AddTool(fetch_doc, fetch_docs_handler)
		// Start the server

		// switch between transport types based on the transport flag
		if transport == "sse" {
			util.Banner()
			sseServer := server.NewSSEServer(mcpServer, server.WithBaseURL("http://localhost:8080"))
			log.Printf("SSE server listening on :8080")
			if err := sseServer.Start(":8080"); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		} else {
			if err := server.ServeStdio(mcpServer); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startMcpServerCmd)

	// Add transport flag which can be either "stdio" or "sse", default is "stdio"
	startMcpServerCmd.Flags().StringVarP(&transport, "transport", "t", "stdio", "The transport to use for the MCP server. Can be either 'stdio' or 'sse'.")
}
