// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"idem/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Scope
	CreateScope(s *model.Scope) error
	GetScope(id string) (*model.Scope, error)
	ListScopes() []*model.Scope
	UpdateScope(s *model.Scope) error
	DeleteScope(id string) error

	// IdemKey
	CreateIdemKey(i *model.IdemKey) error
	GetIdemKey(id string) (*model.IdemKey, error)
	GetIdemKeyByScopeAndKey(scopeID, key string) (*model.IdemKey, error)
	ListIdemKeys() []*model.IdemKey
	UpdateIdemKey(i *model.IdemKey) error
	DeleteIdemKey(id string) error
	DeleteIdemKeysByIDs(ids []string) error

	// RequestRecord
	CreateRequestRecord(r *model.RequestRecord) error
	GetRequestRecord(id string) (*model.RequestRecord, error)
	ListRequestRecords() []*model.RequestRecord
	UpdateRequestRecord(r *model.RequestRecord) error
	DeleteRequestRecord(id string) error

	// DedupHit
	CreateDedupHit(d *model.DedupHit) error
	GetDedupHit(id string) (*model.DedupHit, error)
	ListDedupHits() []*model.DedupHit
	DeleteDedupHit(id string) error
	DeleteDedupHitsByIDs(ids []string) error

	// TTLPolicy
	CreateTTLPolicy(t *model.TTLPolicy) error
	GetTTLPolicy(id string) (*model.TTLPolicy, error)
	ListTTLPolicies() []*model.TTLPolicy
	UpdateTTLPolicy(t *model.TTLPolicy) error
	DeleteTTLPolicy(id string) error

	// CleanupJob
	CreateCleanupJob(c *model.CleanupJob) error
	GetCleanupJob(id string) (*model.CleanupJob, error)
	ListCleanupJobs() []*model.CleanupJob
	UpdateCleanupJob(c *model.CleanupJob) error
	DeleteCleanupJob(id string) error
}
