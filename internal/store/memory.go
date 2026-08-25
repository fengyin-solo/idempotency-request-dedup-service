package store

import (
	"sync"

	"idem/internal/model"
)

type MemoryStore struct {
	mu             sync.RWMutex
	scopes         map[string]*model.Scope
	idemKeys       map[string]*model.IdemKey
	requestRecords map[string]*model.RequestRecord
	dedupHits      map[string]*model.DedupHit
	ttlPolicies    map[string]*model.TTLPolicy
	cleanupJobs    map[string]*model.CleanupJob
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		scopes:         make(map[string]*model.Scope),
		idemKeys:       make(map[string]*model.IdemKey),
		requestRecords: make(map[string]*model.RequestRecord),
		dedupHits:      make(map[string]*model.DedupHit),
		ttlPolicies:    make(map[string]*model.TTLPolicy),
		cleanupJobs:    make(map[string]*model.CleanupJob),
	}
}

var _ Store = (*MemoryStore)(nil)
