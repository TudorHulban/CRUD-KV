package badgerevent

import (
	"os"
	"testing"
	"time"

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
	id, errCre := repo.Insert(ev)
	require.NoError(t, errCre)
	require.Equal(t, ev.ID, id)

	reconstructed, errRe := repo.FindByID(ev.ID)
	require.NoError(t, errRe)

	t.Log("reconstructed:", reconstructed)
	t.Log("event:", ev)

	ev.FetchedFrom = event.FetchedFrom[1]

	require.Zero(t, deep.Equal(reconstructed, ev))

	repo.Delete(id)

	time.Sleep(100 * time.Millisecond) //for cache sync

	shouldBeNil, errNil := repo.FindByID(ev.ID)
	require.Error(t, errNil, shouldBeNil)
	require.Nil(t, shouldBeNil)
}

func TestBadgerEventNotFound(t *testing.T) {
	logger := log.NewLogger(log.DEBUG, os.Stderr, true)

	repo, errNew := NewBadgerEvent(logger)
	require.NoError(t, errNew)
	require.NotNil(t, repo)
	defer repo.store.Close()

	eventNotFound, errNil := repo.FindByID(1)
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
	id1, errCre := repo.Insert(ev1)
	require.NoError(t, errCre)
	require.Equal(t, ev1.ID, id1)

	ev2 := event.NewEvent()
	ev2.ID = ev1.ID

	id2, errUpd := repo.Update(ev2)
	require.NoError(t, errUpd)
	require.Equal(t, id1, id2)
}

func TestBadgerEventFindByIDs(t *testing.T) {
	logger := log.NewLogger(log.DEBUG, os.Stderr, true)

	repo, errNew := NewBadgerEvent(logger)
	require.NoError(t, errNew)
	require.NotNil(t, repo)
	defer repo.store.Close()

	ev1 := event.NewEvent()
	id1, errCre1 := repo.Insert(ev1)
	require.NoError(t, errCre1)
	require.Equal(t, ev1.ID, id1)

	ev2 := event.NewEvent()
	id2, errCre2 := repo.Insert(ev2)
	require.NoError(t, errCre2)
	require.Equal(t, ev2.ID, id2)

	reconstructed1, errRe1 := repo.FindByID(ev1.ID)
	require.NoError(t, errRe1)
	require.Equal(t, event.FetchedFrom[1], reconstructed1.FetchedFrom)

	reconstructed2, errRe2 := repo.FindByID(ev2.ID)
	require.NoError(t, errRe2)
	require.Equal(t, event.FetchedFrom[1], reconstructed2.FetchedFrom)

	time.Sleep(100 * time.Millisecond) //for cache sync

	reconstructed1Cache, errRe1Cache := repo.FindByID(ev1.ID)
	require.NoError(t, errRe1Cache)
	require.Equal(t, event.FetchedFrom[0], reconstructed1Cache.FetchedFrom)

	reconstructed2Cache, errRe2Cache := repo.FindByID(ev2.ID)
	require.NoError(t, errRe2Cache)
	require.Equal(t, event.FetchedFrom[0], reconstructed2Cache.FetchedFrom)
}
