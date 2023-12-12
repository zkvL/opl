# Operator Log (opl) - Yet another operator logging tool

Simple go tool to log commands issued by red team operators. 
The idea is to avoid keeping manual track of timestamps when executing commands and instead just focus on executing. 

## Use
Just precede any command with this tool to get that logged. 
It will keep a registry of files within `$HOME/operator-logs` folder and each log is named with the current date in JSON format:

```bash
opl amass enum -d domain.com
```

To get a timeline of executed commands just use: 
```bash
opl -print $HOME/operator-logs
# OR
opl -print $HOME/operator-logs/YYYY-MM-DD.json
```

### Use in the background

You can use it with a little tweak to log every command without the need to precede with the opl command.

Idea from [c2biz](https://github.com/c2biz)

**ZSH**
```bash
# Add this line to your ~/.zshrc profile
preexec() { opl -runCmd=false "${1}" }
```

**Fish**
```bash
# Add this function to your ~/.config/fish/config.fish file
function log_cmd --on-event fish_prompt
  set cmd $history[1]
  opl -runCmd=false "$cmd"
 end
```

## Install

```bash
go install github.com/zkvL/opl/cmd@latest
```

## TODO
- [x] Install it as a service to capture commands in the background - Less complicated using shell env tweaks
- [ ] Manage complex commands, e.g. when using one command output to pipe into another

```bash
# Will log only the cat targets fragment
opl cat targets | httprobe

# Will log both commands as different log entries
opl cat targets | opl httprobe

# Will log the entire command and pass the execution to the shell environment
# Needs additional configuration, for instance using preexec() in zsh
# Otherwise, it will log the command but won't execute
preexec() { opl -runCmd=false "${1}" }
```
- [ ] Filter out common commands when using in `-runCmd=false` mode