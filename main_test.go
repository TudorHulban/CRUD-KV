package main

import (
	"fmt"
	"testing"
)

func TestEvent(t *testing.T) {
	ev := NewEvent()

	fmt.Println(ev)
}
