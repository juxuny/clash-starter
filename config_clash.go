package main

import "strings"

type ClashConfig struct {
	MixedPort          int                 `yaml:"mixed-port"`
	AllowLan           bool                `yaml:"allow-lan"`
	Mode               string              `yaml:"mode"`
	LogLevel           string              `yaml:"log-level"`
	ExternalController string              `yaml:"external-controller"`
	Secret             string              `yaml:"secret"`
	DNS                *ClashDNS           `yaml:"dns"`
	Proxies            ClashProxyList      `yaml:"proxies,omitempty"`
	ProxyGroups        ClashProxyGroupList `yaml:"proxy-groups"`
	Rules              []string            `yaml:"rules"`
}

type ClashDNS struct {
	Enable         bool                 `yaml:"enable,omitempty"`
	Ipv6           bool                 `yaml:"ipv6,omitempty"`
	Listen         string               `yaml:"listen,omitempty"`
	EnhancedMode   string               `yaml:"enhanced-mode,omitempty"`
	FakeIpFilter   []string             `yaml:"fake-ip-filter,omitempty"`
	Nameserver     []string             `yaml:"nameserver,omitempty"`
	Fallback       []string             `yaml:"fallback,omitempty"`
	FallbackFilter *ClashFallbackFilter `yaml:"fallback-filter,omitempty"`
}

type ClashFallbackFilter struct {
	GeoIP  bool     `yaml:"geoip,omitempty"`
	IpAddr []string `yaml:"ipaddr,omitempty"`
	Domain []string `yaml:"domain,omitempty"`
}

type ClashProxy struct {
	Name           string            `yaml:"name,omitempty"`
	Type           string            `yaml:"type,omitempty"`
	Server         string            `yaml:"server,omitempty"`
	Port           int               `yaml:"port,omitempty"`
	Uuid           string            `yaml:"uuid,omitempty"`
	AlterId        int64             `yaml:"alterId"`
	Cipher         string            `yaml:"cipher,omitempty"`
	UDP            bool              `yaml:"udp,omitempty"`
	TLS            bool              `yaml:"tls,omitempty"`
	SkipCertVerify bool              `yaml:"skip-cert-verify,omitempty"`
	Network        string            `yaml:"network,omitempty"`
	WsHeaders      map[string]string `yaml:"ws-headers,omitempty"`
	WsPath         string            `yaml:"ws-path,omitempty"`
	WsOpts         *WsOpts           `yaml:"ws-opts,omitempty"`
}

type WsOpts struct {
	Path    string            `yaml:"path"`
	Headers map[string]string `yaml:"headers"`
}

type ClashProxyGroup struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Proxies  []string `yaml:"proxies,omitempty"`
	Url      string   `yaml:"url,omitempty"`
	Interval string   `yaml:"interval,omitempty"`
}

func (t *ClashConfig) GetControlPanelEntrypoint() string {
	if t.ExternalController == "" {
		return "http://127.0.0.1:9090"
	}
	if strings.Index(t.ExternalController, "http") == 0 {
		return t.ExternalController
	}
	if strings.Index(t.ExternalController, "[::]") == 0 {
		t.ExternalController = strings.Replace(t.ExternalController, "[::]", "127.0.0.1", 1)
	}
	return "http://" + t.GetExternalController()
}
