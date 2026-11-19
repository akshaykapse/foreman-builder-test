package foremanbuilder

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type LimaInstance struct {
	Name      string `json:"name"`
	Status    string `json:"status"`
	SSHConfig string `json:"sshConfigFile"`
	Arch      string `json:"arch"`
	CPUs      int    `json:"cpus"`
	Memory    int64  `json:"memory"`
	Disk      int64  `json:"disk"`
	Dir       string `json:"dir"`
	Network   []struct {
		VNL      string `json:"vnl"`
		MACAddr  string `json:"macAddress"`
		IPAddr   string `json:"ipAddress"`
	} `json:"network"`
}

type ContainerInfoStruct struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type LimaOptions struct {
	ContainerName string `json:"containerName"`
	Username      string `json:"username"`
}

var execCommand = exec.Command

// IsLimaInstalled checks if limactl is available
func IsLimaInstalled() bool {
	cmd := execCommand("which", "limactl")
	return cmd.Run() == nil
}

// GetLimaInstances returns a list of all Lima instance names
func GetLimaInstances() ([]string, error) {
	cmd := execCommand("limactl", "list", "--format", "json")

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Lima outputs newline-delimited JSON, parse line by line
	var instances []LimaInstance
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var inst LimaInstance
		if err := json.Unmarshal([]byte(line), &inst); err != nil {
			return nil, fmt.Errorf("failed to parse Lima instance: %w", err)
		}
		instances = append(instances, inst)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	var names []string
	for _, inst := range instances {
		names = append(names, inst.Name)
	}

	return names, nil
}

// GetInstanceInfo returns detailed information about a specific Lima instance
func GetInstanceInfo(instanceName string) (*LimaInstance, error) {
	cmd := execCommand("limactl", "list", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("%s", exitErr.Stderr)
		}
		return nil, err
	}

	// Lima outputs newline-delimited JSON, parse line by line
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var inst LimaInstance
		if err := json.Unmarshal([]byte(line), &inst); err != nil {
			continue // Skip malformed lines
		}

		if inst.Name == instanceName {
			return &inst, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("instance not found: %s", instanceName)
}
