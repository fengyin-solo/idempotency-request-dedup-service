package service

import (
	"sort"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

func (s *Service) CreateIdemKey(input model.IdemKey) (*model.IdemKey, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetScope(input.ScopeID); err != nil {
		return nil, model.NewValidationError("scope_id", "作用域不存在")
	}
	i := &model.IdemKey{
		ID:        idgen.Hex(),
		ScopeID:   input.ScopeID,
		Key:       input.Key,
		Owner:     input.Owner,
		Status:    model.IdemKeyStatusProcessing,
		CreatedAt: time.Now(),
		ExpireAt:  input.ExpireAt,
	}
	if err := s.store.CreateIdemKey(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) ListIdemKeys(filter model.IdemKeyFilter, page, size int) ([]*model.IdemKey, int, error) {
	all := s.store.ListIdemKeys()
	matched := make([]*model.IdemKey, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.IdemKey{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetIdemKey(id string) (*model.IdemKey, error) {
	return s.store.GetIdemKey(id)
}

func (s *Service) UpdateIdemKey(id string, input model.IdemKey) (*model.IdemKey, error) {
	i, err := s.store.GetIdemKey(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != i.Status {
		if !model.IdemKeyCanTransition(i.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态流转不合法")
		}
		i.Status = input.Status
		if input.Status == model.IdemKeyStatusCompleted {
			i.CompletedAt = time.Now()
		}
	}
	if input.Key != "" {
		i.Key = input.Key
	}
	if input.Owner != "" {
		i.Owner = input.Owner
	}
	i.ExpireAt = input.ExpireAt
	if err := i.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateIdemKey(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) DeleteIdemKey(id string) error {
	return s.store.DeleteIdemKey(id)
}

func (s *Service) BatchDeleteIdemKeys(ids []string) error {
	return s.store.DeleteIdemKeysByIDs(ids)
}

func (s *Service) ExpireIdemKeys(now time.Time) (int, error) {
	all := s.store.ListIdemKeys()
	expired := 0
	for _, i := range all {
		if i.Status == model.IdemKeyStatusProcessing && !i.ExpireAt.IsZero() && now.After(i.ExpireAt) {
			i.Status = model.IdemKeyStatusExpired
			if err := s.store.UpdateIdemKey(i); err != nil {
				return expired, err
			}
			expired++
		}
	}
	return expired, nil
}
