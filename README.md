# Operator Log (opl) - Yet another operator logging tool

Simple go tool to log red team operations' activities. The idea is to avoid keeping manual track of timestamps when executing commands and instead just focus on executing. 

## Install & enable
Once Go is installed and configured, run:

```bash
❯❯❯ go install github.com/zkvL/opl/cmd@latest
```
If everything worked correctly, you should be able to run `opl -h` and see the help output.

### Fish 
You can create a fish shell function to automatically log shell issued commands (`opl -cmd <COMMAND>`):

```bash
# Fish shell
❯❯❯ function logCmd --on-event fish_prompt; set cmd $history[1]; opl -cmd "$cmd"; end
```

When you want to delete the function you can simply issue:
```bash
❯❯❯ functions -e logCmd
```

### zsh
If using zsh instead, you can add this to the `$HOME/.zshrc` file:

```bash
# $HOME/.zshrc
preexec() { opl -cmd "${1}" }
```

Then just source the file to start logging. When you are done simply remove the `preexec()` function and source again the configuration file.
You may need to restart the zsh shell.

## Use
`opl` will keep a registry of files within the `$HOME/operator-logs` folder. Each log file will be named with the current date in JSON format, as shown below:
```json
[
  {
    "date": "2023-08-05 17:26:09 GMT",
    "command": "amass enum -d DOMAIN.TARGET",
    "ipaddr": "XXX.XXX.XX.XXXX"
  },
  {
    "date": "2023-08-05 17:34:32 GMT",
    "command": "nmap --top-ports 1000 [...]",
    "ipaddr": "XXX.XXX.XX.XXXX",
    "operator": "zkvL"
  },
  [...]
]
```

If you want to log an activity, instead of a command, you can add it manually:
```bash
❯❯❯ opl 'Login to exposed Jenkins using the JenkinsAdmin account'
```

Without the `-cmd` flag, `opl` wont log a source IP address to the log.
```json
[...]
  {
    "date": "2023-08-05 20:16:05 GMT",
    "command": "Login to exposed Jenkins using the JenkinsAdmin account",
    "ipaddr": ""
  },
[...]
```

Finally, you can parse the logs to report activities using the `-print` flag


```bash
❯❯❯ opl -print $HOME/operator-logs
# OR
❯❯❯ opl -print $HOME/operator-logs/YYYY-MM-DD.json

Date                      IPAddr               Operator             Command/Activity             
---------------------------------------------------------------------------------------------------------------------------
2023-08-05 17:26:09 GMT   XXX.XXX.XX.XXXX                           amass enum -d DOMAIN.TARGET
2023-08-05 17:34:32 GMT   XXX.XXX.XX.XXXX      zkvL                 nmap --top-ports 1000 [...]
2023-08-05 20:16:05 GMT                                             Login to exposed Jenkins using the JenkinsAdmin account
[...]
```
The operator field will be added whenever the environment variable `OPERATOR` is set:
```bash
# Fish shell
set -g -x OPERATOR zkvL
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

## Credits
Got the idea of shell env functions to automatically log stuff from [c2biz](https://github.com/c2biz)