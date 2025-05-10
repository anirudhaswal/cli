package tools

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func searchDocsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

func fetchDocsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uri, ok := request.Params.Arguments["uri"].(string)
	if !ok {
		return nil, errors.New("uri must be a string")
	}
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

func newDocumentationTools() []*Tool {
	searchDoc := &Tool{
		Name:        "documentation.search",
		Description: "Enables querying SuprSend documentation",
		MCPTool: mcp.NewTool("search_suprsend_documentation",
			mcp.WithDescription(`Use this tool whenever you need technical guidance or answers related to SuprSend. It is particularly useful when:\n- A user asks about SuprSend's capabilities, features, or integrations (e.g., Workflows, Templates, Tenants, Lists, Vendors, Connectors, etc.)\n- You're unsure how a specific SuprSend functionality or API works\n- You're writing or debugging an integration with SuprSend\nHow to use:\n- Frame your queries using precise technical terms (avoid vague phrasing)\n- This tool returns a JSON array in text format — each item contains:\n  - uri: the documentation path\n  - snippet: a short excerpt relevant to the query\nAnswering process:\n1. Review all snippets in order to try to answer the question.\n2. If by using all the snippets, you still don't have enough information to confidently answer:\n  - Use the corresponding uri to fetch full documentation content vy using the fetch_suprsend_documentation tool.\n3. Process each resource one-by-one in the order returned.\n4. Use information from both snippets and full documentation to construct the final answer.`),
			mcp.WithString("query",
				mcp.Description(`Search query. The query should: 1. Identify the core concepts and intent 2. Add relevant synonyms and related terms 3. Structure the query to emphasize key terms 4. Include technical or domain-specific terminology if applicable`),
				mcp.Required(),
			),
		),
		Handler: searchDocsHandler,
	}

	fetchDoc := &Tool{
		Name:        "documentation.fetch",
		Description: "Fetch the full documentation content for the given uri.",
		MCPTool: mcp.NewTool("fetch_suprsend_documentation",
			mcp.WithDescription(`Use this tool to fetch the full documentation content for the given uri.`),
			mcp.WithString("uri",
				mcp.Description(`The uri of the documentation to fetch.`),
				mcp.Required(),
			),
		),
		Handler: fetchDocsHandler,
	}

	return []*Tool{searchDoc, fetchDoc}
}

func init() {
	for _, t := range newDocumentationTools() {
		RegisterTool(t, "documentation")
	}
}
