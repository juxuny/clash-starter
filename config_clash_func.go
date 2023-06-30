package main

func (t *ClashConfig) patchOverride(autoGenProxyGroups *AutoGenProxyGroup, override *ClashConfig) {
	if t == nil {
		return
	}
	if override == nil {
		return
	}
	if override.GetMixedPort() > 0 {
		t.MixedPort = override.GetMixedPort()
	}
	if override.GetExternalController() != "" {
		t.ExternalController = override.GetExternalController()
	}
	if t.GetDNS() == nil {
		t.DNS = override.GetDNS()
	}
	t.GetDNS().Patch(override.GetDNS())
	if override.GetProxies() != nil {
		t.Proxies = override.GetProxies()
	}
	if autoGenProxyGroups != nil {
		proxyGroup := ClashProxyGroup{
			Name:     autoGenProxyGroups.Name,
			Type:     autoGenProxyGroups.Type,
			Proxies:  nil,
			Url:      autoGenProxyGroups.UrlTest,
			Interval: autoGenProxyGroups.Interval,
		}
		for _, proxyItem := range t.Proxies {
			proxyGroup.Proxies = append(proxyGroup.Proxies, proxyItem.Name)
		}
		t.ProxyGroups = []ClashProxyGroup{proxyGroup}
	} else if override.GetProxyGroups() != nil {
		t.ProxyGroups = override.GetProxyGroups()
	}
	if override.GetRules() != nil {
		t.Rules = override.GetRules()
	}
}

func (t *ClashConfig) patchMerge(merge *MergeConfig) {
	if t == nil {
		return
	}
	if merge == nil {
		return
	}
	if merge.GetProxies() != nil {
		if t.Proxies == nil {
			t.Proxies = merge.GetProxies()
		} else {
			t.Proxies = append(merge.GetProxies(), t.Proxies...)
		}
	}
	t.autoRenameDuplicatedProxy()
	if merge.GetProxyGroups() != nil {
		if t.ProxyGroups == nil {
			t.ProxyGroups = merge.GetProxyGroups()
		} else {
			t.ProxyGroups = append(merge.GetProxyGroups(), t.ProxyGroups...)
		}
	}
	if merge.GetRules() != nil {
		if t.Rules == nil {
			t.Rules = merge.GetRules()
		} else {
			t.Rules = append(merge.GetRules(), t.Rules...)
		}
	}
}

func (t *ClashConfig) autoRenameDuplicatedProxy() {
	usedName := make(map[string]bool)
	for _, item := range t.Proxies {
		item.Name = autoGenNoDuplicatedName(usedName, item.Name)
		usedName[item.Name] = true
	}
}

func (t *ClashConfig) Patch(autoGenProxyGroups *AutoGenProxyGroup, override *ClashConfig, merge *MergeConfig) {
	t.patchOverride(autoGenProxyGroups, override)
	t.patchMerge(merge)
}

func (t *ClashConfig) GetMixedPort() int {
	if t == nil {
		return 0
	}
	return t.MixedPort
}

func (t *ClashConfig) GetExternalController() string {
	if t == nil {
		return ""
	}
	return t.ExternalController
}

func (t *ClashConfig) GetDNS() *ClashDNS {
	if t == nil {
		return nil
	}
	return t.DNS
}

func (t *ClashDNS) GetListen() string {
	if t == nil {
		return ""
	}
	return t.Listen
}

func (t *ClashDNS) Patch(d *ClashDNS) {
	if t == nil {
		return
	}
	if d == nil {
		return
	}
	if d.GetListen() != "" {
		t.Listen = d.GetListen()
	}
}

func (t *ClashConfig) GetProxies() []ClashProxy {
	if t == nil {
		return nil
	}
	return t.Proxies
}

func (t *ClashConfig) GetProxyGroups() []ClashProxyGroup {
	if t == nil {
		return nil
	}
	return t.ProxyGroups
}

func (t *ClashConfig) GetRules() []string {
	if t == nil {
		return nil
	}
	return t.Rules
}

func (t *MergeConfig) GetProxies() []ClashProxy {
	if t == nil {
		return nil
	}
	return t.Proxies
}

func (t *MergeConfig) GetProxyGroups() []ClashProxyGroup {
	if t == nil {
		return nil
	}
	return t.ProxyGroups
}

func (t *MergeConfig) GetRules() []string {
	if t == nil {
		return nil
	}
	return t.Rules
}
