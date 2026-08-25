package handler

import (
	"net/http"
	"strconv"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerTTLPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/ttlpolicies", s.createTTLPolicy)
	mux.HandleFunc("GET /api/ttlpolicies", s.listTTLPolicies)
	mux.HandleFunc("GET /api/ttlpolicies/{id}", s.getTTLPolicy)
	mux.HandleFunc("PUT /api/ttlpolicies/{id}", s.updateTTLPolicy)
	mux.HandleFunc("DELETE /api/ttlpolicies/{id}", s.deleteTTLPolicy)
}

type createTTLPolicyRequest struct {
	ScopeID     string `json:"scope_id"`
	TTLSeconds  int    `json:"ttl_seconds"`
	MaxEntries  int    `json:"max_entries"`
	AutoCleanup bool   `json:"auto_cleanup"`
}

func (s *Server) createTTLPolicy(w http.ResponseWriter, r *http.Request) {
	var req createTTLPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTTLPolicy(model.TTLPolicy{
		ScopeID:     req.ScopeID,
		TTLSeconds:  req.TTLSeconds,
		MaxEntries:  req.MaxEntries,
		AutoCleanup: req.AutoCleanup,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTTLPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TTLPolicyFilter{
		ScopeID: r.URL.Query().Get("scope_id"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	if v := r.URL.Query().Get("auto_cleanup"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.AutoCleanup = &b
	}
	items, total, err := s.svc.ListTTLPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTTLPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetTTLPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) updateTTLPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createTTLPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTTLPolicy(id, model.TTLPolicy{
		TTLSeconds:  req.TTLSeconds,
		MaxEntries:  req.MaxEntries,
		AutoCleanup: req.AutoCleanup,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTTLPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteTTLPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
