package lxc

import (
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"log"
	"errors"
	"os"
	"path/filepath"
	"fmt"
)

// The unique ID for this builder
const BuilderId = "mitchellh.lxc"

type config struct {
	common.PackerConfig `mapstructure:",squash"`

	Distribution      string     `mapstructure:"distribution"`
	Release           string     `mapstructure:"release"`
	MirrorUrl         string     `mapstructure:"mirror_url"`
	SecurityMirrorUrl string     `mapstructure:"security_mirror_url"`
	OutputDir         string     `mapstructure:"output_directory"`
	ContainerName     string     `mapstructure:"container_name"`
	CommandWrapper    string     `mapstructure:"command_wrapper"`

	tpl *packer.ConfigTemplate
}

type wrappedCommandTemplate struct {
	Command string
}

type Builder struct {
	config config
	runner multistep.Runner
}


func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {
	md, err := common.DecodeConfig(&b.config, raws...)
	if err != nil {
		return nil, err
	}

	b.config.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, err
	}
	b.config.tpl.UserVars = b.config.PackerUserVars

	// Accumulate any errors
	errs := common.CheckUnusedConfig(md)

	if b.config.OutputDir == "" {
		b.config.OutputDir = fmt.Sprintf("output-%s", b.config.PackerBuildName)
	}

	if b.config.ContainerName == "" {
		b.config.ContainerName = fmt.Sprintf("packer-%s", b.config.PackerBuildName)
	}

	if b.config.CommandWrapper == "" {
		b.config.CommandWrapper = "{{.Command}}"
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	return nil, nil
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	wrappedCommand := func(command string) (string, error) {
		return b.config.tpl.Process(
			b.config.CommandWrapper, &wrappedCommandTemplate{
				Command: command,
			})
	}

	steps := []multistep.Step{
		new(stepPrepareOutputDir),
		new(stepLxcCreate),
		new(StepChrootProvision),
		new(stepExport),
	}

	// Setup the state bag
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("cache", cache)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("wrappedCommand", CommandWrapper(wrappedCommand))

	// Run
	if b.config.PackerDebug {
		b.runner = &multistep.DebugRunner{
			Steps:   steps,
			PauseFn: common.MultistepDebugFn(ui),
		}
	} else {
		b.runner = &multistep.BasicRunner{Steps: steps}
	}

	b.runner.Run(state)

	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	// Compile the artifact list
	files := make([]string, 0, 5)
	visit := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}

		return err
	}

	if err := filepath.Walk(b.config.OutputDir, visit); err != nil {
		return nil, err
	}

	artifact := &Artifact{
		dir: b.config.OutputDir,
		f:   files,
	}

	return artifact, nil
}

func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
