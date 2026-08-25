package service

import (
	"sort"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

func (s *Service) CreateScope(input model.Scope) (*model.Scope, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sc := &model.Scope{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}
	if err := s.store.CreateScope(sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *Service) ListScopes(filter model.ScopeFilter, page, size int) ([]*model.Scope, int, error) {
	all := s.store.ListScopes()
	matched := make([]*model.Scope, 0, len(all))
	for _, sc := range all {
		if filter.Match(sc) {
			matched = append(matched, sc)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Scope{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetScope(id string) (*model.Scope, error) {
	return s.store.GetScope(id)
}

func (s *Service) UpdateScope(id string, input model.Scope) (*model.Scope, error) {
	sc, err := s.store.GetScope(id)
	if err != nil {
		return nil, err
	}
	sc.Name = input.Name
	sc.Description = input.Description
	if err := sc.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateScope(sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *Service) DeleteScope(id string) error {
	return s.store.DeleteScope(id)
}
