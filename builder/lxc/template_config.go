package lxc

import (
	"fmt"
	"github.com/mitchellh/packer/packer"
)

type TemplateConfig struct {
	Name           string   `mapstructure:"template_name"`
	Parameters     []string `mapstructure:"template_parameters"`
	EnvVars        []string `mapstructure:"template_environment_vars"`
	TargetRunlevel int      `mapstructure:"target_runlevel"`
}

func (c *TemplateConfig) Prepare(t *packer.ConfigTemplate) []error {
	errs := make([]error, 0)

	if c.Parameters == nil {
		c.Parameters = make([]string, 0)
	}
	for i, arg := range c.Parameters {
		if err := t.Validate(arg); err != nil {
			errs = append(errs, fmt.Errorf("Error processing parameters[%d]: %s", i, err))
		}
	}

	if c.EnvVars == nil {
		c.EnvVars = make([]string, 0)
	}
	for i, arg := range c.EnvVars {
		if err := t.Validate(arg); err != nil {
			errs = append(errs, fmt.Errorf("Error processing environment_vars[%d]: %s", i, err))
		}
	}

	return errs
}
