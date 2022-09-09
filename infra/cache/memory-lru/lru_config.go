package lru

type CfgCache struct {
	Name     string
	Capacity int
}

func NewCachesForMethods(cacheConfigurations ...CfgCache) map[string]*CacheLRU {
	res := make(map[string]*CacheLRU)

	for _, cacheConfiguration := range cacheConfigurations {
		res[cacheConfiguration.Name] = NewCacheLRU(cacheConfiguration.Capacity)
	}

	return res
}
