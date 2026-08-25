package store

import "idem/internal/model"

func (s *MemoryStore) CreateTTLPolicy(t *model.TTLPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.ttlPolicies {
		if exist.ScopeID == t.ScopeID {
			return ErrConflict
		}
	}
	s.ttlPolicies[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTTLPolicy(id string) (*model.TTLPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.ttlPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTTLPolicies() []*model.TTLPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TTLPolicy, 0, len(s.ttlPolicies))
	for _, t := range s.ttlPolicies {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTTLPolicy(t *model.TTLPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ttlPolicies[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.ttlPolicies {
		if exist.ID != t.ID && exist.ScopeID == t.ScopeID {
			return ErrConflict
		}
	}
	s.ttlPolicies[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTTLPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.ttlPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.ttlPolicies, id)
	return nil
}
