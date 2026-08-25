package service

import (
	"sort"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

func (s *Service) CreateTTLPolicy(input model.TTLPolicy) (*model.TTLPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetScope(input.ScopeID); err != nil {
		return nil, model.NewValidationError("scope_id", "作用域不存在")
	}
	t := &model.TTLPolicy{
		ID:          idgen.Hex(),
		ScopeID:     input.ScopeID,
		TTLSeconds:  input.TTLSeconds,
		MaxEntries:  input.MaxEntries,
		AutoCleanup: input.AutoCleanup,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateTTLPolicy(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTTLPolicies(filter model.TTLPolicyFilter, page, size int) ([]*model.TTLPolicy, int, error) {
	all := s.store.ListTTLPolicies()
	matched := make([]*model.TTLPolicy, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TTLPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetTTLPolicy(id string) (*model.TTLPolicy, error) {
	return s.store.GetTTLPolicy(id)
}

func (s *Service) UpdateTTLPolicy(id string, input model.TTLPolicy) (*model.TTLPolicy, error) {
	t, err := s.store.GetTTLPolicy(id)
	if err != nil {
		return nil, err
	}
	t.TTLSeconds = input.TTLSeconds
	t.MaxEntries = input.MaxEntries
	t.AutoCleanup = input.AutoCleanup
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTTLPolicy(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTTLPolicy(id string) error {
	return s.store.DeleteTTLPolicy(id)
}
