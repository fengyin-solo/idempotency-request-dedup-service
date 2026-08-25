package handler

import (
	"net/http"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerDedupHitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/deduphits", s.createDedupHit)
	mux.HandleFunc("GET /api/deduphits", s.listDedupHits)
	mux.HandleFunc("GET /api/deduphits/{id}", s.getDedupHit)
	mux.HandleFunc("DELETE /api/deduphits/{id}", s.deleteDedupHit)
	mux.HandleFunc("POST /api/deduphits/batch-delete", s.batchDeleteDedupHits)
}

type createDedupHitRequest struct {
	IdemKeyID         string `json:"idem_key_id"`
	OriginalRequestID string `json:"original_request_id"`
	RequestHash       string `json:"request_hash"`
}

func (s *Server) createDedupHit(w http.ResponseWriter, r *http.Request) {
	var req createDedupHitRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	d, err := s.svc.CreateDedupHit(model.DedupHit{
		IdemKeyID:         req.IdemKeyID,
		OriginalRequestID: req.OriginalRequestID,
		RequestHash:       req.RequestHash,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, d)
}

func (s *Server) listDedupHits(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.DedupHitFilter{
		IdemKeyID: r.URL.Query().Get("idem_key_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListDedupHits(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getDedupHit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	d, err := s.svc.GetDedupHit(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, d)
}

func (s *Server) deleteDedupHit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteDedupHit(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) batchDeleteDedupHits(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if err := s.svc.BatchDeleteDedupHits(req.IDs); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
