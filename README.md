# logmerge
Merges logfiles from webservers into a single one

## Usage
```
Merges several logfiles into a big file

Usage:
  logmerge [command]

Available Commands:
  help        Help about any command
  merge       Merges the log files
  testLogfile Parses a logfile and outputs some statistics and not parsable log lines

Flags:
      --config string   config file (default is $HOME/.logmerge.yaml)
  -h, --help            help for logmerge

Use "logmerge [command] --help" for more information about a command.
```

```
Merges the log files

Usage:
  logmerge merge [flags]

Flags:
  -f, --file stringArray   Access files to merge (default [access.log])
  -h, --help               help for merge
  -o, --output string      Access files to write to (default "access-out.log")

Global Flags:
      --config string   config file (default is $HOME/.logmerge.yaml)
```
