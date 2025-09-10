## suprsend schema list

List schemas

### Synopsis

List schemas in a workspace

```
suprsend schema list [flags]
```

### Options

```
  -h, --help          help for list
  -l, --limit int     Limit the number of schemas to list (default 20)
  -m, --mode string   Mode of schemas to list (default "live")
  -f, --offset int    Offset the number of schemas to list (default: 0)
```

### Options inherited from parent commands

```
      --config string          config file (default: $HOME/.suprsend.yaml)
  -n, --no-color               Disable color output (default: $NO_COLOR)
  -o, --output string          Output Tyle (pretty, yaml, json) (default "pretty")
  -s, --service-token string   Service token (default: $SUPRSEND_SERVICE_TOKEN)
  -v, --verbosity string       Log level (debug, info, warn, error, fatal, panic) (default "info")
  -w, --workspace string       Workspace to use (default "staging")
```

### SEE ALSO

* [suprsend schema](suprsend_schema.md)	 - Manage schema

