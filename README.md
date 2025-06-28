# ⚡ Thop - Tmux Hopper
Fast and lightweight interactive CLI for defining and jumping between projects / tmux sessions.

## About
Light and quick to use way of managing tmux sessions

### Why not tmuxinator?
Tmuxinator is a great tool, but I found it to be too troublesome and too complex for my own need.
Thop is designed to be lightweight, simple to install, and extremely quick to use.

### Features:
- Fast navigation to desired project / session from anywhere (including from inside of a Tmux session)
- Easy to edit yaml templates
- Execute shell commands in all/desired windows/panes

## Dependencies
- [fzf](https://github.com/junegunn/fzf)
- [tmux](https://github.com/tmux/tmux) 1.8+ (except for 2.5)

## Installation
Run below script in your terminal to install the latest release:

```bash
curl --proto '=https' --tlsv1.2 -sSf https://raw.githubusercontent.com/wezik/thop/develop/install.sh | sh
```

## Usage
```bash
thop [command]
```

### Commands:
```
create [name]          Creates a session template.
delete [name]          Deletes a session template.
edit [name]            Edits a session template.
help                   Shows help message.
kill [name]            Kills a session.
open [name]            Opens a session template.
```

`[name]` argument is always optional, if not provided thop will use defaults and (when needed) launch selector powered by fzf

### Editor

Thop uses your shell's default editor for opening files stored in `$EDITOR`
To change it your best option is to override it in your `.bashrc` or other shell config file:

```bashrc
export EDITOR='vim'
```

### Aliases

You can use aliases to make your life easier:

```
thop create:            thop c, thop new, thop add, thop a
thop delete:            thop d
thop edit:              thop e
thop kill:              thop k,
thop open:              thop o, thop select, thop s, thop
```

### Templates
Templates are blue-prints for your sessions, they are stored in `$XDG_CONFIG/thop/templates/`, edit such template using `thop edit` command

Example template:
```yaml
name: Example project name                  # (required) Name used for opening / selecting the project
version: 1                                  # Version of the template, used for migrations (leave it as is)
template:
  name: Optional session name               # Name of the session, will use project name if not present
  root: /home/foobar/projects/some_project  # Root directory for this session
  run:                                      # List of commands to be executed in all windows
  - echo 'Hello world'
  windows:                                  # List of windows to be created
  - name: main                              # Name of the window
    root: /optional/root/dir                # Root directory for this window
    layout: tiled                           # Layout for this window (tiled, main-vertical, main-horizontal, even-vertical, even-horizontal)
    run:                                    # List of commands to be executed in all panes
    - ls
    panes:                                  # List of panes to be created
    - run:
      - nvim                                # List of commands to be executed in this pane
    - root: /optional/root/dir/pane         # Root directory for this pane
      active: true                          # Set as active pane
  - name: logs
    active: true                            # Set as active window
    run:
    - mkdir ./tmp
    - touch ./tmp/logs.txt
    - tail -F ./tmp/logs.txt
```

All fields are optional unless stated as `(required)`

## Current state
This project follows [Semantic Versioning](https://semver.org/), but currently it's at version v0 as it's in development.
Destination is set but things can still change and break backwards compatibility; That includes templates they are not getting migrations until v1.

### TODO's:
- Integration tests
- Review the Makefile

### Ideas:
- A general config file
- Video showcase in README
- Opt in clearing of executed shell commands
- Setup and destroy blocks
- Make run parameter take in multi line string instead of array

### Post v1 ideas:
- Migrations for template versions
