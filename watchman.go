package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var serverCmd *exec.Cmd

func startServer(serverFilePath string, port string) {
	port = "8081"
	pid, error := checkPort(port)
	if error != nil {
		log.Println(error)
	}
	if len(pid) != 0 {
		killProcess(pid)
	}

	serverCmd = exec.Command("go", "run", serverFilePath)
	serverCmd.Stdout = os.Stdout
	serverCmd.Stderr = os.Stderr

	err := serverCmd.Start()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func stopServer() {
	if serverCmd != nil && serverCmd.Process != nil {
		err := serverCmd.Process.Kill()
		if err != nil {
			log.Fatalf("Failed to stop server: %v", err)
		}
		serverCmd.Wait()
	}
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run watchman.go  <port> <server-filepath>")
		os.Exit(1)
	}

	log.Println(len(os.Args))
	serverFile := os.Args[2]
	port := os.Args[1]
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file:", event.Name)
					stopServer()
					startServer(serverFile, port)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()
	excludeDirs := []string{".git"}
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		for _, excludeDir := range excludeDirs {
			if strings.Contains(path, excludeDir) {
				return nil
			}
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	startServer(serverFile, port)

	<-done
}

// checkPort checks if a process is running on the specified port and returns its PID
func checkPort(port string) ([]int, error) {
	cmd := exec.Command("lsof", "-t", "-i", ":"+port)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return nil, nil
	}

	pidStrs := strings.Fields(string(output))
	pids := make([]int, len(pidStrs))
	for i, pidStr := range pidStrs {
		var pid int
		fmt.Sscanf(pidStr, "%d", &pid)
		pids[i] = pid
	}
	return pids, nil
}

// killProcess kills the process with the specified PID
func killProcess(pids []int) error {
	for _, pid := range pids {
		process, err := os.FindProcess(pid)
		if err != nil {
			return err
		}

		if err := process.Kill(); err != nil {
			return err
		}
	}
	return nil
}
