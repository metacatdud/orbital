package storage

type Storage interface {
	Get(key string, value any) error
	Set(key string, value any) error
	Del(key string) error
	Exist(key string) bool
}
