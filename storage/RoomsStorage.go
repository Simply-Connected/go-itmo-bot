package storage

type RoomsStorage interface {
	GetRoom(id string) (map[string]int, error)
	IncUserPoints(roomID string, userName string) error
	SetUserPoints(roomID string, userName string, value int) error
	GetUserPoints(roomID string, userName string) (int, bool)
}
