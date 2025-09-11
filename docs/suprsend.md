## suprsend

CLI to interact with SuprSend, a Notification Infrastructure

### Synopsis

SuprSend is a robust notification infrastructure that helps you deploy multi-channel product notifications effortlessly and take care of user experience.

This CLI lets you interact with your SuprSend workspace and do actions like fetching/modifying template, workflows etc.

### Options

```
      --config string          config file (default: $HOME/.suprsend.yaml)
  -h, --help                   help for suprsend
  -n, --no-color               Disable color output (default: $NO_COLOR)
  -o, --output string          Output Tyle (pretty, yaml, json) (default "pretty")
  -s, --service-token string   Service token (default: $SUPRSEND_SERVICE_TOKEN)
  -v, --verbosity string       Log level (debug, info, warn, error, fatal, panic) (default "info")
  -w, --workspace string       Workspace to use (default "staging")
```

### SEE ALSO

* [suprsend category](suprsend_category.md)	 - Manage preference categories
* [suprsend completion](suprsend_completion.md)	 - Generate the autocompletion script for the specified shell
* [suprsend generate-types](suprsend_generate-types.md)	 - Generate Types from JSON Schema
* [suprsend profile](suprsend_profile.md)	 - Manage Profile
* [suprsend schema](suprsend_schema.md)	 - Manage schema
* [suprsend start-mcp-server](suprsend_start-mcp-server.md)	 - Starts MCP server for SuprSend
* [suprsend sync](suprsend_sync.md)	 - Sync all your SuprSend assets locally
* [suprsend version](suprsend_version.md)	 - Print the CLI version
* [suprsend workflow](suprsend_workflow.md)	 - Manage workflows

