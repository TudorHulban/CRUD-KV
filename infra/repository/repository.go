package repository

import "github.com/TudorHulban/CRUD-KV/domain/event"

type IRepositoryEvent interface {
	Insert(event *event.Event) (uint64, error)
	FindByID(id uint64) (*event.Event, error)
	FindByIDs(ids ...uint64) ([]*event.Event, error)
	Update(event *event.Event) (uint64, error)
	Delete(id uint64) error
}
