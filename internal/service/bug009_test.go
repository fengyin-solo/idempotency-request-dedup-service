package service

import (
	"sync"
	"testing"
	"time"

	"idem/internal/config"
	"idem/internal/model"
	"idem/internal/store"
	"idem/pkg/logger"
)

func TestConcurrentRequestRecordUpdatesAreRaceFree(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{MaxPageSize: 20})
	sc, _ := svc.CreateScope(model.Scope{Name: "race-scope"})
	ik, _ := svc.CreateIdemKey(model.IdemKey{ScopeID: sc.ID, Key: "race", Owner: "caller"})
	rec, err := svc.CreateRequestRecord(model.RequestRecord{IdemKeyID: ik.ID, RequestHash: "initial"})
	if err != nil {
		t.Logf("setup error: %v", err)
		t.FailNow()
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for _, hash := range []string{"left", "right"} {
		wg.Add(1)
		go func(h string) {
			defer wg.Done()
			<-start
			_, _ = svc.UpdateRequestRecord(rec.ID, model.RequestRecord{RequestHash: h, Status: model.RequestRecordStatusSuccess})
		}(hash)
	}
	close(start)
	wg.Wait()
	time.Sleep(time.Millisecond)
	if rec.RequestHash != "initial" || rec.Status != model.RequestRecordStatusPending {
		t.Fatalf("caller snapshot changed during concurrent updates: %#v", rec)
	}
	got, err := svc.GetRequestRecord(rec.ID)
	if err != nil {
		t.Logf("setup error: %v", err)
		t.FailNow()
	}
	if got.Status != model.RequestRecordStatusSuccess {
		t.Fatalf("unexpected final status: %#v", got)
	}
}
