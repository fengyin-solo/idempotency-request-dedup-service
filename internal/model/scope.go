package model

import (
	"strings"
	"time"
)

type Scope struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Scope) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	if s.Name == "" {
		return NewValidationError("name", "作用域名称不能为空")
	}
	return nil
}

type ScopeFilter struct {
	Keyword string
}

func (f ScopeFilter) Match(s *Scope) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Description), k) {
			return false
		}
	}
	return true
}
