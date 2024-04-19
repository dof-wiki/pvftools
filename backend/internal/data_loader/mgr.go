package data_loader

import "github.com/samber/lo"

type CheckUpdateResult struct {
	Key        string `json:"key"`
	Name       string `json:"name"`
	NeedUpdate bool   `json:"need_update"`
}

type dataLoaderMgr struct {
	loaders map[string]ILoader
}

func (m *dataLoaderMgr) LoadData(keys ...string) {
	if len(keys) == 0 {
		keys = lo.Keys(m.loaders)
	}
	for _, k := range keys {
		loader := m.loaders[k]
		if loader == nil {
			continue
		}
		loader.Load()
	}
}

func (m *dataLoaderMgr) CheckUpdate() []*CheckUpdateResult {
	ret := make([]*CheckUpdateResult, 0, len(m.loaders))
	for k, loader := range m.loaders {
		ret = append(ret, &CheckUpdateResult{
			Key:        k,
			Name:       loader.Name(),
			NeedUpdate: loader.CheckUpdate(),
		})
	}
	return ret
}

var DataLoaderMgr = &dataLoaderMgr{
	loaders: make(map[string]ILoader),
}

func register(key string, loader ILoader) {
	DataLoaderMgr.loaders[key] = loader
}
