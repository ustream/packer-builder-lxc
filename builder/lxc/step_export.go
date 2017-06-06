package lxc

import (
	"github.com/mitchellh/multistep"
	"fmt"
	"github.com/mitchellh/packer/packer"
	"bytes"
	"os/exec"
	"log"
	"strings"
	"path/filepath"
	"os"
	"io"
	"encoding/json"
)

type stepExport struct{}

type Metadata struct {
	Provider string `json:"provider"`
	Version string  `json:"version"`
}

func (s *stepExport) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	name := config.ContainerName

	containerDir := fmt.Sprintf("/var/lib/lxc/%s", name)
	outputPath := filepath.Join(config.OutputDir, "rootfs.tar.gz")
	configFilePath := filepath.Join(config.OutputDir, "lxc-config")
	metadataFilePath := filepath.Join(config.OutputDir, "metadata.json")

	metadata := Metadata{"lxc", "1.0.0"}
	metadataFile, err := os.Create(metadataFilePath)

	if err != nil {
		err := fmt.Errorf("Error creating metadata file: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		err := fmt.Errorf("Error marshaling metadata : %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	_, err = metadataFile.Write(metadataJson)
	metadataFile.Sync()

	if err != nil {
		err := fmt.Errorf("Error writing metadata file : %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	configFile, err := os.Create(configFilePath)

	if err != nil {
		err := fmt.Errorf("Error creating config file: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	originalConfigFile, err := os.Open(config.ConfigFile)

	if err != nil {
		err := fmt.Errorf("Error opening config file: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	_, err = io.Copy(configFile, originalConfigFile)

	commands := make([][]string, 4)
	commands[0] = []string{
		"lxc-stop", "--name", name,
	}
	commands[1] = []string{
		"tar", "-C", containerDir, "--numeric-owner", "--anchored", "--exclude=./rootfs/dev/log", "-czf", outputPath, "./rootfs",
	}
	commands[2] = []string{
		"chmod", "+x", configFilePath,
	}
	commands[3] = []string{
		"sh", "-c", "chown $USER:`id -gn` " + filepath.Join(config.OutputDir, "*"),
	}

	ui.Say("Exporting container...")
	for _, command := range commands {
		err := s.SudoCommand(command...)
		if err != nil {
			err := fmt.Errorf("Error exporting container: %s", err)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		}
	}

	return multistep.ActionContinue
}

func (s *stepExport) Cleanup(state multistep.StateBag) {}


func (s *stepExport) SudoCommand(args ...string) error {
	var stdout, stderr bytes.Buffer

	log.Printf("Executing sudo command: %#v", args)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		err = fmt.Errorf("Sudo command error: %s", stderrString)
	}

	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)

	return err
}
