package model

import (
	"strings"
	"time"
)

const (
	RequestRecordStatusPending = "pending"
	RequestRecordStatusSuccess = "success"
	RequestRecordStatusFailed  = "failed"
)

type RequestRecord struct {
	ID           string    `json:"id"`
	IdemKeyID    string    `json:"idem_key_id"`
	RequestHash  string    `json:"request_hash"`
	ResponseHash string    `json:"response_hash"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func (r *RequestRecord) Validate() error {
	r.RequestHash = strings.TrimSpace(r.RequestHash)
	r.ResponseHash = strings.TrimSpace(r.ResponseHash)
	if r.IdemKeyID == "" {
		return NewValidationError("idem_key_id", "幂等键 ID 不能为空")
	}
	if r.RequestHash == "" {
		return NewValidationError("request_hash", "请求哈希不能为空")
	}
	if r.Status == "" {
		r.Status = RequestRecordStatusPending
	}
	if r.Status != RequestRecordStatusPending && r.Status != RequestRecordStatusSuccess && r.Status != RequestRecordStatusFailed {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type RequestRecordFilter struct {
	IdemKeyID   string
	Status      string
	RequestHash string
}

func (f RequestRecordFilter) Match(r *RequestRecord) bool {
	if f.IdemKeyID != "" && r.IdemKeyID != f.IdemKeyID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.RequestHash != "" && r.RequestHash != f.RequestHash {
		return false
	}
	return true
}
