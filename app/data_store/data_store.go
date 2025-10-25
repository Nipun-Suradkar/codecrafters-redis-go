package data_store

import (
	"sync"
	"time"
)

type DataStore struct {
	DbFilename string
	DbDir      string
	dataMap    sync.Map
}

type Value struct {
	Data      interface{}
	ExpiresAt time.Time
}

func (s *DataStore) Set(key string, value interface{}, ttl time.Duration) {
	expire := time.Time{}
	if ttl > 0 {
		expire = time.Now().Add(ttl)
	}
	s.dataMap.Store(key, Value{Data: value, ExpiresAt: expire})
}

func (s *DataStore) Get(key string) (interface{}, bool) {
	val, ok := s.dataMap.Load(key)
	if !ok {
		return nil, false
	}
	v := val.(Value)
	if !v.ExpiresAt.IsZero() && v.ExpiresAt.Before(time.Now()) {
		s.dataMap.Delete(key)
		return nil, false
	}
	return v.Data, true
}
