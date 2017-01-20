package cache

import (
	"bytes"
	//"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

func NewJsonFileCache(dsn string) Cache {
	cache := &fileCache{}
	cache.lockers = make(map[string]*sync.Mutex)
	cache.data = make(map[string][]byte)
	cache.path = dsn

	cache.load()
	return cache
}

type fileCache struct {
	gmutex  sync.RWMutex
	lockers map[string]*sync.Mutex
	data    map[string][]byte
	path    string
}

func (cache *fileCache) load() error {
	b, err := ioutil.ReadFile(cache.path)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewReader(b)).Decode(&cache.data)
}

func (cache *fileCache) save() error {
	w, err := os.OpenFile(cache.path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer w.Close()

	return json.NewEncoder(w).Encode(cache.data)
}

func (cache *fileCache) Get(key string, val interface{}) error {
	cache.gmutex.RLock()
	defer cache.gmutex.RUnlock()
	if b, ok := cache.data[key]; ok {
		return json.NewDecoder(bytes.NewReader(b)).Decode(val)
	}
	return errors.New("key not exists")
}

func (cache *fileCache) Set(key string, val interface{}) error {
	cache.gmutex.Lock()
	defer cache.gmutex.Unlock()

	w := bytes.NewBuffer(nil)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		return err
	}
	cache.data[key] = w.Bytes()
	return cache.save()
}

func (cache *fileCache) Lock(key string) error {
	cache.gmutex.RLock()
	mu, ok := cache.lockers[key]
	cache.gmutex.RUnlock()

	if ok {
		mu.Lock()
	} else {
		cache.gmutex.Lock()
		defer cache.gmutex.Unlock()
		if mu, ok = cache.lockers[key]; ok {
			mu.Lock()
		} else {
			mu = &sync.Mutex{}
			cache.lockers[key] = mu
			mu.Lock()
		}
	}

	return nil
}

func (cache *fileCache) Unlock(key string) error {
	cache.gmutex.RLock()
	mu, ok := cache.lockers[key]
	cache.gmutex.RUnlock()
	if ok {
		mu.Unlock()
	}
	return nil
}
