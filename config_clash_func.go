package main

import (
	"fmt"
	"strings"
)

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

func (t *ClashConfig) ApplyProxySelector(selectors []*ProxySelector) {
	list := make([]ClashProxyGroup, 0)
	for _, selector := range selectors {
		clashProxyGroup := ClashProxyGroup{
			Type:     selector.Type,
			Name:     selector.Name,
			Url:      selector.UrlTest,
			Interval: selector.Interval,
			Proxies:  nil,
		}
		for _, proxyItem := range t.Proxies {
			if proxyItem.IsMatchFilter(selector.ProxyFilter) {
				clashProxyGroup.Proxies = append(clashProxyGroup.Proxies, proxyItem.Name)
			}
		}
		fmt.Println("selector: ", selector.Name, " len: ", len(clashProxyGroup.Proxies))
		if len(clashProxyGroup.Proxies) > 0 {
			list = append(list, clashProxyGroup)
		}
	}
	t.ProxyGroups = append(t.ProxyGroups, list...)
}

func (t *ClashConfig) RunFilter(filter *ProxyFilter) {
	if filter == nil {
		return
	}
	t.Proxies = t.Proxies.Filter(func(item ClashProxy) bool {
		return item.IsMatchFilter(filter)
	})
	t.RemoveInvalidProxies()
}

func (t *ClashConfig) RemoveInvalidProxies() {
	proxyMapper := make(map[string]bool)
	for _, proxy := range t.GetProxies() {
		proxyMapper[proxy.Name] = true
	}

	for i := range t.ProxyGroups {
		proxies := make([]string, 0)
		for _, name := range t.ProxyGroups[i].Proxies {
			if proxyMapper[name] {
				proxies = append(proxies, name)
			}
		}
		t.ProxyGroups[i].Proxies = proxies
	}
	t.ProxyGroups = t.ProxyGroups.Filter(func(group ClashProxyGroup) bool {
		return len(group.Proxies) > 0
	})
}

type ClashProxyList []ClashProxy

func (t ClashProxyList) Filter(filter func(item ClashProxy) bool) ClashProxyList {
	ret := make(ClashProxyList, 0)
	for _, item := range t {
		if filter(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

type ClashProxyGroupList []ClashProxyGroup

func (t ClashProxyGroupList) Filter(filter func(group ClashProxyGroup) bool) ClashProxyGroupList {
	ret := make(ClashProxyGroupList, 0)
	for _, group := range t {
		if filter(group) {
			ret = append(ret, group)
		}
	}
	return ret
}

func (t *ClashProxy) IsMatchFilter(filter *ProxyFilter) bool {
	if filter == nil {
		return false
	}
	if filter.BlockName != nil {
		isBlocked := false
		for _, searchKey := range filter.BlockName {
			if strings.Contains(t.Name, searchKey) {
				isBlocked = true
				break
			}
		}
		if isBlocked {
			return false
		}
	}
	if filter.Matcher != nil {
		isMatched := false
		for _, searchKey := range filter.Matcher {
			if strings.Contains(t.Name, searchKey) {
				isMatched = true
				break
			}
		}
		return isMatched
	}
	return true
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
