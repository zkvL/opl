# Operator Log (opl) - Yet another operator logging tool

Simple go tool to log activities issued by red team operators. 
The idea is to avoid keeping manual track of timestamps when executing commands and instead just focus on executing. 

## Install & enable
You can simply download the binary and place it within your `$PATH`. The following command will enable logging all text input to the terminal:

```bash
opl -enable [fish|zsh]
source ~/.config/fish/config.fish
# OR
source  ~/.zshrc
```

Effectively it will change the ~/.zshrc or ~/.config/fish/config.fish configuration to add the following:

```bash
# Fish sehll
function logCmd --on-event fish_prompt
  set cmd $history[1]
  opl "$cmd"
end
```

```bash
# ZSH shell
preexec() { opl "${1}" }
```

Note that currently, it only changes `zsh` or `fish` shell configuration. 

## Use
When enabled, `opl` will keep a registry of files within the `$HOME/operator-logs` folder. Each log file will be named with the current date in JSON format, as shown below:
```bash
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
opl [-noip] 'Login to exposed Jenkins using the JenkinsAdmin account'
```

The `-noip` flag will prevent the command from adding an IP address to the log (in case you executed such activity from another endpoint than the one registering the activity).
Finally, you can parse the logs to report activities using the `-print` flag


```bash
opl -print $HOME/operator-logs
# OR
opl -print $HOME/operator-logs/YYYY-MM-DD.json

Date                      IPAddr               Operator             Command             
-----------------------------------------------------------------------------------------------
2023-08-05 17:26:09 GMT   XXX.XXX.XX.XXXX                           amass enum -d DOMAIN.TARGET
2023-08-05 17:34:32 GMT   XXX.XXX.XX.XXXX      zkvL                 nmap --top-ports 1000 [...]
[...]
```
The operator field will be added whenever the environment variable `OPERATOR` is set:
```bash
# Fish shell
set -g -x OPERATOR zkvL
```

## Disable
`opl -disable [fish|zsh]` will disable logging every input. Needs to open a new shell.

## TODO
- [x] Use shell env tweaks to capture commands in the background. Idea from [c2biz](https://github.com/c2biz)
- [ ] Filter out common commands like `ls`, `cd`, `mkdir`, `opl`` itself, etc.