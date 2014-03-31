package lxc

import (
	"github.com/mitchellh/multistep"
	"fmt"
	"github.com/mitchellh/packer/packer"
	"bytes"
	"os/exec"
	"log"
	"strings"
)

type stepLxcCreate struct{}

func (s *stepLxcCreate) Run(state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*config)
	ui := state.Get("ui").(packer.Ui)

	name := config.ContainerName

	rootfs := fmt.Sprintf("/var/lib/lxc/%s/rootfs", name)

	command := []string{
		fmt.Sprintf("MIRROR=%s", config.MirrorUrl), "lxc-create",
			"-n", fmt.Sprintf("%s", name), "-t", config.Distribution, "--", "-r", config.Release,
	}

	ui.Say("Creating containter...")
	err := s.SudoCommand(command...)
	if err != nil {
		err := fmt.Errorf("Error creating container: %s", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	state.Put("mount_path", rootfs)

	return multistep.ActionContinue
}

func (s *stepLxcCreate) Cleanup(state multistep.StateBag) {
	config := state.Get("config").(*config)
	ui := state.Get("ui").(packer.Ui)

	command := []string{
		"lxc-destroy", "-f", "-n", config.ContainerName,
	}

	ui.Say("Unregistering and deleting virtual machine...")
	if err := s.SudoCommand(command...); err != nil {
		ui.Error(fmt.Sprintf("Error deleting virtual machine: %s", err))
	}
}


func (s *stepLxcCreate) SudoCommand(args ...string) error {
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