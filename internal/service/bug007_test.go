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

func TestProcessRequestHonorsCancellation(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 20})
	sc, _ := svc.CreateScope(model.Scope{Name: "process-scope"})
	ik, err := svc.CreateIdemKey(model.IdemKey{ScopeID: sc.ID, Key: "request", Owner: "caller"})
	if err != nil {
		t.Logf("setup error: %v", err)
		t.FailNow()
	}
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	out := make(chan struct {
		rec *model.RequestRecord
		err error
	}, 1)
	go func() {
		rec, err := svc.ProcessRequest(ctx, ik.ID, "hash", func(workCtx context.Context) (string, error) {
			close(started)
			select {
			case <-workCtx.Done():
				return "", workCtx.Err()
			case <-time.After(time.Second):
				return "late", nil
			}
		})
		out <- struct {
			rec *model.RequestRecord
			err error
		}{rec, err}
	}()
	select {
	case <-started:
		cancel()
	case <-time.After(time.Second):
		t.Fatal("work was not invoked")
	}
	var outcome struct {
		rec *model.RequestRecord
		err error
	}
	select {
	case outcome = <-out:
	case <-time.After(time.Second):
		t.Fatal("request did not stop after cancellation")
	}
	rec, err := outcome.rec, outcome.err
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want cancellation, got %v", err)
	}
	if rec != nil && rec.Status != model.RequestRecordStatusFailed {
		t.Fatalf("want failed record, got %#v", rec)
	}
}
