package config

type Configuration struct {
	Tel           string `yaml:"tel"`
	Password      string `yaml:"password"`
	CheckMode     string `yaml:"check-mode"`
	CheckInterval uint   `yaml:"check-interval"`
}
