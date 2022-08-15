package badgerevent

import (
	"testing"

	"github.com/TudorHulban/CRUD-KV/domain/event"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func TestBadgerEvent(t *testing.T) {
	repo, errNew := NewBadgerEvent()
	require.NoError(t, errNew)
	require.NotNil(t, repo)
	defer repo.store.Close()

	ev := event.NewEvent()
	id, errCre := repo.Create(ev)
	require.NoError(t, errCre)
	require.Equal(t, ev.ID, id)

	reconstructed, errRe := repo.Read(ev.ID)
	require.NoError(t, errRe)

	t.Log("reconstructed:", reconstructed)
	t.Log("event:", ev)

	require.Zero(t, deep.Equal(reconstructed, ev))
}
