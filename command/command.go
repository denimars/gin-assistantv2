package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-assistantv2/code"
	"github.com/gin-assistantv2/helper"
)

type DirectoryConfig struct {
	Name  string
	Files []FileConfig
}

type FileConfig struct {
	Name    string
	Content func() string
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func packageList() []string {
	return []string{
		"github.com/gin-gonic/gin",
		"github.com/gin-contrib/cors",
		"gorm.io/gorm",
		"gorm.io/driver/mysql",
		"github.com/joho/godotenv",
		"gorm.io/gorm/logger",
	}
}

func initProject(project string) {
	var output []byte
	var err error
	cmd := exec.Command("go", "mod", "init", project)
	if output, err = cmd.CombinedOutput(); err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}

func installPackage() {
	var output []byte
	var err error
	for _, p := range packageList() {
		cmd := exec.Command("go", "get", p)
		if output, err = cmd.CombinedOutput(); err != nil {
			panic(err)
		}
		fmt.Println(string(output))
	}
}

func getDirectoryConfig() []DirectoryConfig {
	return []DirectoryConfig{
		{
			Name: "db",
			Files: []FileConfig{
				{Name: "connection.go", Content: code.Connection},
			},
		},
		{
			Name: "helper",
			Files: []FileConfig{
				{Name: "helper.go", Content: code.Helper},
			},
		},
		{
			Name:  "model",
			Files: []FileConfig{},
		},
		{
			Name: "service",
			Files: []FileConfig{
				{Name: "validator.go", Content: code.Validator},
				{Name: "clonestruct.go", Content: code.CloneStruct},
				{Name: "base.go", Content: code.Base},
			},
		},
	}
}

func createDirectoriesAndFiles(appDir string, config []DirectoryConfig) {
	for _, dir := range config {
		dirPath := fmt.Sprintf("%s/%s", appDir, dir.Name)
		checkError(helper.CreateDirectory(dirPath))

		for _, file := range dir.Files {
			content := strings.TrimSpace(file.Content())
			checkError(helper.CreateFile(dirPath, file.Name, content))
		}
	}
}

func InitProject(dir string) {
	projectName := helper.GetProjectName(dir)
	appDir := fmt.Sprintf("%s/app", dir)

	if _, err := os.Stat(helper.Path(appDir)); os.IsNotExist(err) {
		initProject(projectName)
		installPackage()

		checkError(helper.CreateDirectory(helper.Path(appDir)))
		checkError(helper.CreateFile(appDir, "run.go", strings.TrimSpace(code.Run())))

		createDirectoriesAndFiles(appDir, getDirectoryConfig())

		checkError(helper.CreateFile(dir, "main.go", strings.TrimSpace(code.Main(projectName))))
	} else {
		fmt.Println("folder app exist...")
	}
}

func Service(dir string, serviceName string) {
	dirService := fmt.Sprintf("%s/app/service", dir)
	finalDir := fmt.Sprintf("%s/%s", dirService, serviceName)
	checkError(helper.CreateDirectory(helper.Path(finalDir)))
	checkError(helper.CreateFile(finalDir, "repository.go", strings.TrimSpace(code.Repository(serviceName))))
	checkError(helper.CreateFile(finalDir, "service.go", strings.TrimSpace(code.Service(serviceName))))
	checkError(helper.CreateFile(finalDir, "router.go", strings.TrimSpace(code.Router(serviceName))))

}
