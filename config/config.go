package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

// Mode for trumpet
var Mode = "debug"

// Version for trumpet
var Version = "0.0.0"

// BuildHash from vcs
var BuildHash = ""

// AppName for trumpet
var AppName = "trumpet"

// IConfiguration : configuration interface
type IConfiguration interface {
	Init()
}

// ConfigurationType : configuration type
type ConfigurationType struct {
	Debug             bool
	ConfigurationPath string `yaml:"configuration_path" validate:"nonzero"`

	StrictInclude bool     `yaml:"strict_include"`
	Includes      []string `yaml:"includes,omitempty"`

	Logging Logging `yaml:"logging"`
	HTTP    HTTP    `yaml:"http"`
}

// Init : init ConfigurationType
func (c *ConfigurationType) Init() {
	c.Debug = Mode == "debug"
	c.ConfigurationPath = fmt.Sprintf("%s.yaml", AppName)
	c.StrictInclude = false

	c.Logging.Init()
	c.HTTP.Init()
}

// ReadFrom : read configuration from path
func (c *ConfigurationType) ReadFrom(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	return nil
}

// Read : read configuration
func (c *ConfigurationType) Read() error {
	var err error
	confPath := c.ConfigurationPath

	err = c.ReadFrom(confPath)
	if err != nil {
		return err
	}

	dirPath, _ := filepath.Split(confPath)
	includes := c.Includes
	strictInclude := c.StrictInclude

	for _, p := range includes {
		if !filepath.IsAbs(p) {
			p, err = filepath.Abs(filepath.Join(dirPath, p))
			if strictInclude && err != nil {
				return err
			}
		}
		err = c.ReadFrom(p)
		if strictInclude && err != nil {
			return err
		}
	}
	c.Includes = includes
	c.StrictInclude = strictInclude
	c.ConfigurationPath = confPath
	return nil
}

// Validate :
func (c *ConfigurationType) Validate() error {
	return validator.Validate(c)
}

// Dumps :
func (c *ConfigurationType) Dumps() (string, error) {
	data, err := yaml.Marshal(Configuration)
	return string(data), err
}

// Configuration : global configuration
var Configuration = ConfigurationType{}

func init() {
	Configuration.Init()
}
