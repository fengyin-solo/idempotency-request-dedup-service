package model

import (
	"strings"
	"time"
)

const (
	CleanupJobStatusRunning = "running"
	CleanupJobStatusDone    = "done"
	CleanupJobStatusFailed  = "failed"
)

var cleanupJobTransitions = map[string]map[string]bool{
	CleanupJobStatusRunning: {CleanupJobStatusDone: true, CleanupJobStatusFailed: true},
}

func CleanupJobCanTransition(from, to string) bool {
	if m, ok := cleanupJobTransitions[from]; ok {
		return m[to]
	}
	return false
}

type CleanupJob struct {
	ID           string    `json:"id"`
	PolicyID     string    `json:"policy_id"`
	DeletedCount int       `json:"deleted_count"`
	Status       string    `json:"status"`
	Message      string    `json:"message"`
	RanAt        time.Time `json:"ran_at"`
}

func (c *CleanupJob) Validate() error {
	c.Message = strings.TrimSpace(c.Message)
	if c.PolicyID == "" {
		return NewValidationError("policy_id", "策略 ID 不能为空")
	}
	if c.Status == "" {
		c.Status = CleanupJobStatusRunning
	}
	if c.Status != CleanupJobStatusRunning && c.Status != CleanupJobStatusDone && c.Status != CleanupJobStatusFailed {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type CleanupJobFilter struct {
	PolicyID string
	Status   string
	Keyword  string
}

func (f CleanupJobFilter) Match(c *CleanupJob) bool {
	if f.PolicyID != "" && c.PolicyID != f.PolicyID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Message), k) {
			return false
		}
	}
	return true
}
