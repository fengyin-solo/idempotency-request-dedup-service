package model

import (
	"strings"
	"time"
)

const (
	IdemKeyStatusProcessing = "processing"
	IdemKeyStatusCompleted  = "completed"
	IdemKeyStatusExpired    = "expired"
)

var idemKeyTransitions = map[string]map[string]bool{
	IdemKeyStatusProcessing: {IdemKeyStatusCompleted: true, IdemKeyStatusExpired: true},
}

func IdemKeyCanTransition(from, to string) bool {
	if m, ok := idemKeyTransitions[from]; ok {
		return m[to]
	}
	return false
}

type IdemKey struct {
	ID          string    `json:"id"`
	ScopeID     string    `json:"scope_id"`
	Key         string    `json:"key"`
	Owner       string    `json:"owner"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	ExpireAt    time.Time `json:"expire_at"`
	CompletedAt time.Time `json:"completed_at"`
}

func (i *IdemKey) Validate() error {
	i.Key = strings.TrimSpace(i.Key)
	i.Owner = strings.TrimSpace(i.Owner)
	if i.ScopeID == "" {
		return NewValidationError("scope_id", "作用域 ID 不能为空")
	}
	if i.Key == "" {
		return NewValidationError("key", "幂等键不能为空")
	}
	if i.Owner == "" {
		return NewValidationError("owner", "所有者不能为空")
	}
	if i.Status == "" {
		i.Status = IdemKeyStatusProcessing
	}
	if i.Status != IdemKeyStatusProcessing && i.Status != IdemKeyStatusCompleted && i.Status != IdemKeyStatusExpired {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type IdemKeyFilter struct {
	ScopeID string
	Status  string
	Keyword string
}

func (f IdemKeyFilter) Match(i *IdemKey) bool {
	if f.ScopeID != "" && i.ScopeID != f.ScopeID {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(i.Key), k) &&
			!strings.Contains(strings.ToLower(i.Owner), k) {
			return false
		}
	}
	return true
}
