package chainstore

type MemCache struct {
	data map[string][]byte
}

func NewMemCacheStore(capacity int64) (store *LRUManager, err error) {
	mem := &MemCache{}
	mem.Open()
	store, err = NewLRUManager(mem, capacity)
	return
}

func (s *MemCache) Open() (err error) {
	s.data = make(map[string][]byte, 1000)
	return nil
}

func (s *MemCache) Close() error {
	// s.num = 0
	return nil
}

func (s *MemCache) Put(key string, obj []byte) (err error) {
	if !IsValidKey(key) {
		return ErrInvalidKey
	}
	s.data[key] = obj
	return nil
	// s.num++
}

func (s *MemCache) Get(key string) (obj []byte, err error) {
	return s.data[key], nil
}

func (s *MemCache) Del(key string) (err error) {
	delete(s.data, key)
	return nil
}
