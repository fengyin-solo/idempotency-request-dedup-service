package service

import (
	"sort"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

func (s *Service) CreateRequestRecord(input model.RequestRecord) (*model.RequestRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetIdemKey(input.IdemKeyID); err != nil {
		return nil, model.NewValidationError("idem_key_id", "幂等键不存在")
	}
	r := &model.RequestRecord{
		ID:           idgen.Hex(),
		IdemKeyID:    input.IdemKeyID,
		RequestHash:  input.RequestHash,
		ResponseHash: input.ResponseHash,
		Status:       input.Status,
		CreatedAt:    time.Now(),
	}
	if err := s.store.CreateRequestRecord(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) ListRequestRecords(filter model.RequestRecordFilter, page, size int) ([]*model.RequestRecord, int, error) {
	all := s.store.ListRequestRecords()
	matched := make([]*model.RequestRecord, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RequestRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetRequestRecord(id string) (*model.RequestRecord, error) {
	return s.store.GetRequestRecord(id)
}

func (s *Service) UpdateRequestRecord(id string, input model.RequestRecord) (*model.RequestRecord, error) {
	// GetRequestRecord 返回的是存储内部对象的拷贝，在此处独立构建完整的新状态，
	// 不与其它并发更新共享同一对象，避免 data race 与字段级中间态相互覆盖。
	r, err := s.store.GetRequestRecord(id)
	if err != nil {
		return nil, err
	}
	if input.RequestHash != "" {
		r.RequestHash = input.RequestHash
	}
	r.ResponseHash = input.ResponseHash
	r.Status = input.Status
	if err := r.Validate(); err != nil {
		return nil, err
	}
	// 整体替换写入：一次 Update 对应一次原子的全量赋值，
	// 读侧永远不会观察到字段级中间态，最终落库的哈希必来自某一次完整更新。
	if err := s.store.UpdateRequestRecord(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRequestRecord(id string) error {
	return s.store.DeleteRequestRecord(id)
}
