package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	foremanbuilder "github.com/akshaykapse/foreman-builder-test/foreman-builder"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new container environment",
	Run: func(cmd *cobra.Command, args []string) {
		foremanUser.runCreate()
	},
}

var containerType = "lima"

func (u User) runCreate() {
	currentUser, err := user.Current()
	if err != nil {
		foremanbuilder.Logger.Fatalf("Failed to get current user: %s", err)
	}
	username := currentUser.Username
	fmt.Println("Username:", username)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Container name (default foreman): ")
	containerName, _ := reader.ReadString('\n')
	containerName = strings.TrimSpace(containerName)
	if containerName == "" {
		containerName = "foreman"
	}

	containerNameExists, err := foremanbuilder.GetLineInFile(u.containersPath, containerName, "")
	if err != nil {
		if err.Error() != "not_found" {
			fmt.Println("An Error has occured searching dotfile")
			os.Exit(1)
		}
	}
	foremanbuilder.Logger.Debug("container name: %s \n", containerName)
	foremanbuilder.Logger.Debug("container name EXISTS: %s \n", containerNameExists)
	if containerNameExists != "" {
		fmt.Println("Container name must be unique")
		os.Exit(1)
	}

	err = foremanbuilder.AppendToFile(u.containersPath, fmt.Sprintf("%s::%s", containerName, containerType))
	if err != nil {
		foremanbuilder.Logger.Error("Failed to write container to container file")
	}

	foremanbuilder.Logger.Info("Starting environment creation")

	if containerType == "lima" {
		limaOpts := foremanbuilder.LimaOptions{
			Username:      username,
			ContainerName: containerName,
		}
		err := foremanUser.createLimaInstance(limaOpts)
		if err != nil {
			// better error message to show?
			fmt.Println("An error has occured during Lima instance creation")
			os.Exit(1)
		}
	}

}

func (u User) createLimaInstance(opts foremanbuilder.LimaOptions) error {
	config, err := foremanbuilder.GetYmlValues("./config.yml")
	if err != nil {
		foremanbuilder.Logger.Info("No config file found, skipping")
	}

	data := foremanbuilder.LimaConfigData{
		Username: opts.Username,
		Packages: config.Packages,
	}

	confsDir := filepath.Join(u.dotFilePath, "confs")
	if err := os.MkdirAll(confsDir, 0755); err != nil {
		foremanbuilder.Logger.Errorf("Failed to create confs directory: %v", err)
		return err
	}

	pathName := filepath.Join(confsDir, fmt.Sprintf("lima-foreman-%s.yml", data.Username))
	foremanbuilder.Logger.Info("using", pathName)
	err = foremanbuilder.GenerateContainerConfig(data, pathName)
	if err != nil {
		foremanbuilder.Logger.Errorf("failed to generate Lima config, err : %v\n", err)
		return err
	}

	// Create Lima instance
	foremanbuilder.Logger.Info("Creating Lima instance:", opts.ContainerName)
	createCmd := exec.Command("limactl", "create", "--name", opts.ContainerName, pathName)
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr
	if err := createCmd.Run(); err != nil {
		foremanbuilder.Logger.Errorf("Error creating Lima instance: %s", err)
		return err
	}

	// Start Lima instance
	foremanbuilder.Logger.Info("Starting Lima instance:", opts.ContainerName)
	startCmd := exec.Command("limactl", "start", opts.ContainerName)
	startCmd.Stdout = os.Stdout
	startCmd.Stderr = os.Stderr
	if err := startCmd.Run(); err != nil {
		foremanbuilder.Logger.Errorf("Error starting Lima instance: %s", err)
		return err
	}

	fmt.Println("✓ Lima instance has been created and started!")
	fmt.Printf("\nTo access the instance:\n")
	fmt.Printf("  limactl shell %s\n\n", opts.ContainerName)
	fmt.Printf("Once inside, start Foreman:\n")
	fmt.Printf("  cd ~/foreman && bundle exec foreman start\n\n")
	fmt.Printf("Access Foreman at: https://foreman-dev.example.com\n")
	fmt.Printf("Credentials: admin / changeme\n")

	return nil

}
