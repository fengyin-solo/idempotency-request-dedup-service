package service

import (
	"sort"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

func (s *Service) CreateCleanupJob(input model.CleanupJob) (*model.CleanupJob, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTTLPolicy(input.PolicyID); err != nil {
		return nil, model.NewValidationError("policy_id", "策略不存在")
	}
	c := &model.CleanupJob{
		ID:           idgen.Hex(),
		PolicyID:     input.PolicyID,
		DeletedCount: 0,
		Status:       model.CleanupJobStatusRunning,
		Message:      input.Message,
		RanAt:        time.Now(),
	}
	if err := s.store.CreateCleanupJob(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListCleanupJobs(filter model.CleanupJobFilter, page, size int) ([]*model.CleanupJob, int, error) {
	all := s.store.ListCleanupJobs()
	matched := make([]*model.CleanupJob, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].RanAt.After(matched[j].RanAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.CleanupJob{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCleanupJob(id string) (*model.CleanupJob, error) {
	return s.store.GetCleanupJob(id)
}

func (s *Service) UpdateCleanupJob(id string, input model.CleanupJob) (*model.CleanupJob, error) {
	c, err := s.store.GetCleanupJob(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" && input.Status != c.Status {
		if !model.CleanupJobCanTransition(c.Status, input.Status) {
			return nil, model.NewValidationError("status", "状态流转不合法")
		}
		c.Status = input.Status
	}
	c.DeletedCount = input.DeletedCount
	c.Message = input.Message
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCleanupJob(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCleanupJob(id string) error {
	return s.store.DeleteCleanupJob(id)
}
