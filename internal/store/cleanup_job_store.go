package store

import "idem/internal/model"

func (s *MemoryStore) CreateCleanupJob(c *model.CleanupJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cleanupJobs[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCleanupJob(id string) (*model.CleanupJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.cleanupJobs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListCleanupJobs() []*model.CleanupJob {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.CleanupJob, 0, len(s.cleanupJobs))
	for _, c := range s.cleanupJobs {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCleanupJob(c *model.CleanupJob) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cleanupJobs[c.ID]; !ok {
		return ErrNotFound
	}
	s.cleanupJobs[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCleanupJob(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cleanupJobs[id]; !ok {
		return ErrNotFound
	}
	delete(s.cleanupJobs, id)
	return nil
}
