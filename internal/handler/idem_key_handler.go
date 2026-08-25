package handler

import (
	"net/http"
	"time"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerIdemKeyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/idemkeys", s.createIdemKey)
	mux.HandleFunc("GET /api/idemkeys", s.listIdemKeys)
	mux.HandleFunc("GET /api/idemkeys/{id}", s.getIdemKey)
	mux.HandleFunc("PUT /api/idemkeys/{id}", s.updateIdemKey)
	mux.HandleFunc("DELETE /api/idemkeys/{id}", s.deleteIdemKey)
	mux.HandleFunc("POST /api/idemkeys/batch-delete", s.batchDeleteIdemKeys)
	mux.HandleFunc("POST /api/idemkeys/expire", s.expireIdemKeys)
}

type createIdemKeyRequest struct {
	ScopeID  string    `json:"scope_id"`
	Key      string    `json:"key"`
	Owner    string    `json:"owner"`
	ExpireAt time.Time `json:"expire_at"`
}

type updateIdemKeyRequest struct {
	Key      string    `json:"key"`
	Owner    string    `json:"owner"`
	Status   string    `json:"status"`
	ExpireAt time.Time `json:"expire_at"`
}

func (s *Server) createIdemKey(w http.ResponseWriter, r *http.Request) {
	var req createIdemKeyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.CreateIdemKey(model.IdemKey{ScopeID: req.ScopeID, Key: req.Key, Owner: req.Owner, ExpireAt: req.ExpireAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, i)
}

func (s *Server) listIdemKeys(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.IdemKeyFilter{
		ScopeID: r.URL.Query().Get("scope_id"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListIdemKeys(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getIdemKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := s.svc.GetIdemKey(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

func (s *Server) updateIdemKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateIdemKeyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.UpdateIdemKey(id, model.IdemKey{Key: req.Key, Owner: req.Owner, Status: req.Status, ExpireAt: req.ExpireAt})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

func (s *Server) deleteIdemKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteIdemKey(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchDeleteRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteIdemKeys(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchDeleteIdemKeys(req.IDs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) expireIdemKeys(w http.ResponseWriter, r *http.Request) {
	count, err := s.svc.ExpireIdemKeys(time.Now())
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"expired_count": count})
}
