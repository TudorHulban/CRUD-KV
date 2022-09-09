package event

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/TudorHulban/epochid"
)

type EventData struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	ValidTo int64  `json:"validto"`
	Status  uint8  `json:"status"`
}

type Event struct {
	EventData

	FetchedFrom string
	ID          uint64 `json:"id"`
}

const _secondsValidity = 1800

var FetchedFrom = map[int]string{
	0: "cache",
	1: "persistance",
}

var _status = map[uint8]string{
	0: "draft",
	1: "active",
	2: "inactive",
}

func (e Event) String() string {
	res := []string{}

	res = append(res, fmt.Sprintf("ID: %d", e.ID))
	res = append(res, fmt.Sprintf("Fetched From: %s", e.FetchedFrom))
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

func NewEvent(options ...EventOption) *Event {
	res := Event{
		ID: epochid.NewIDIncremental10KWCoCorrection(),
	}

	for _, option := range options {
		option(&res)
	}

	now := time.Now()
	stringNow := strconv.Itoa(int(now.UnixNano()))

	res.EventData = EventData{
		Title:   "Title " + stringNow,
		Content: "Content " + stringNow,
		ValidTo: now.Add(time.Second * time.Duration(_secondsValidity)).Unix(),
		Status:  0,
	}

	return &res
}
