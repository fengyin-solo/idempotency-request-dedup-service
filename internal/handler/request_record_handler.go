package handler

import (
	"net/http"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerRequestRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/requestrecords", s.createRequestRecord)
	mux.HandleFunc("GET /api/requestrecords", s.listRequestRecords)
	mux.HandleFunc("GET /api/requestrecords/{id}", s.getRequestRecord)
	mux.HandleFunc("PUT /api/requestrecords/{id}", s.updateRequestRecord)
	mux.HandleFunc("DELETE /api/requestrecords/{id}", s.deleteRequestRecord)
}

type createRequestRecordRequest struct {
	IdemKeyID    string `json:"idem_key_id"`
	RequestHash  string `json:"request_hash"`
	ResponseHash string `json:"response_hash"`
	Status       string `json:"status"`
}

func (s *Server) createRequestRecord(w http.ResponseWriter, r *http.Request) {
	var req createRequestRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateRequestRecord(model.RequestRecord{
		IdemKeyID:    req.IdemKeyID,
		RequestHash:  req.RequestHash,
		ResponseHash: req.ResponseHash,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listRequestRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RequestRecordFilter{
		IdemKeyID:   r.URL.Query().Get("idem_key_id"),
		Status:      r.URL.Query().Get("status"),
		RequestHash: r.URL.Query().Get("request_hash"),
	}
	items, total, err := s.svc.ListRequestRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRequestRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rec, err := s.svc.GetRequestRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) updateRequestRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req createRequestRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.UpdateRequestRecord(id, model.RequestRecord{
		RequestHash:  req.RequestHash,
		ResponseHash: req.ResponseHash,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) deleteRequestRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRequestRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
