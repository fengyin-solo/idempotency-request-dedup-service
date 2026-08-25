package service

import (
	"sort"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

func (s *Service) CreateDedupHit(input model.DedupHit) (*model.DedupHit, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetIdemKey(input.IdemKeyID); err != nil {
		return nil, model.NewValidationError("idem_key_id", "幂等键不存在")
	}
	d := &model.DedupHit{
		ID:                idgen.Hex(),
		IdemKeyID:         input.IdemKeyID,
		OriginalRequestID: input.OriginalRequestID,
		RequestHash:       input.RequestHash,
		HitAt:             time.Now(),
	}
	if err := s.store.CreateDedupHit(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListDedupHits(filter model.DedupHitFilter, page, size int) ([]*model.DedupHit, int, error) {
	all := s.store.ListDedupHits()
	matched := make([]*model.DedupHit, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].HitAt.After(matched[j].HitAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.DedupHit{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetDedupHit(id string) (*model.DedupHit, error) {
	return s.store.GetDedupHit(id)
}

func (s *Service) DeleteDedupHit(id string) error {
	return s.store.DeleteDedupHit(id)
}

func (s *Service) BatchDeleteDedupHits(ids []string) error {
	return s.store.DeleteDedupHitsByIDs(ids)
}
