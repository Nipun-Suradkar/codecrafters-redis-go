package redis_server

import (
	"sync"
	"time"
)

// Store defines the behavior of a key-value store.
type Store interface {
	Set(key string, value interface{}, ttl time.Duration)
	Get(key string) (interface{}, bool)
	Delete(key string)
}

type Value struct {
	Data      interface{}
	ExpiresAt time.Time
}

type DataStore struct {
	data *sync.Map
}

var (
	globalDataStore Store
	once            sync.Once
)

func GetDataStore() Store {
	once.Do(func() {
		globalDataStore = &DataStore{
			data: &sync.Map{},
		}
	})
	return globalDataStore
}

func (d *DataStore) Set(key string, value interface{}, ttl time.Duration) {
	expire := time.Time{}
	if ttl > 0 {
		expire = time.Now().Add(ttl)
	}
	d.data.Store(key, Value{Data: value, ExpiresAt: expire})
}

func (d *DataStore) Get(key string) (interface{}, bool) {
	val, ok := d.data.Load(key)
	if !ok {
		return nil, false
	}
	v := val.(Value)
	if !v.ExpiresAt.IsZero() && v.ExpiresAt.Before(time.Now()) {
		d.data.Delete(key)
		return nil, false
	}
	return v.Data, true
}

func (d *DataStore) Delete(key string) {
	d.data.Delete(key)
}
