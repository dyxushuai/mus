package db


type IStorage interface {
	Keys(string) ([]string, error)

	GetServer(string) ([]byte, error)

	GetServers(string) ([][]byte, error)

	SetServer(string, []byte) error

	DelServer(string) error

	IncrSize(string, int) (int64, error)

	GetSize(string) (int64, error)
}
