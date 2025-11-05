package main

import (
	"fmt"
	"os"
)

var (
	_args = os.Args
	prog  = _args[0]
	args  = _args[1:]
)

// wir könnten auch map["action"]actionFunc() machen
func main() {
	if len(args) == 0 {
		flagUsage()
		return
	}

	switch args[0] {
	case "run", "r":
		flagModeRun()

	case "get", "g":
		flagModeGetDataVol()

	case "nachbuchen", "buy":
		flagModeKaufDataVol()

	case "help", "h":
		flagUsage()

	case "version", "ver", "v":
		flagVersion()

	default:
		fmt.Printf("%s: unknown action\n", args[0])
		fmt.Printf("Run '%s help'\n", prog)
	}
}
