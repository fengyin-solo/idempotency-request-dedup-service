package store

import "idem/internal/model"

func (s *MemoryStore) CreateIdemKey(i *model.IdemKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.idemKeys {
		if exist.ScopeID == i.ScopeID && exist.Key == i.Key {
			return ErrConflict
		}
	}
	s.idemKeys[i.ID] = i
	return nil
}

func (s *MemoryStore) GetIdemKey(id string) (*model.IdemKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.idemKeys[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) GetIdemKeyByScopeAndKey(scopeID, key string) (*model.IdemKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, i := range s.idemKeys {
		if i.ScopeID == scopeID && i.Key == key {
			return i, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListIdemKeys() []*model.IdemKey {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.IdemKey, 0, len(s.idemKeys))
	for _, i := range s.idemKeys {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateIdemKey(i *model.IdemKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.idemKeys[i.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.idemKeys {
		if exist.ID != i.ID && exist.ScopeID == i.ScopeID && exist.Key == i.Key {
			return ErrConflict
		}
	}
	s.idemKeys[i.ID] = i
	return nil
}

func (s *MemoryStore) DeleteIdemKey(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.idemKeys[id]; !ok {
		return ErrNotFound
	}
	delete(s.idemKeys, id)
	return nil
}

func (s *MemoryStore) DeleteIdemKeysByIDs(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		delete(s.idemKeys, id)
	}
	return nil
}
