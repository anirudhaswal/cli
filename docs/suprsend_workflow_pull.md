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
  -m, --mode string         Mode of workflows to pull from (draft, live), default: live (default "live")
  -d, --output-dir string   Output directory for workflows
  -g, --slug string         Slug of the workflow to pull
```

### Options inherited from parent commands

```
      --config string          config file (default: $HOME/.suprsend.yaml)
  -n, --no-color               Disable color output (default: $NO_COLOR)
  -s, --service-token string   Service token (default: $SUPRSEND_SERVICE_TOKEN)
  -v, --verbosity string       Log level (debug, info, warn, error, fatal, panic) (default "info")
  -w, --workspace string       Workspace to use (default "staging")
```

### SEE ALSO

* [suprsend workflow](suprsend_workflow.md)	 - Manage workflows

