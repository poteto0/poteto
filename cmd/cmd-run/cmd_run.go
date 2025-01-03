package cmdrun

import (
	"fmt"
	"os"

	"github.com/poteto-go/poteto/cmd/core"
	"github.com/poteto-go/poteto/cmd/engine"
	"github.com/poteto-go/poteto/utils"
)

func loadOption() core.RunnerOption {
	configFile, err := os.Open("./poteto.yaml")
	defer configFile.Close()

	if err != nil {
		utils.PotetoPrint("you can use poteto.yaml")
		return core.DefaultRunnerOption
	}

	configBytes := make([]byte, 1024)
	n, err := configFile.Read(configBytes)
	if err != nil || n == 0 {
		utils.PotetoPrint("warning error on reading poteto.yaml, use default option")
		return core.DefaultRunnerOption
	}

	var option core.RunnerOption
	err = utils.YamlParse(configBytes[:n], &option)
	if err != nil {
		utils.PotetoPrint("warning error on reading poteto.yaml, use default option")
		return core.DefaultRunnerOption
	}

	return option
}

func CommandRun() {
	option := loadOption()

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

	engine.RunRun(option)
}

func help() {
	fmt.Println("poteto-cli run: hot-reload run api server")
	fmt.Println("https://github.com/poteto-go/poteto")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help: Display help (this is this)")
}
