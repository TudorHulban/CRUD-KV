package event

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventPrint(t *testing.T) {
	ev := NewEvent()

	fmt.Println(ev)
}

func TestEventJSON(t *testing.T) {
	ev := NewEvent()

	data, errMa := json.Marshal(ev)
	require.NoError(t, errMa)

	fmt.Println("Raw JSON: ", string(data))

	var reconstructed Event

	errUn := json.Unmarshal(data, &reconstructed)
	require.NoError(t, errUn)

	fmt.Println("Data:")
	fmt.Println(reconstructed)
}
