package badgerevent

import (
	"fmt"

	"github.com/TudorHulban/CRUD-KV/domain/event"
	lru "github.com/TudorHulban/CRUD-KV/infra/cache/memory-lru"
)

var _cachesForMethods map[string]*lru.CacheLRU

var _cacheConfigObjects = []lru.CfgCache{
	{Name: "FindByID", Capacity: 10},
	{Name: "FindByIDs", Capacity: 10},
	{Name: "FindByName", Capacity: 10},
	{Name: "FindByNames", Capacity: 10},
}

func getCachesForObjectsMethods() map[string]*lru.CacheLRU {
	if _cachesForMethods == nil {
		_cachesForMethods = lru.NewCachesForMethods(_cacheConfigObjects...)
	}

	return _cachesForMethods
}

// Read would return a key not found error if key is missing.
func (b *BadgerEvent) findByID(id uint64) (*event.Event, error) {
	eventEncoded, errGet := b.store.GetVByK([]byte(fmt.Sprintf("%d", id)))
	if errGet != nil {
		return nil, errGet
	}

	var data event.EventData

	errDec := Decode(eventEncoded, &data)
	if errDec != nil {
		return nil, errDec
	}

	return &event.Event{
		ID:        id,
		EventData: data,
	}, nil
}

// Read would return a key not found error if key is missing.
func (b *BadgerEvent) FindByID(id uint64) (*event.Event, error) {
	cacheLRU := getCachesForObjectsMethods()["FindByID"]

	objectLRU := cacheLRU.Get(id)
	if objectLRU != nil {
		obj := objectLRU.(*event.Event)
		obj.FetchedFrom = event.FetchedFrom[0]

		return obj, nil
	}

	objectRetrieved, errFind := b.findByID(id)
	if errFind != nil {
		return nil, errFind
	}

	go cacheLRU.Put(id, objectRetrieved)

	objectRetrieved.FetchedFrom = event.FetchedFrom[1]

	return objectRetrieved, nil
}

// FindByIDs uses FindByID if value not cached.
// If result not found for one ID it would not error.
func (b *BadgerEvent) FindByIDs(ids ...uint64) ([]*event.Event, error) {
	var res []*event.Event

	for _, id := range ids {
		ev, errFind := b.FindByID(id)
		if errFind != nil {
			fmt.Println(errFind)

			continue
		}

		res = append(res, ev)
	}

	return res, nil
}
