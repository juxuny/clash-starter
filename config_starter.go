package main

type StarterConfig struct {
	Port      int         `yaml:"port"`
	Link      string      `yaml:"link"`
	Bin       string      `yaml:"bin"`
	ConfigDir string      `yaml:"config-dir"`
	Test      *TestConfig `yaml:"test"`
}

type TestConfig struct {
	Port     int    `yaml:"port"`
	TestUrl  string `yaml:"test-url"`
	Interval string `yaml:"interval"`
}
