package service

import (
	"context"
	"time"

	"idem/internal/model"
)

// CleanupExpiredWithContext expires old processing keys in a cancellable sweep.
// When ctx is cancelled, the sweep stops immediately and leaves the remaining
// keys untouched, so a cancellation signal actually halts further mutations.
func (s *Service) CleanupExpiredWithContext(ctx context.Context, now time.Time) (int, error) {
	items := s.store.ListIdemKeys()
	count := 0
	for _, item := range items {
		// Honor the cancellation signal before doing any work, so a cancelled
		// sweep does not keep flipping keys to "expired".
		if err := ctx.Err(); err != nil {
			return count, err
		}
		time.Sleep(time.Millisecond)
		if item.Status == model.IdemKeyStatusProcessing && !item.ExpireAt.IsZero() && now.After(item.ExpireAt) {
			item.Status = model.IdemKeyStatusExpired
			if err := s.store.UpdateIdemKey(item); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}
