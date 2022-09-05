package badgerevent

import (
	"os"
	"testing"

	"github.com/TudorHulban/CRUD-KV/domain/event"
	"github.com/TudorHulban/log"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

func TestBadgerEvent(t *testing.T) {
	logger := log.NewLogger(log.DEBUG, os.Stderr, true)

	repo, errNew := NewBadgerEvent(logger)
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

	repo.Delete(id)

	shouldBeNil, errNil := repo.Read(ev.ID)
	require.Error(t, errNil)
	require.Nil(t, shouldBeNil)
}

func TestBadgerEventNotFound(t *testing.T) {
	logger := log.NewLogger(log.DEBUG, os.Stderr, true)

	repo, errNew := NewBadgerEvent(logger)
	require.NoError(t, errNew)
	require.NotNil(t, repo)
	defer repo.store.Close()

	eventNotFound, errNil := repo.Read(1)
	require.Equal(t, errNil.Error(), "Key not found")
	require.Nil(t, eventNotFound)
}

func TestBadgerEventUpdate(t *testing.T) {
	logger := log.NewLogger(log.DEBUG, os.Stderr, true)

	repo, errNew := NewBadgerEvent(logger)
	require.NoError(t, errNew)
	require.NotNil(t, repo)
	defer repo.store.Close()

	ev1 := event.NewEvent()
	id1, errCre := repo.Create(ev1)
	require.NoError(t, errCre)
	require.Equal(t, ev1.ID, id1)

	ev2 := event.NewEvent()
	ev2.ID = ev1.ID

	id2, errUpd := repo.Update(ev2)
	require.NoError(t, errUpd)
	require.Equal(t, id1, id2)
}
