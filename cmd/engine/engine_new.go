package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/poteto0/poteto/cmd/template"
)

type EngineNewParam struct {
	ProjectName string
	IsFast      bool
	IsDocker    bool
	IsJSONRPC   bool
}

func RunNew(param EngineNewParam) error {
	prevDir, _ := filepath.Abs(".")
	defer os.Chdir(prevDir)

	fmt.Println("1. generating project: ", param.ProjectName)

	dirArr := strings.Split(param.ProjectName, "/")
	dirname := param.ProjectName
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
	if err := exec.Command("go", "mod", "init", param.ProjectName).Run(); err != nil {
		return err
	}

	fmt.Println("2. generating main.go")
	if err := createMain(param); err != nil {
		return err
	}

	fmt.Println("3. go mod tidy")
	if err := exec.Command("go", "mod", "tidy").Run(); err != nil {
		return err
	}

	if !param.IsDocker {
		return nil
	}

	fmt.Println("4. generating docker")
	if err := createDockerfile(); err != nil {
		return err
	}
	if err := createDockerCompose(); err != nil {
		return err
	}
	return nil
}

func createMain(param EngineNewParam) error {
	templateFile := choiceTemplateFile(param)

	return createAndWrite("main.go", templateFile)
}

func createDockerfile() error {
	return createAndWrite("Dockerfile", template.DockerTemplate)
}

func createDockerCompose() error {
	return createAndWrite("docker-compose.yaml", template.DockerComposeTemplate)
}

func createAndWrite(filename, templateFile string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	templateFileByte := []byte(templateFile)
	if _, err := f.Write(templateFileByte); err != nil {
		return err
	}

	return nil
}

func choiceTemplateFile(param EngineNewParam) string {
	if param.IsFast && !param.IsJSONRPC {
		return template.FastTemplate
	}

	if param.IsJSONRPC && !param.IsFast {
		return template.JSONRPCTemplate
	}

	if param.IsJSONRPC && param.IsFast {
		return template.JSONRPCFastTemplate
	}

	return template.DefaultTemplate
}