package model

import (
	"strings"
	"time"
)

type DedupHit struct {
	ID                string    `json:"id"`
	IdemKeyID         string    `json:"idem_key_id"`
	OriginalRequestID string    `json:"original_request_id"`
	RequestHash       string    `json:"request_hash"`
	HitAt             time.Time `json:"hit_at"`
}

func (d *DedupHit) Validate() error {
	d.RequestHash = strings.TrimSpace(d.RequestHash)
	if d.IdemKeyID == "" {
		return NewValidationError("idem_key_id", "幂等键 ID 不能为空")
	}
	if d.OriginalRequestID == "" {
		return NewValidationError("original_request_id", "原始请求 ID 不能为空")
	}
	if d.RequestHash == "" {
		return NewValidationError("request_hash", "请求哈希不能为空")
	}
	return nil
}

type DedupHitFilter struct {
	IdemKeyID string
	Keyword   string
}

func (f DedupHitFilter) Match(d *DedupHit) bool {
	if f.IdemKeyID != "" && d.IdemKeyID != f.IdemKeyID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.RequestHash), k) &&
			!strings.Contains(strings.ToLower(d.OriginalRequestID), k) {
			return false
		}
	}
	return true
}
