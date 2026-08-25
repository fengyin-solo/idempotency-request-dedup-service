package store

import "idem/internal/model"

func (s *MemoryStore) CreateRequestRecord(r *model.RequestRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 存入拷贝，避免外部指针成为存储的内部状态，
	// 防止调用方在创建后继续修改对象从而引发并发安全问题。
	cp := *r
	s.requestRecords[r.ID] = &cp
	return nil
}

func (s *MemoryStore) GetRequestRecord(id string) (*model.RequestRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.requestRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	// 返回拷贝，确保读取方拿到的对象与存储内部隔离，
	// 不会被并发的更新改写到同一对象而触发 data race。
	cp := *r
	return &cp, nil
}

func (s *MemoryStore) ListRequestRecords() []*model.RequestRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RequestRecord, 0, len(s.requestRecords))
	for _, r := range s.requestRecords {
		cp := *r
		list = append(list, &cp)
	}
	return list
}

func (s *MemoryStore) UpdateRequestRecord(r *model.RequestRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.requestRecords[r.ID]; !ok {
		return ErrNotFound
	}
	// 以整体替换的方式写入拷贝：单次 Update 对应一次原子的全量赋值，
	// 不会把字段级中间态暴露给后续读取，从而避免哈希来自中间状态。
	cp := *r
	s.requestRecords[r.ID] = &cp
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
