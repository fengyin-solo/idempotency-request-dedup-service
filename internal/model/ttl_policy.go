package model

import (
	"strings"
	"time"
)

type TTLPolicy struct {
	ID          string    `json:"id"`
	ScopeID     string    `json:"scope_id"`
	TTLSeconds  int       `json:"ttl_seconds"`
	MaxEntries  int       `json:"max_entries"`
	AutoCleanup bool      `json:"auto_cleanup"`
	CreatedAt   time.Time `json:"created_at"`
}

func (t *TTLPolicy) Validate() error {
	if t.ScopeID == "" {
		return NewValidationError("scope_id", "作用域 ID 不能为空")
	}
	if t.TTLSeconds <= 0 {
		return NewValidationError("ttl_seconds", "TTL 秒数必须大于 0")
	}
	if t.MaxEntries <= 0 {
		return NewValidationError("max_entries", "最大条目数必须大于 0")
	}
	return nil
}

type TTLPolicyFilter struct {
	ScopeID     string
	AutoCleanup *bool
	Keyword     string
}

func (f TTLPolicyFilter) Match(t *TTLPolicy) bool {
	if f.ScopeID != "" && t.ScopeID != f.ScopeID {
		return false
	}
	if f.AutoCleanup != nil && t.AutoCleanup != *f.AutoCleanup {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.ScopeID), k) {
			return false
		}
	}
	return true
}
