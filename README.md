# Operator Log (opl) - Yet another operator logging tool

Go utility to automate red team activity logging, allowing operators to focus on task execution without manual timestamp tracking.

## Install

You can just download the binary from [releases](https://github.com/zkvL/opl/releases) and place it into your `$PATH` (e.g. in `/usr/local/bin`).

Or, if you're a Go user, use `go install`
```shell
go install github.com/zkvL/opl@latest
```

## Usage 
`opl` will keep a registry of files within the `$HOME/.oplogs` folder. Each log file will be named with the current date in JSON format, as shown below:
```json
[
  {
    "date": "2025-07-30 00:19:41 UTC",
    "activity": "Sent email phishing campaing",
    "ipaddr": [
      "187.130.216.101"
    ],
    "operator": ""
  },
  {
    "date": "2025-07-30 00:19:53 UTC",
    "activity": "Login to exposed Jenkins using the JenkinsAdmin account",
    "ipaddr": [
      "189.100.12.11"
    ],
    "operator": "zkvL"
  },
  {
    "date": "2025-07-30 00:20:05 UTC",
    "activity": "Run distributed amass scan",
    "ipaddr": [
      "189.100.12.11",
      "189.100.12.12",
      "189.100.12.13"
    ],
    "operator": ""
  },
  {
    "date": "2025-07-30 00:20:21 UTC",
    "activity": "nmap -p- -T4 -sV -sC -oA output -vv -Pn 123.123.123.4,192.168.1.1",
    "ipaddr": [
      "187.130.216.91"
    ],
    "operator": "zkvL"
  }
]
```

> **NOTE:** The operator field will be added whenever the environment variable `OPERATOR` is set.

### Log activities
- If you need to log an activity, not necessarly an executed command you can just:
```shell
opl -a "Sent email phishing campaing"
```

- Maybe you need to log where this activity came from:
```shell
opl -a "Login to exposed Jenkins using the JenkinsAdmin account" -i 189.100.12.11
opl -a "Run distributed amass scan" -i 189.100.12.11,189.100.12.12,189.100.12.13
```

By default, if no `-i` option is specified `opl` will log the public IP address from where the utility is executed. Which is useful for the next use case.

### Log commands
To ensure that `opl` solely focused on logging, command execution and logging can be automated using shell environment functions - cause we don't have time to execute `opl -a` for every activity we do right?:

#### fish shell
- You can create a fish shell function to automatically log shell issued commands (`opl -act <COMMAND>`):

```shell
set -g -x OPERATOR zkvL
function logCmd --on-event fish_prompt; set cmd $history[1]; opl -act "$cmd"; end
```

- When you want to delete the function you can simply issue:
```shell
functions -e logCmd
```

#### zsh shell
- If using zsh, you can add this to the `$HOME/.zshrc` file:

```bash
# $HOME/.zshrc
preexec() { opl -act "${1}" }
```

Then just source the file to start logging. When you are done simply remove the `preexec()` function and source again the configuration file.
You may need to restart the zsh shell.

### Reporting
To parse the logs you can use the `-show` flag:

```shell
# Parse logs from default location $HOME/.oplogs or from a given folder or file
opl -s [-l date-project/operator/logs/YYYY-MM-DD.json]
Operator Operator IP(s)                              Timestamp (UTC)         Command/Activity                                                 
-------------------------------------------------------------------------------------------------------------------------------------------
         187.130.216.101                             2025-07-30 00:19:41 UTC Sent email phishing campaing                                     
zkvL     189.100.12.11                               2025-07-30 00:19:53 UTC Login to exposed Jenkins using the JenkinsAdmin account          
         189.100.12.11, 189.100.12.12, 189.100.12.13 2025-07-30 00:20:05 UTC Run distributed amass scan                                       
zkvL     187.130.216.91                              2025-07-30 00:20:21 UTC nmap -p- -T4 -sV -sC -oA output -vv -Pn 123.123.123.4,192.168.1.1
[...]
```

Currently `opl` also supports formatting the output to markdown:
```shell
opl -s -f md
| Operator | Operator IP(s)                              | Timestamp (UTC)         | Command/Activity                                                  |
|----------|---------------------------------------------|-------------------------|-------------------------------------------------------------------|
|          | 187.130.216.101                             | 2025-07-30 00:19:41 UTC | Sent email phishing campaing                                      |
| zkvL     | 189.100.12.11                               | 2025-07-30 00:19:53 UTC | Login to exposed Jenkins using the JenkinsAdmin account           |
|          | 189.100.12.11, 189.100.12.12, 189.100.12.13 | 2025-07-30 00:20:05 UTC | Run distributed amass scan                                        |
| zkvL     | 187.130.216.91                              | 2025-07-30 00:20:21 UTC | nmap -p- -T4 -sV -sC -oA output -vv -Pn 123.123.123.4,192.168.1.1 |
[...]
```

Or XLSX:
```shell
opl -s -f xlsx
[+] Logs written successfully to file: opl-timeline.xlsx
```

## Filtering

`opl` filters out the following common commands from logging:
- alias
- cd
- chmod
- chown
- cp
- exit
- find
- id
- kill
- ls
- locate
- make
- man
- mkdir
- mv
- nano
- opl
- ps
- pwd
- uname
- vim
- which
- whoami

## Get Help
```shell
opl -h
Yet another operator logging tool for Red Teamers

Usage:
  opl [flags]

Flags:
  -a, --act string        Logs an activity
  -i, --ip string         Comma-separated IPs to log from where the activity was performed. If not set, public IP will be logged
  -s, --show              Print logs from the default location or any specified with -l
  -l, --location string   Individual file or folder name to print logs from (default "/Users/ybasurto/.oplogs")
  -f, --format string     Print logs in specified format (currently supported md and xlsx)
  -d, --debug             Print debug error messages
  -h, --help              help for opl
```

## Credits
Got the idea of shell env functions to automatically log stuff from [c2biz](https://github.com/c2biz)