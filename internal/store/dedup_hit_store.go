package store

import "idem/internal/model"

func (s *MemoryStore) CreateDedupHit(d *model.DedupHit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dedupHits[d.ID] = d
	return nil
}

func (s *MemoryStore) GetDedupHit(id string) (*model.DedupHit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.dedupHits[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListDedupHits() []*model.DedupHit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.DedupHit, 0, len(s.dedupHits))
	for _, d := range s.dedupHits {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) DeleteDedupHit(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.dedupHits[id]; !ok {
		return ErrNotFound
	}
	delete(s.dedupHits, id)
	return nil
}

func (s *MemoryStore) DeleteDedupHitsByIDs(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		delete(s.dedupHits, id)
	}
	return nil
}
