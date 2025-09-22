## suprsend workflow push

push workflows from local to suprsend

### Synopsis

push workflows from local to suprsend dashboard

```
suprsend workflow push [flags]
```

### Options

```
  -c, --commit                  Commit the workflows
  -m, --commit-message string   Commit message for the workflows
  -h, --help                    help for push
  -p, --path string             Output directory for workflows
  -g, --slug string             Slug of the workflow to push
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

