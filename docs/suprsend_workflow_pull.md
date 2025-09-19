## suprsend workflow pull

Pull workflows from workspace to local directory

### Synopsis

pull workflows from workspace to local directory of each workflow

```
suprsend workflow pull [flags]
```

### Options

```
  -h, --help                help for pull
  -m, --mode string         Mode of workflows to pull from (default "live")
  -d, --output-dir string   Output directory for workflows
```

### Options inherited from parent commands

```
      --config string          config file (default: $HOME/.suprsend.yaml)
  -n, --no-color               Disable color output (default: $NO_COLOR)
  -o, --output string          Output Style (pretty, yaml, json) (default "pretty")
  -s, --service-token string   Service token (default: $SUPRSEND_SERVICE_TOKEN)
  -v, --verbosity string       Log level (debug, info, warn, error, fatal, panic) (default "info")
  -w, --workspace string       Workspace to use (default "staging")
```

### SEE ALSO

* [suprsend workflow](suprsend_workflow.md)	 - Manage workflows

