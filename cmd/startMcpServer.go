/*
Copyright © 2025 SuprSend
*/
package cmd

import (
	"strings"

	"github.com/mark3labs/mcp-go/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	toolset "suprsend-cli/cmd/tools"
	"suprsend-cli/util"
)

var (
	transport string
	tools     string
)

func getSelectedTools(toolsFlag string) ([]*toolset.Tool, error) {
	// selected tools
	selected := []*toolset.Tool{}
	supportedTools := toolset.GetAllTools()
	// if toolsFlag is "all", return all the tools
	if toolsFlag == "all" {
		return supportedTools, nil
	}
	// get the tools mentioned in toolsFlag
	tools := strings.Split(toolsFlag, ",")

	// if tool name is `type`.* include all the tools that have same type
	for _, tool := range tools {
		if strings.Contains(tool, ".*") {
			toolType := strings.Split(tool, ".*")[0]
			for _, t := range supportedTools {
				if t.Type == toolType {
					selected = append(selected, t)
				}
			}
		} else {
			for _, t := range supportedTools {
				if t.Name == tool {
					selected = append(selected, t)
				}
			}
		}
	}
	return selected, nil
}

// startMcpServerCmd represents the startMcpServer command
var startMcpServerCmd = &cobra.Command{
	Use:   "start-mcp-server",
	Short: "Starts MCP server for SuprSend",
	Long: `Starts the MCP server for SuprSend.
This server will handle all the requests from user about SuprSend capabilities and data.`,
	Run: func(cmd *cobra.Command, args []string) {
		selectedTools, err := getSelectedTools(tools)
		if err != nil {
			log.Fatalf("%v", err)
		}
		// Print a readable string representation of selectedTools
		var toolStrs []string
		for _, t := range selectedTools {
			toolStrs = append(toolStrs, t.Type+":"+t.Name)
		}
		log.Debugf("Selected tools: [%s]", strings.Join(toolStrs, ", "))

		mcpServer := server.NewMCPServer(
			"SuprSend",
			"0.1",
			server.WithResourceCapabilities(true, true),
			server.WithPromptCapabilities(true),
			server.WithToolCapabilities(true),
			server.WithLogging(),
		)

		for _, t := range selectedTools {
			mcpServer.AddTool(t.MCPTool, t.Handler)
		}

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

// Add a subcommand to startMcpServerCmd to list all the tools supported by the server
var listToolsCmd = &cobra.Command{
	Use:   "list-tools",
	Short: "List all the tools supported by the server",
	Run: func(cmd *cobra.Command, args []string) {
		type toolListResponse struct {
			Tool_Type        string `json:"tool_type"`
			Tool_Name        string `json:"tool_name"`
			Tool_Description string `json:"tool_description"`
		}
		var resp []toolListResponse
		for _, t := range toolset.GetAllTools() {
			resp = append(resp, toolListResponse{Tool_Type: t.Type, Tool_Name: t.Name, Tool_Description: t.Description})
		}
		outputType, _ := cmd.Flags().GetString("output")
		util.OutputData(resp, outputType)
	},
}

func init() {
	startMcpServerCmd.AddCommand(listToolsCmd)
	rootCmd.AddCommand(startMcpServerCmd)

	startMcpServerCmd.Flags().StringVarP(&transport, "transport", "t", "stdio", "The transport to use for the MCP server. Can be either 'stdio' or 'sse'.")
	startMcpServerCmd.Flags().StringVarP(&tools, "tools", "T", "all", "The types of tools to use. Can be either 'all' or comma separated list of tool names.")
}
