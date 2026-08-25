package store

import "idem/internal/model"

func (s *MemoryStore) CreateScope(sc *model.Scope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.scopes {
		if exist.Name == sc.Name {
			return ErrConflict
		}
	}
	s.scopes[sc.ID] = sc
	return nil
}

func (s *MemoryStore) GetScope(id string) (*model.Scope, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sc, ok := s.scopes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sc, nil
}

func (s *MemoryStore) ListScopes() []*model.Scope {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Scope, 0, len(s.scopes))
	for _, sc := range s.scopes {
		list = append(list, sc)
	}
	return list
}

func (s *MemoryStore) UpdateScope(sc *model.Scope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scopes[sc.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.scopes {
		if exist.ID != sc.ID && exist.Name == sc.Name {
			return ErrConflict
		}
	}
	s.scopes[sc.ID] = sc
	return nil
}

func (s *MemoryStore) DeleteScope(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scopes[id]; !ok {
		return ErrNotFound
	}
	delete(s.scopes, id)
	return nil
}
