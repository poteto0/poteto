package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/poteto0/poteto/cmd/template"

	"github.com/manifoldco/promptui"
)

var isFast = false

func CommandNew() {
	fmt.Println("You can also use poteto-cli -h | --help")
	for i := 2; i < len(os.Args); i++ {
		switch {
		case os.Args[i] == "-h", os.Args[i] == "--help":
			help()
			os.Exit(-1)
		case os.Args[i] == "-f", os.Args[i] == "--fast":
			isFast = true
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

	err := run(projectName)
	if err != nil {
		panic(err)
	}
}

func run(projectName string) error {
	prevDir, _ := filepath.Abs(".")
	defer os.Chdir(prevDir)

	fmt.Println("1. generating project: ", projectName)

	dirArr := strings.Split(projectName, "/")
	dirname := projectName
	if len(dirArr) > 1 {
		dirname = dirArr[(len(dirArr) - 1)]
	}
	if err := os.Mkdir(dirname, 0755); err != nil {
		if !os.IsExist(err) {
			// 「ディレクトリが既に存在する」以外のエラー
			return err
		}
	}
	if err := os.Chdir(dirname); err != nil {
		return err
	}
	if err := exec.Command("go", "mod", "init", projectName).Run(); err != nil {
		return err
	}

	fmt.Println("2. generating main.go")
	if err := createMain(); err != nil {
		return err
	}

	fmt.Println("3. go mod tidy")
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return err
	}
	return nil
}

func createMain() error {
	f, err := os.Create("main.go")
	if err != nil {
		return err
	}
	defer f.Close()

	templateFile := []byte(template.DefaultTemplate)
	if isFast {
		templateFile = []byte(template.FastTemplate)
	}

	mainGoByte := templateFile
	if _, err := f.Write(mainGoByte); err != nil {
		return err
	}

	return nil
}

func help() {
	fmt.Println("poteto-cli new: support creating poteto-app")
	fmt.Println("https://github.com/poteto0/poteto")
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help: Display help (this is this)")
	fmt.Println("  -f, --fast: fast mode api (doesn't gen requestId automatic)")
}
