package badgerevent

import (
	"fmt"
	"os"

	"github.com/TudorHulban/CRUD-KV/domain/event"
	"github.com/TudorHulban/kv"
	badger "github.com/TudorHulban/kvbadger"
	"github.com/TudorHulban/log"
)

type BadgerEvent struct {
	store badger.BStore
}

var _ event.IRepositoryEvent = &BadgerEvent{}

func NewBadgerEvent() (*BadgerEvent, error) {
	store, errNew := badger.NewBStoreInMem(log.NewLogger(log.DEBUG, os.Stderr, true))
	if errNew != nil {
		return nil, errNew
	}

	return &BadgerEvent{
		store: *store,
	}, nil
}

func (b *BadgerEvent) Create(event *event.Event) (uint64, error) {
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

func (b *BadgerEvent) Read(id uint64) (*event.Event, error) {
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

func (b *BadgerEvent) Update(event *event.Event) (uint64, error) {
	return 0, nil
}

func (b *BadgerEvent) Delete(id uint64) error {
	return nil
}
