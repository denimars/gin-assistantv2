package helper

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func ClearPackageName(serviceName string) string {
	serviceNameSplit := strings.Split(serviceName, "/")
	return strings.NewReplacer(".", "", "/", "", "\\", "", "-", "").Replace(strings.ToLower(serviceNameSplit[len(serviceNameSplit)-1]))
}

func getEndSliceDirectory(directory string) string {
	var split_ []string
	switch runtime.GOOS {
	case "windows":
		split_ = strings.Split(directory, "\\")
	default:
		split_ = strings.Split(directory, "/")
	}

	return split_[len(split_)-1]
}

func GetServiceName(directory string) string {
	return getEndSliceDirectory(directory)
}

func GetProjectName(directory string) string {
	return getEndSliceDirectory(directory)
}

func Path(path string) string {
	switch runtime.GOOS {
	case "windows":
		path_ := strings.ReplaceAll(path, `\`, `\\`)
		path__ := strings.ReplaceAll(path_, `/`, `\\`)
		return path__

	default:
		return path
	}
}

func CreateFile(path, fileName, text string) error {
	path_ := fmt.Sprintf("%s/%s", path, fileName)
	var err error
	if file, err := os.Create(Path(path_)); err == nil {
		defer file.Close()
		if _, err = file.Write([]byte(text)); err == nil {
			fmt.Println("assistant created file: " + fileName)
		}
	}
	return err

}
func sliceDirectory(path string) []string {
	switch runtime.GOOS {
	case "windows":
		return strings.Split(path, `\\`)
	default:
		return strings.Split(path, "/")
	}
}

func CreateDirectory(path string) error {
	pathSlice := sliceDirectory(path)
	var err error
	if err = os.Mkdir(path, 0755); err == nil {
		fmt.Println("assistant created directory: " + pathSlice[len(pathSlice)-1])
	}
	return err
}

func reWriteFile(dir string, content []string) error {
	connection, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer connection.Close()
	for _, line := range content {
		_, err := connection.WriteString(line + "\n")
		if err != nil {

			return err
		}
	}
	return nil
}

func readFile(searchFile string, fixedFile string, isALl bool, directory string) []string {
	var result []string
	file, err := os.Open(directory)
	if err != nil {
		fmt.Println("file not open", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchFile) {
			if isALl {
				result = append(result, fixedFile)
			} else {
				modifiedLine := strings.Replace(line, searchFile, fixedFile, 1)
				result = append(result, modifiedLine)
			}

		} else {
			result = append(result, line)
		}

	}
	return result
}

func Port(port string, directory string) {
	dir := Path(fmt.Sprintf("%v/app/run.go", directory))
	newPort := fmt.Sprintf("port := \"%v\"", port)
	addPort := readFile("port := \"", newPort, true, dir)
	reWriteFile(dir, addPort)

}
