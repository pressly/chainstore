package chainstore

import (
	"container/list"
	"errors"
)

type LRUManager struct {
	store    Store
	capacity int64 // in bytes
	cushion  int64 // 10% of bytes of the capacity, to free up this much if it hits

	items map[string]*lruItem
	list  *list.List
}

type lruItem struct {
	key         string
	size        int64
	listElement *list.Element
}

func NewLRUManager(store Store, capacity int64) (*LRUManager, error) {
	if capacity < 10 {
		return nil, errors.New("Invalid capacity, must be >= 10 bytes")
	}

	return &LRUManager{
		store:    store,
		capacity: capacity,
		cushion:  int64(float64(capacity) * 0.1),
		items:    make(map[string]*lruItem, 10000),
		list:     list.New(),
	}, nil
}

func (s *LRUManager) Open() (err error) {
	// TODO: the items list will be empty after restarting a server
	// with an existing db. We should ask the store for a list of
	// keys and their size to seed this list. Keys are easy,
	// but having a generic way to get the size of each object quickly
	// from each kind of store is challenging / over-kill (ie. s3).
	// we could persist the LRU list of keys/objects somewhere..
	// perhaps using a bolt bucket.
	return // noop
}

func (s *LRUManager) Close() (err error) {
	s.store.Close()
	return // noop
}

func (s *LRUManager) Put(key string, value []byte) (err error) {
	defer s.prune() // free up space

	valueSize := int64(len(value))

	if item, exists := s.items[key]; exists {
		s.list.MoveToFront(item.listElement)
		s.capacity += (item.size - valueSize)
		item.size = valueSize
		s.promote(item)
	} else {
		s.addItem(key, valueSize)
	}

	// TODO: what if the value is larger then even the initial capacity?
	// ..error..
	return s.store.Put(key, value)
}

func (s *LRUManager) Get(key string) (value []byte, err error) {
	value, err = s.store.Get(key)
	valueSize := len(value)
	if item, exists := s.items[key]; exists {
		s.promote(item)
	} else if valueSize > 0 {
		s.addItem(key, int64(valueSize))
	}
	return
}

func (s *LRUManager) Del(key string) (err error) {
	if item, exists := s.items[key]; exists {
		s.evict(item)
	}
	return s.store.Del(key)
}

//--

func (s *LRUManager) Capacity() int64 {
	return s.capacity
}

func (s *LRUManager) Cushion() int64 {
	return s.cushion
}

func (s *LRUManager) NumItems() int {
	return s.list.Len()
}

func (s *LRUManager) addItem(key string, size int64) {
	item := &lruItem{key: key, size: size}
	item.listElement = s.list.PushFront(item)
	s.items[key] = item
	s.capacity -= size
}

func (s *LRUManager) promote(item *lruItem) {
	s.list.MoveToFront(item.listElement)
}

func (s *LRUManager) evict(item *lruItem) {
	s.list.Remove(item.listElement)
	delete(s.items, item.key)
	s.capacity += item.size
}

func (s *LRUManager) prune() {
	if s.capacity > 0 {
		return
	}

	for s.capacity < s.cushion {
		tail := s.list.Back()
		if tail == nil {
			return
		}
		item := tail.Value.(*lruItem)
		s.Del(item.key)
	}

	if s.capacity < 0 {
		s.prune()
	}
}
