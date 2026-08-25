package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"idem/internal/config"
	"idem/internal/model"
	"idem/internal/store"
	"idem/pkg/logger"
)

func TestCleanupSweepStopsBeforeChangingKeysAfterCancel(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 100})
	sc, _ := svc.CreateScope(model.Scope{Name: "cleanup-scope"})
	for i := 0; i < 25; i++ {
		if _, err := svc.CreateIdemKey(model.IdemKey{ScopeID: sc.ID, Key: "key-" + string(rune('a'+i)), Owner: "worker", ExpireAt: time.Now().Add(-time.Minute)}); err != nil {
			t.Logf("setup error: %v", err)
			t.FailNow()
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	count, err := svc.CleanupExpiredWithContext(ctx, time.Now())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want cancellation, got count=%d err=%v", count, err)
	}
	if count != 0 {
		t.Fatalf("cancelled sweep changed %d keys", count)
	}
	items := stList(svc)
	for _, item := range items {
		if item.Status != model.IdemKeyStatusProcessing {
			t.Fatalf("key changed after cancellation: %#v", item)
		}
	}
}

func stList(svc *Service) []*model.IdemKey { return svc.store.ListIdemKeys() }
