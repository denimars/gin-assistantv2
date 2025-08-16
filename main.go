package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-assistantv2/command"
	"github.com/gin-assistantv2/helper"
	"github.com/joho/godotenv"
)

const (
	usageInit    = "gin-assistant2 init"
	usageService = "gin-assistant2 service <service_name>"
	usageRun     = "gin-assistant2 run [port]"
	usageHelp    = "gin-assistant2 [init|service|run|help]"
)

var serverPort string

func setPort(args []string) error {
	godotenv.Load()

	if envPort := os.Getenv("PORT"); envPort != "" {
		serverPort = envPort
		return nil
	}

	if len(args) >= 2 {
		port := args[1]
		if num, err := strconv.Atoi(port); err == nil && num >= 7000 && num <= 9999 {
			serverPort = port
			return nil
		}
		return fmt.Errorf("port must be a number between 7000-9999")
	}
	serverPort = "8080"
	return nil
}

func isPortInUse(port string) bool {
	switch runtime.GOOS {
	case "windows":
		return isPortInUseWindows(port)
	default:
		return isPortInUseUnix(port)
	}
}

func isPortInUseWindows(port string) bool {
	cmd := exec.Command("cmd", "/C", "netstat -an | findstr :"+port)
	output, err := cmd.Output()
	return err == nil && len(output) > 0
}

func isPortInUseUnix(port string) bool {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("lsof -i:%s", port))
	err := cmd.Run()
	return err == nil
}

func handleInitCommand(projectDir string) {
	fmt.Println("Initializing project...")
	command.InitProject(projectDir)
	fmt.Println("Project initialized successfully!")
}

func handleServiceCommand(args []string, projectDir string) error {
	if len(args) < 2 {
		return fmt.Errorf("service name is required\nUsage: %s", usageService)
	}

	serviceName := args[1]
	fmt.Printf("Creating service: %s\n", serviceName)
	command.Service(projectDir, serviceName)
	fmt.Printf("Service '%s' created successfully!\n", serviceName)
	return nil
}

func handleRunCommand(args []string, projectDir string) error {
	if err := setPort(args); err != nil {
		return err
	}
	helper.Port(serverPort, projectDir)

	if isPortInUse(serverPort) {
		fmt.Printf("⚠️ Port %s is already in use, attempting to free it...\n", serverPort)
		freePort()
		time.Sleep(2 * time.Second)
		if isPortInUse(serverPort) {
			return fmt.Errorf("failed to free port %s", serverPort)
		}
	}

	fmt.Printf("🚀 Starting server on port %s...\n", serverPort)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	defer watcher.Close()

	appDir := helper.Path("./app")
	if err := watchFiles(watcher, appDir); err != nil {
		return fmt.Errorf("failed to watch files: %w", err)
	}

	fmt.Println("🚀 Starting server...")
	run()
	fmt.Println("🚀 Watching for file changes...")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			fmt.Println("🔄 File changed:", event.Name)
			if event.Op&fsnotify.Create != 0 {
				if fileInfo, err := os.Stat(event.Name); err == nil && fileInfo.IsDir() {
					fmt.Println("📂 New folder detected, adding to watcher:", event.Name)
					if err := watcher.Add(event.Name); err != nil {
						fmt.Println("❌ Error adding new folder to watcher:", err)
					}
				}
			}
			debounceRestart()
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Println("❌ File watcher error:", err)
		}
	}
}

func showHelp() {
	fmt.Println("Gin Assistant v2 - Go Gin Framework Project Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("  %s\n", usageInit)
	fmt.Printf("  %s\n", usageService)
	fmt.Printf("  %s\n", usageRun)
	fmt.Printf("  %s\n", "gin-assistant2 help")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init     Initialize a new Gin project in current directory")
	fmt.Println("  service  Generate a new service with repository, service and router files")
	fmt.Println("  run      Start development server with hot reload")
	fmt.Println("  help     Show this help message")
}

func getCurrentDirectory() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return dir, nil
}

var cmd *exec.Cmd
var mutex sync.Mutex
var debounceTimer *time.Timer

func stopServer() {
	if cmd != nil && cmd.Process != nil {
		fmt.Println("🛑 Stopping server...")
		if runtime.GOOS == "windows" {
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Println("❌ Error killing process:", err)
			}
		} else {
			err := cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				fmt.Println("⚠️ Error sending SIGTERM:", err)
			}
			_, _ = cmd.Process.Wait()
		}

		if processExists(cmd.Process.Pid) {
			fmt.Println("⚠️ Force killing process...")
			err := cmd.Process.Kill()
			if err != nil {
				fmt.Println("❌ Error force killing process:", err)
			}
		}
	}
	freePort()
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func freePort() {
	switch runtime.GOOS {
	case "windows":
		freePortWindows()
	default:
		freePortUnix()
	}
}

func freePortWindows() {
	cmd := exec.Command("cmd", "/C", "netstat -ano | findstr :"+serverPort)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("⚠️ Error checking port:", err)
		return
	}
	for _, line := range strings.Split(string(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 4 {
			pid := fields[len(fields)-1]
			fmt.Println("🔄 Killing process using port:", serverPort, "PID:", pid)
			killCmd := exec.Command("taskkill", "/F", "/PID", pid)
			killCmd.Run()
		}
	}
}

func freePortUnix() {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("lsof -t -i:%s | xargs kill -9", serverPort))
	cmd.Run()
}

func run() {
	mutex.Lock()
	defer mutex.Unlock()
	stopServer()
	time.Sleep(1 * time.Second)

	// Check port availability before starting
	if isPortInUse(serverPort) {
		fmt.Printf("⚠️ Port %s is in use, freeing it...\n", serverPort)
		freePort()
		time.Sleep(1 * time.Second)
	}

	cmd = exec.Command("go", "run", "main.go")
	cmd.Env = append(os.Environ(), "PORT="+serverPort)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
	fmt.Printf("🚀 Server restarted on port %s at %v\n", serverPort, time.Now())
}

func watchFiles(watcher *fsnotify.Watcher, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			fmt.Println("👀 Watching:", path)
			return watcher.Add(path)
		}
		return nil
	})
}

func debounceRestart() {
	if debounceTimer != nil {
		debounceTimer.Stop()
	}
	debounceTimer = time.AfterFunc(1*time.Second, run)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Printf("No command provided\nUsage: %s\n", usageHelp)
		os.Exit(1)
	}

	projectDir, err := getCurrentDirectory()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		handleInitCommand(projectDir)
	case "service":
		if err := handleServiceCommand(args, projectDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		showHelp()
	case "run":
		if err := handleRunCommand(args, projectDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\nUsage: %s\n", args[0], usageHelp)
		os.Exit(1)
	}
}
