package lxc

import (
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
	"fmt"
	"time"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	TemplateConfig      `mapstructure:",squash"`
	ConfigFile          string            `mapstructure:"config_file"`
	OutputDir           string            `mapstructure:"output_directory"`
	ContainerName       string            `mapstructure:"container_name"`
	CommandWrapper      string            `mapstructure:"command_wrapper"`
	RawInitTimeout      string            `mapstructure:"init_timeout"`
	InitTimeout         time.Duration

	tpl *packer.ConfigTemplate
}

func NewConfig(raws ...interface{}) (*Config, error) {
	c := new(Config)
	md, err := common.DecodeConfig(c, raws...)
	if err != nil {
		return nil, err
	}

	c.tpl, err = packer.NewConfigTemplate()
	if err != nil {
		return nil, err
	}
	c.tpl.UserVars = c.PackerUserVars

	// Accumulate any errors
	errs := common.CheckUnusedConfig(md)
	errs = packer.MultiErrorAppend(errs, c.TemplateConfig.Prepare(c.tpl)...)

	if c.OutputDir == "" {
		c.OutputDir = fmt.Sprintf("output-%s", c.PackerBuildName)
	}

	if c.ContainerName == "" {
		c.ContainerName = fmt.Sprintf("packer-%s", c.PackerBuildName)
	}

	if c.CommandWrapper == "" {
		c.CommandWrapper = "{{.Command}}"
	}

	if c.RawInitTimeout == "" {
		c.RawInitTimeout = "20s"
	}

	c.InitTimeout, err = time.ParseDuration(c.RawInitTimeout)
	if err != nil {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("Failed parsing init_timeout: %s", err))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	return c, nil
}