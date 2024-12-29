package cmdnew

import (
	"fmt"
	"os"

	"github.com/poteto0/poteto/cmd/engine"

	"github.com/manifoldco/promptui"
)

func CommandNew() {
	param := engine.EngineNewParam{}

	fmt.Println("You can also use poteto-cli -h | --help")
	for i := 2; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == "-h", os.Args[i] == "--help":
			help()
			os.Exit(-1)
		case os.Args[i] == "-f", os.Args[i] == "--fast":
			param.IsFast = true
		case os.Args[i] == "-d", os.Args[i] == "--docker":
			param.IsDocker = true
		case os.Args[i] == "-j", os.Args[i] == "--jsonrpc":
			param.IsJSONRPC = true
		default:
			fmt.Println("unknown command or option:", os.Args[i])
			os.Exit(-1)
		}
	}

	wd, _ := os.Getwd()
	fmt.Println("Generate New Poteto App @", wd)

	prompt := promptui.Prompt{
		Label: "your project [github.com/github/poteto-api]", // 表示する文言
	}
	projectName, _ := prompt.Run()
	if len(projectName) == 0 {
		projectName = "github.com/github/poteto-api"
	}

	param.ProjectName = projectName

	err := engine.RunNew(param)
	if err != nil {
		panic(err)
	}
}

func help() {
	fmt.Println("poteto-cli new: support creating poteto-app")
	fmt.Println("https://github.com/poteto0/poteto")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help: Display help (this is this)")
	fmt.Println("  -f, --fast: fast mode api (doesn't gen requestId automatic)")
	fmt.Println("  -d, --docker: with Dockerfile & docker-compose w golang@1.23")
	fmt.Println("  -j, --jsonrpc: jsonrpc template")
}
