package service

import (
	"testing"
	"time"

	"idem/internal/config"
	"idem/internal/store"
	"idem/pkg/logger"
)

func TestEnqueuedHashesKeepRequestSnapshot(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 20})
	hashes := []string{"first", "second"}
	queued := svc.EnqueueRequestHashes(hashes)
	hashes[0] = "changed-by-next-request"
	hashes[1] = "also-changed"
	select {
	case got := <-queued:
		if len(got) != 2 || got[0] != "first" || got[1] != "second" {
			t.Fatalf("queued request was polluted: %#v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("queued request did not complete")
	}
}
