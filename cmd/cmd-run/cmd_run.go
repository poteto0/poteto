package cmdrun

import (
	"fmt"
	"os"

	"github.com/poteto0/poteto/cmd/engine"
)

// TODO: setting yaml file
func CommandRun() {
	param := engine.EngineRunParam{}

	fmt.Println("You can also use poteto-cli run -h | --help")
	for i := 2; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == "-h", os.Args[i] == "--help":
			help()
			os.Exit(-1)
		default:
			fmt.Println("unknown command or option:", os.Args[i])
			os.Exit(-1)
		}
	}

	engine.RunRun(param)
}

func help() {
	fmt.Println("poteto-cli run: hot-reload run api server")
	fmt.Println("https://github.com/poteto0/poteto")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help: Display help (this is this)")
}
