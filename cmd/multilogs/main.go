package main

import "multilogs/ui"

func main() {
	cmds := []string{
		"podman run --rm -it docker.io/nginx",
		"podman run --rm -it docker.io/metal3d/xmrig",
		"watch -n 1 date",
		"highlight -O ansi ui/ui.go",
	}
	ui.UI(cmds, 3)
}
