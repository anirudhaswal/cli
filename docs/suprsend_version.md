## suprsend version

Print the CLI version

```
suprsend version [flags]
```

### Examples

```

/Users/gaurav/Library/Caches/go-build/8c/8c3755e3687779410a1634bee45f2150f0413ce12d31b91f42292eb5784aff91-d/main version
/Users/gaurav/Library/Caches/go-build/8c/8c3755e3687779410a1634bee45f2150f0413ce12d31b91f42292eb5784aff91-d/main version -o=json
/Users/gaurav/Library/Caches/go-build/8c/8c3755e3687779410a1634bee45f2150f0413ce12d31b91f42292eb5784aff91-d/main version -o=yaml
/Users/gaurav/Library/Caches/go-build/8c/8c3755e3687779410a1634bee45f2150f0413ce12d31b91f42292eb5784aff91-d/main version -o=short

```

### Options

```
  -h, --help            help for version
  -o, --output string   Output format. One of: json | pretty | short | yaml (default "pretty")
```

### Options inherited from parent commands

```
      --config string          config file (default: $HOME/.suprsend.yaml)
  -n, --no-color               Disable color output (default: $NO_COLOR)
  -s, --service-token string   Service token (default: $SUPRSEND_SERVICE_TOKEN)
  -v, --verbosity string       Log level (debug, info, warn, error, fatal, panic) (default "info")
```

### SEE ALSO

* [suprsend](suprsend.md)	 - CLI to interact with SuprSend, a Notification Infrastructure

