package store

import "idem/internal/model"

func (s *MemoryStore) CreateRequestRecord(r *model.RequestRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.requestRecords[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRequestRecord(id string) (*model.RequestRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.requestRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRequestRecords() []*model.RequestRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RequestRecord, 0, len(s.requestRecords))
	for _, r := range s.requestRecords {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRequestRecord(r *model.RequestRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.requestRecords[r.ID]; !ok {
		return ErrNotFound
	}
	s.requestRecords[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRequestRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.requestRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.requestRecords, id)
	return nil
}
