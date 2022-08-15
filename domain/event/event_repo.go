package event

type IRepositoryEvent interface {
	Create(event *Event) (uint64, error)
	Read(id uint64) (*Event, error)
	Update(event *Event) (uint64, error)
	Delete(id uint64) error
}
