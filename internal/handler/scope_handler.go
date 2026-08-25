package handler

import (
	"net/http"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerScopeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/scopes", s.createScope)
	mux.HandleFunc("GET /api/scopes", s.listScopes)
	mux.HandleFunc("GET /api/scopes/{id}", s.getScope)
	mux.HandleFunc("PUT /api/scopes/{id}", s.updateScope)
	mux.HandleFunc("DELETE /api/scopes/{id}", s.deleteScope)
}

type createScopeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createScope(w http.ResponseWriter, r *http.Request) {
	var req createScopeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sc, err := s.svc.CreateScope(model.Scope{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sc)
}

func (s *Server) listScopes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ScopeFilter{Keyword: r.URL.Query().Get("keyword")}
	items, total, err := s.svc.ListScopes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getScope(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sc, err := s.svc.GetScope(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sc)
}

func (s *Server) updateScope(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createScopeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sc, err := s.svc.UpdateScope(id, model.Scope{Name: req.Name, Description: req.Description})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sc)
}

func (s *Server) deleteScope(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteScope(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
