# logmerge
Merges logfiles from webservers into a single one. This is very useful for HA servers that serve requests from two 
different machines for the same site. The files can be either plain text or gzip compressed.

At the moment it only supports the VCOMBINED and VCOMMON logfile format, but it should be easy to extend to 
support more formats.

The source log files need to have the log lines sorted by the timestamps.

The order in which the files are specified does not matter. 

The application reads a line from each input file and writes the oldest one. 
It then reads the next line from the file it just got the line from and the process repeats.
##Install

You can download the latest release by running:

    $ curl https://raw.githubusercontent.com/UnikumAB/logmerge/master/install.sh | sudo bash 
    
This will install the binary for your system to /usr/local/bin .

As always, you should have a look at the install.sh script to see what it is doing before you execute it on your system!

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
