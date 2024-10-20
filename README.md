# PifPaf log

Simple command to log several commands at once in a tiled view in the terminal.

Stop all the command simultaneously with `Ctrl+C`.

Easy, useful.

## Usage

Give the commands as argument (use quotes if needed to identify each command).
And leave the rest to `pifpaf`.

```bash
pifpaf "command1" "command2" "command3"

# Less columns ?
pifpaf -c 2 "command1" "command2" "command3"
```

Examples:

```bash
pifpiaf \
  "podman run --rm -it nginx" \
  "podman run --rm -it metal3d/xmrig" \
  "dmesg -w"
```

![3 columns by default](assets/pifpaf1.png)

Using `-c 2` to have only 2 columns:

```bash
pifpaf -c 2 \
  "podman run --rm -it nginx" \
  "podman run --rm -it metal3d/xmrig" \
  "dmesg -w"
```

![With 2 columns max](assets/pifpaf2.png)

## Installation

Get the release from the [releases page](https://github.com/metal3d/pifpaf/releases) and put it in your `PATH`.

You may, also, use `go inteall` if you have Go installed:

```bash
go install github.com/metal3d/pifpaf@latest
```

## Known issues

> Sometimes, pressing `CTRL+C` does not stop the commands. The terminal doesn't give back the prompt.

It's a known issue, and I'm working on it. But no panic, only press `CTRL+C` again, and it will stop the commands.

> Colors are sometimes not well displayed

It's a limitation of the `tview` library. I'm working on it.

> Some commands need to refresh the screen to display the output, like in `watch`, or for example `top` or `htopt`...

At this time, PifPaf does not support this kind of command. It's a "log" tool, not a "monitor" tool. But I try to find a
way to support this kind of command.

## Thanks to

- Thanks to Rivo for [TView](https://github.com/rivo/tview) library, it's a great tool to build terminal applications.
- And thanks to the [Cobra](https://github.com/spf13/cobra) library for the CLI part.
