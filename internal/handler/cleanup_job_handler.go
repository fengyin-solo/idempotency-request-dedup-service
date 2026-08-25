package handler

import (
	"net/http"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerCleanupJobRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/cleanupjobs", s.createCleanupJob)
	mux.HandleFunc("GET /api/cleanupjobs", s.listCleanupJobs)
	mux.HandleFunc("GET /api/cleanupjobs/{id}", s.getCleanupJob)
	mux.HandleFunc("PUT /api/cleanupjobs/{id}", s.updateCleanupJob)
	mux.HandleFunc("DELETE /api/cleanupjobs/{id}", s.deleteCleanupJob)
}

type createCleanupJobRequest struct {
	PolicyID string `json:"policy_id"`
	Message  string `json:"message"`
}

type updateCleanupJobRequest struct {
	DeletedCount int    `json:"deleted_count"`
	Status       string `json:"status"`
	Message      string `json:"message"`
}

func (s *Server) createCleanupJob(w http.ResponseWriter, r *http.Request) {
	var req createCleanupJobRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateCleanupJob(model.CleanupJob{PolicyID: req.PolicyID, Message: req.Message})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listCleanupJobs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.CleanupJobFilter{
		PolicyID: r.URL.Query().Get("policy_id"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListCleanupJobs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getCleanupJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := s.svc.GetCleanupJob(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateCleanupJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateCleanupJobRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateCleanupJob(id, model.CleanupJob{
		DeletedCount: req.DeletedCount,
		Status:       req.Status,
		Message:      req.Message,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteCleanupJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteCleanupJob(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
