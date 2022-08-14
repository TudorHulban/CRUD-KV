package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TudorHulban/epochid"
)

type Event struct {
	Title   string `json:"title"`
	Content string `json:"content"`

	ValidTo int64 `json:"validto"`

	ID     uint64 `json:"id"`
	Status uint8  `json:"status"`
}

type RepositoryEvent interface {
	Create(event *Event) (uint, error)
	Read(id uint) (*Event, error)
	Update(event *Event) (uint, error)
	Delete(id uint) error
}

const _secondsValidity = 1800

var _status = map[uint8]string{
	0: "draft",
	1: "active",
	2: "inactive",
}

func (e Event) String() string {
	res := []string{}

	res = append(res, fmt.Sprintf("ID: %d", e.ID))
	res = append(res, fmt.Sprintf("Title: %s", e.Title))
	res = append(res, fmt.Sprintf("Content: %s", e.Content))
	res = append(res, fmt.Sprintf("Valid to: %s", time.Unix(e.ValidTo, 0)))
	res = append(res, fmt.Sprintf("Status: %s", _status[e.Status]))

	return strings.Join(res, "\n")
}

// MarshalJSON provides custom marshalling, also avoiding memory allignment issues.
func (e Event) MarshalJSON() ([]byte, error) {
	event := struct {
		ID      uint64 `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Status  uint8  `json:"status"`
		ValidTo int64  `json:"validto"`
	}{
		ID:      e.ID,
		Title:   e.Title,
		Content: e.Content,
		Status:  e.Status,
		ValidTo: e.ValidTo,
	}

	return json.Marshal(&event)
}

func NewEvent() *Event {
	now := time.Now()
	stringNow := strconv.Itoa(int(now.UnixNano()))

	return &Event{
		ID:      epochid.NewIDIncremental10KWCoCorrection(),
		Title:   "Title " + stringNow,
		Content: "Content " + stringNow,
		ValidTo: now.Add(time.Second * time.Duration(_secondsValidity)).Unix(),
		Status:  0,
	}
}

func main() {}
