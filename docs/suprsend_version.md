## suprsend version

Print the CLI version

```
suprsend version [flags]
```

### Examples

```

/Users/gaurav/Library/Caches/go-build/4c/4cbe520b26bc85b6507978c00a94507f3b30d2fcfffea8cf04302120a2e2cb10-d/main version
/Users/gaurav/Library/Caches/go-build/4c/4cbe520b26bc85b6507978c00a94507f3b30d2fcfffea8cf04302120a2e2cb10-d/main version -o=json
/Users/gaurav/Library/Caches/go-build/4c/4cbe520b26bc85b6507978c00a94507f3b30d2fcfffea8cf04302120a2e2cb10-d/main version -o=yaml
/Users/gaurav/Library/Caches/go-build/4c/4cbe520b26bc85b6507978c00a94507f3b30d2fcfffea8cf04302120a2e2cb10-d/main version -o=short

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

