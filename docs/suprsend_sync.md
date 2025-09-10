## suprsend sync

Sync all your SuprSend assets locally

### Synopsis

Sync all your SuprSend assets locally with the server

```
suprsend sync [flags]
```

### Options

```
  -a, --assets string   Assets to sync (all, workflows, schemas) (default "all")
  -f, --from string     Source workspace (required) (default "staging")
  -h, --help            help for sync
  -m, --mode string     Mode to sync assets (draft, live) (default "live")
  -t, --to string       Destination workspace (required) (default "production")
```

### Options inherited from parent commands

```
      --config string          config file (default: $HOME/.suprsend.yaml)
  -n, --no-color               Disable color output (default: $NO_COLOR)
  -o, --output string          Output Tyle (pretty, yaml, json) (default "pretty")
  -s, --service-token string   Service token (default: $SUPRSEND_SERVICE_TOKEN)
  -v, --verbosity string       Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [suprsend](suprsend.md)	 - CLI to interact with SuprSend, a Notification Infrastructure

