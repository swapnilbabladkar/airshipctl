## airshipctl config set-context

Manage contexts

### Synopsis

Create or modify a context in the airshipctl config files.


```
airshipctl config set-context NAME [flags]
```

### Examples

```

# Create a new context named "exampleContext"
airshipctl config set-context exampleContext \
  --manifest=exampleManifest \

# Update the manifest of the current-context
airshipctl config set-context \
  --current \
  --manifest=exampleManifest

```

### Options

```
      --current                    update the current context
  -h, --help                       help for set-context
      --management-config string   set the management config for the specified context
      --manifest string            set the manifest for the specified context
```

### Options inherited from parent commands

```
      --airshipconf string   Path to file for airshipctl configuration. (default "$HOME/.airship/config")
      --debug                enable verbose output
```

### SEE ALSO

* [airshipctl config](airshipctl_config.md)	 - Manage the airshipctl config file

