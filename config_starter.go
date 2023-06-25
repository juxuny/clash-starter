package main

type StarterConfig struct {
	Link                  string             `yaml:"link"`
	Bin                   string             `yaml:"bin"`
	KeepDurationInSeconds int                `yaml:"keep-duration-in-seconds"`
	KeepNumOfFile         int                `yaml:"keep-num-of-file"`
	ConfigDir             string             `yaml:"config-dir"`
	WorkingDirectory      string             `yaml:"working-directory"`
	ReloadInterval        int                `yaml:"reload-interval"`
	AutoProxyGroup        *AutoGenProxyGroup `yaml:"auto-proxy-group"`
	Override              *ClashConfig       `yaml:"override"`
	Merge                 *MergeConfig       `yaml:"merge"`
}
type MergeConfig struct {
	Proxies     []ClashProxy      `yaml:"proxies"`
	ProxyGroups []ClashProxyGroup `yaml:"proxy-groups"`
	Rules       []string          `yaml:"rules"`
}

type AutoGenProxyGroup struct {
	Type     string `yaml:"type"`
	Name     string `yaml:"name"`
	UrlTest  string `yaml:"url-test"`
	Interval string `yaml:"interval"`
}
