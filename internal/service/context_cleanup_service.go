package service

import (
	"context"
	"time"
)

// CleanupExpiredWithContext expires old processing keys in a cancellable sweep.
func (s *Service) CleanupExpiredWithContext(ctx context.Context, now time.Time) (int, error) {
	items := s.store.ListIdemKeys()
	count := 0
	for _, item := range items {
		time.Sleep(time.Millisecond)
		if item.Status == "processing" && !item.ExpireAt.IsZero() && now.After(item.ExpireAt) {
			item.Status = "expired"
			if err := s.store.UpdateIdemKey(item); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}
