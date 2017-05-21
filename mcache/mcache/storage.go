package mcache

import (
	"sync"
	"time"
)

// StoreValue - keeps value, ttl.
type StoreValue struct {
	val interface {}

	// cancel ttlDelete if key reset
	cancleTTLDelete chan struct{}
}

// Storage - implements memory key value storage.
type Storage struct {
	storeMux sync.Mutex
	store map[string]*StoreValue
}

// NewStorage - returns new storage object.
func NewStorage() *Storage {
	s := &Storage{
		store: make(map[string]*StoreValue),
	}
	return s
}

// Set - set value.
func (s *Storage) Set(k string, v interface{}, ttl time.Duration) {

	switch v.(type) {
	case string, []string, map[string]string:
		s.storeMux.Lock()

		v, ok := s.store[k]
		if ok {
			v.cancleTTLDelete <- struct{}{}
		}

		ch := make(chan struct{}, 1)
		s.store[k] = &StoreValue{
			val: v,
			cancleTTLDelete: ch,
		}
		go s.ttlDelete(k, ttl, ch)

		s.storeMux.Unlock()

	default:
		panic("you can set only values with type - string, []string, map[string]string")
	}
}

// Get - return *StoreValue.
func (s *Storage) Get(k string) *StoreValue {
	s.storeMux.Lock()
	v, ok := s.store[k]
	s.storeMux.Unlock()
	if !ok {
		return nil
	}
	return v
}

// Remove - delete value by key.
func (s *Storage) Remove(k string) {
	s.storeMux.Lock()
	v, ok := s.store[k]
	if ok {
		v.cancleTTLDelete <- struct{}{}
	}
	delete(s.store, k)
	s.storeMux.Unlock()
}

// Keys - return store keys.
func (s *Storage) Keys() []string {
	keys := []string{}
	s.storeMux.Lock()
	for k := range s.store {
		keys = append(keys, k)
	}
	s.storeMux.Unlock()
	return keys
}

// ttlDelete - delete from store after ttl expired, if key reset goroutine is canceled.
func (s *Storage) ttlDelete(k string, ttl time.Duration, cl chan struct{}) {
	//
	time.Sleep(ttl)

	//
	select {
	case <- cl:
	default:
		s.storeMux.Lock()
		delete(s.store, k)
		s.storeMux.Unlock()
	}
}
