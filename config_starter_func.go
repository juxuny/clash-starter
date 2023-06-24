package main

func (t *StarterConfig) GetOverride() *ClashConfig {
	if t == nil {
		return nil
	}
	return t.Override
}

func (t *StarterConfig) GetMerge() *MergeConfig {
	if t == nil {
		return nil
	}
	return t.Merge
}

func (t *StarterConfig) GetAutoProxyGroup() *AutoGenProxyGroup {
	if t == nil {
		return nil
	}
	return t.AutoProxyGroup
}
