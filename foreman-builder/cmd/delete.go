package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	foremanbuilder "github.com/akshaykapse/foreman-builder-test/foreman-builder"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete container via name.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 1 {
			foremanbuilder.Logger.Errorf("Arg Error: Got %d args, when 1 was expected", len(args))
			fmt.Printf("Arg Error: Got %d args, when 1 was expected \n", len(args))
			os.Exit(1)
		}
		containerName := args[0]
		foremanUser.runDelete(fmt.Sprint(strings.Split(containerName, "::")[0]))

	},
}

func (u User) runDelete(containerName string) {

	// check if item actually exists first AND foreman-builder created it
	containersPath := u.containersPath
	containers, err := foremanbuilder.GetAllLines(containersPath, "::")
	foremanbuilder.Logger.Debugf("containers: %v\n", containers)
	foremanbuilder.Logger.Debugf("containerName: %s \n", containerName)

	if !slices.Contains(containers, containerName) {
		foremanbuilder.Logger.Fatalf("%s was not created with foreman-builder, foreman-builder will not delete it", containerName)
		os.Exit(1)

	}
	instanceInfo, err := foremanbuilder.GetInstanceInfo(containerName)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("instance not found: %s", containerName)) {
			foremanbuilder.Logger.Debugf("Instance does not exist in Lima anymore, just delete the line")
			foremanbuilder.DeleteLineInFile(containersPath, containerName)
			fmt.Printf("%s has been deleted\n", containerName)
			os.Exit(1)
		} else {
			foremanbuilder.Logger.Fatalf("Failed to get instance info: %v", err)
			os.Exit(1)
		}

	}
	if instanceInfo.Status == "Running" {
		var input string
		for input != "y" && input != "n" {
			fmt.Println("Lima instance is currently running! \n Do you want to stop the instance and delete? \n [y/n]")
			fmt.Scanln(&input)
		}
		if input == "n" {
			os.Exit(1)
		}
		// stop instance
		stopCmd := exec.Command("limactl", "stop", containerName)
		stopCmd.Stdout = os.Stdout
		stopCmd.Stderr = os.Stderr
		stopCmd.Run()

	}
	fmt.Printf("Deleting Lima instance...\n")
	deleteCmd := exec.Command("limactl", "delete", containerName, "-f")
	deleteCmd.Stdout = os.Stdout
	deleteCmd.Stderr = os.Stderr
	deleteCmd.Run()

	// delete container in dotfile
	foremanbuilder.DeleteLineInFile(containersPath, containerName)
	fmt.Printf("%s has been deleted\n", containerName)
}
