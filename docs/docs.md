# pwsh-go

## Usage
> pwsh-go is a tool to update your powershell version automatically

pwsh-go

## Flags
|Flag|Usage|
|----|-----|
|`--debug`|enable debug messages|
|`--disable-update-checks`|disables update checks|
|`-i, --interactive`|Choose the PowerShell version interactively|
|`--raw`|print unstyled raw output (set it if output is written to a file)|

## Commands
|Command|Usage|
|-------|-----|
|`pwsh-go check`|Check if there's an update for your version|
|`pwsh-go completion`|Generate the autocompletion script for the specified shell|
|`pwsh-go help`|Help about any command|
|`pwsh-go update`|A brief description of your command|
# ... check
`pwsh-go check`

## Usage
> Check if there's an update for your version

pwsh-go check
# ... completion
`pwsh-go completion`

## Usage
> Generate the autocompletion script for the specified shell

pwsh-go completion

## Description

```
Generate the autocompletion script for pwsh-go for the specified shell.
See each sub-command's help for details on how to use the generated script.

```

## Commands
|Command|Usage|
|-------|-----|
|`pwsh-go completion bash`|Generate the autocompletion script for bash|
|`pwsh-go completion fish`|Generate the autocompletion script for fish|
|`pwsh-go completion powershell`|Generate the autocompletion script for powershell|
|`pwsh-go completion zsh`|Generate the autocompletion script for zsh|
# ... completion bash
`pwsh-go completion bash`

## Usage
> Generate the autocompletion script for bash

pwsh-go completion bash

## Description

```
Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(pwsh-go completion bash)

To load completions for every new session, execute once:

#### Linux:

	pwsh-go completion bash > /etc/bash_completion.d/pwsh-go

#### macOS:

	pwsh-go completion bash > /usr/local/etc/bash_completion.d/pwsh-go

You will need to start a new shell for this setup to take effect.

```

## Flags
|Flag|Usage|
|----|-----|
|`--no-descriptions`|disable completion descriptions|
# ... completion fish
`pwsh-go completion fish`

## Usage
> Generate the autocompletion script for fish

pwsh-go completion fish

## Description

```
Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	pwsh-go completion fish | source

To load completions for every new session, execute once:

	pwsh-go completion fish > ~/.config/fish/completions/pwsh-go.fish

You will need to start a new shell for this setup to take effect.

```

## Flags
|Flag|Usage|
|----|-----|
|`--no-descriptions`|disable completion descriptions|
# ... completion powershell
`pwsh-go completion powershell`

## Usage
> Generate the autocompletion script for powershell

pwsh-go completion powershell

## Description

```
Generate the autocompletion script for powershell.

To load completions in your current shell session:

	pwsh-go completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.

```

## Flags
|Flag|Usage|
|----|-----|
|`--no-descriptions`|disable completion descriptions|
# ... completion zsh
`pwsh-go completion zsh`

## Usage
> Generate the autocompletion script for zsh

pwsh-go completion zsh

## Description

```
Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions for every new session, execute once:

#### Linux:

	pwsh-go completion zsh > "${fpath[1]}/_pwsh-go"

#### macOS:

	pwsh-go completion zsh > /usr/local/share/zsh/site-functions/_pwsh-go

You will need to start a new shell for this setup to take effect.

```

## Flags
|Flag|Usage|
|----|-----|
|`--no-descriptions`|disable completion descriptions|
# ... help
`pwsh-go help`

## Usage
> Help about any command

pwsh-go help [command]

## Description

```
Help provides help for any command in the application.
Simply type pwsh-go help [path to command] for full details.
```
# ... update
`pwsh-go update`

## Usage
> A brief description of your command

pwsh-go update

## Description

```
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.
```


---
> **Documentation automatically generated with [PTerm](https://github.com/pterm/cli-template) on 20 July 2023**
