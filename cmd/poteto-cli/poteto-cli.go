package main

import (
	"fmt"
	"os"

	"github.com/poteto0/poteto/cmd/engine"
)

func main() {
	if len(os.Args) == 1 {
		help()
		os.Exit(-1)
	}

	if os.Args[1] == "new" {
		engine.New()
		os.Exit(-1)
	}

	for i := 1; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == "-h", os.Args[i] == "--help":
			help()
			os.Exit(-1)
		case os.Args[i] == "new":
			continue
		default:
			fmt.Println("unknown command or option:", os.Args[i])
			os.Exit(-1)
		}
	}
}

func help() {
	fmt.Println("poteto-cli: support creating poteto-app")
	fmt.Println("https://github.com/poteto0/poteto")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("Command: poteto-cli [command]")
	fmt.Println("  new:        create new poteto app")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help: Display help (this is this)")
}
