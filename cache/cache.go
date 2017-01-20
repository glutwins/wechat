package cache

type Cache interface {
	Set(string, interface{}) error
	Get(string, interface{}) error
	Lock(string) error
	Unlock(string) error
}
