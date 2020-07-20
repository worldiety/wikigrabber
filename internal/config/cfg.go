package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Config represents the grabbers settings.
// Example yaml:
//    # your remote git repos, which will be pulled
//    remote:
//    - https://github.com/worldiety/wikigrabber
//    - https://github.com/worldiety/goup
//
//    # your local folders, which are just interpreted as-is
//    local:
//    - /Users/user/home/git/my/project
//
//    # either an https git link or a local folder which contains the template files
//    template:
//
//    # the local server to launch and serve the merged results from
//    port: 8080
type Config struct {
	Remote   []string `yaml:"remote"`
	Local    []string `yaml:"local"`
	Port     int      `yaml:"port"`
	TmpDir   string   `yaml:"tmpDir"`
	Template string   `yaml:"template"`
}

// Load parses the file as yaml
func (c *Config) Load(fname string) error {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}

	return nil
}
