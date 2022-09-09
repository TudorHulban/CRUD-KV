package badgerevent

import (
	"fmt"

	"github.com/TudorHulban/CRUD-KV/domain/event"
	"github.com/TudorHulban/CRUD-KV/infra/repository"
	"github.com/TudorHulban/kv"
	badger "github.com/TudorHulban/kvbadger"
	"github.com/TudorHulban/log"
)

type BadgerEvent struct {
	store badger.BStore
}

var _ repository.IRepositoryEvent = &BadgerEvent{}

func NewBadgerEvent(logger *log.Logger) (*BadgerEvent, error) {
	store, errNew := badger.NewBStoreInMem(logger)
	if errNew != nil {
		return nil, errNew
	}

	return &BadgerEvent{
		store: *store,
	}, nil
}

func (b *BadgerEvent) Insert(event *event.Event) (uint64, error) {
	eventEncoded, errEnc := Encode(event.EventData)
	if errEnc != nil {
		return 0, errEnc
	}

	kv := kv.KV{
		Key:   []byte(fmt.Sprintf("%d", event.ID)),
		Value: eventEncoded,
	}

	errSet := b.store.Set(kv)
	if errSet != nil {
		return 0, errSet
	}

	return event.ID, nil
}

func (b *BadgerEvent) Update(event *event.Event) (uint64, error) {
	_, errGet := b.store.GetVByK([]byte(fmt.Sprintf("%d", event.ID)))
	if errGet != nil {
		return 0, errGet
	}

	eventEncoded, errEnc := Encode(event.EventData)
	if errEnc != nil {
		return 0, errEnc
	}

	kv := kv.KV{
		Key:   []byte(fmt.Sprintf("%d", event.ID)),
		Value: eventEncoded,
	}

	errSet := b.store.Set(kv)
	if errSet != nil {
		return 0, errSet
	}

	go func() {
		cacheLRU := getCachesForObjectsMethods()["FindByID"]

		cacheLRU.Put(event.ID, event.EventData)
	}()

	return event.ID, nil
}

func (b *BadgerEvent) Delete(id uint64) error {
	go func() {
		cacheLRU := getCachesForObjectsMethods()["FindByID"]

		cacheLRU.Delete(id)
	}()

	return b.store.DeleteKVByK([]byte(fmt.Sprintf("%d", id)))
}
