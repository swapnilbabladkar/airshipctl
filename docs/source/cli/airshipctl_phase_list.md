## airshipctl phase list

List phases

### Synopsis

List life-cycle phases which were defined in document model by group.
Phases within a group are executed sequentially. Multiple phase groups
are executed in parallel.


```
airshipctl phase list PHASE_NAME [flags]
```

### Examples

```

# List phases of phasePlan
airshipctl phase list --plan phasePlan

# To output the contents to table (default operation)
airshipctl phase list --plan phasePlan -o table

# To output the contents to yaml
airshipctl phase list --plan phasePlan -o yaml

# List all phases
airshipctl phase list

# List phases with clustername
airshipctl phase list --cluster-name clustername

```

### Options

```
  -c, --cluster-name string   filter documents by cluster name
  -h, --help                  help for list
  -o, --output string         'table' and 'yaml' are available output formats (default "table")
      --plan string           Plan name of a plan
```

### Options inherited from parent commands

```
      --airshipconf string   Path to file for airshipctl configuration. (default "$HOME/.airship/config")
      --debug                enable verbose output
```

### SEE ALSO

* [airshipctl phase](airshipctl_phase.md)	 - Manage phases

