package pudge

import (
	"github.com/Skazzi00/go-itmo-bot/storage"
	"github.com/recoilme/pudge"
	"sync"
)

const dbName = "rooms"

type RoomsStorage struct {
	db   *pudge.Db
	lock *sync.RWMutex
}

func NewRoomsStorage() storage.RoomsStorage {
	db, _ := pudge.Open(dbName, &pudge.Config{SyncInterval: 3})
	return &RoomsStorage{db, &sync.RWMutex{}}
}

func (r *RoomsStorage) getRoomUnsafe(id string) (map[string]int, error) {
	room := make(map[string]int)
	err := r.db.Get(id, &room)
	if err == pudge.ErrKeyNotFound {
		return room, r.db.Set(id, room)
	}
	return room, err
}

func (r *RoomsStorage) GetRoom(id string) (map[string]int, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.getRoomUnsafe(id)
}

func (r *RoomsStorage) IncUserPoints(roomID string, userName string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	room, err := r.getRoomUnsafe(roomID)
	if err != nil {
		return err
	}
	room[userName]++
	return r.db.Set(roomID, room)
}

func (r *RoomsStorage) GetUserPoints(roomID string, userName string) (int, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	room, err := r.getRoomUnsafe(roomID)
	if err != nil {
		return 0, false
	}
	val, ok := room[userName]
	return val, ok
}

func (r *RoomsStorage) SetUserPoints(roomID string, userName string, value int) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	room, err := r.getRoomUnsafe(roomID)
	if err != nil {
		return err
	}
	room[userName] = value
	return r.db.Set(roomID, room)
}
