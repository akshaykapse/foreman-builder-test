package main

import (
	"os"
	"path/filepath"
	"runtime"

	foremanbuilder "github.com/akshaykapse/foreman-builder-test/foreman-builder"
	"github.com/akshaykapse/foreman-builder-test/foreman-builder/cmd"
)

func main() {

	if !isSystemSupported() {
		os.Exit(1)
	}

	if !foremanbuilder.IsLimaInstalled() {
		foremanbuilder.Logger.Fatal("Lima must be installed. Install with: brew install lima")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		foremanbuilder.Logger.Fatalf("Failed to get home directory: %v", err)
	}

	dotFolderPath := filepath.Join(home, ".foreman-builder")
	containersPath := filepath.Join(dotFolderPath, "containers")

	dotfolderExists, err := foremanbuilder.DoesFileOrDirectoryExist(dotFolderPath)
	if err != nil {
		foremanbuilder.Logger.Errorf("Error checking dotfolder: %v", err)
	}

	containersExists, err := foremanbuilder.DoesFileOrDirectoryExist(containersPath)
	if err != nil {
		foremanbuilder.Logger.Errorf("Error checking containers file: %v", err)
	}

	if !dotfolderExists {
		foremanbuilder.Logger.Info("Creating dotfolder...")
		err := os.MkdirAll(dotFolderPath, 0755)
		if err != nil {
			foremanbuilder.Logger.Fatalf("Failed to create dotfolder: %v", err)
		}
	}

	if !containersExists {
		foremanbuilder.Logger.Info("Creating containers file...")
		err := os.WriteFile(containersPath, []byte(""), 0644)
		if err != nil {
			foremanbuilder.Logger.Fatalf("Failed to create container file: %v", err)
		}
	}
	cmd.SyncContainers()

	cmd.Execute()
}

func isSystemSupported() bool {
	// Lima supports both arm64 (Apple Silicon) and amd64 (Intel)
	userCPU := runtime.GOARCH

	if userCPU != "arm64" && userCPU != "amd64" {
		foremanbuilder.Logger.Fatalf("Unsupported architecture: %s. Lima supports arm64 and amd64.", userCPU)
		return false
	}

	foremanbuilder.Logger.Infof("Detected architecture: %s", userCPU)
	return true
}
