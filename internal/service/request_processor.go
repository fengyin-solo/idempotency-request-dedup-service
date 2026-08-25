package service

import (
	"context"
	"time"

	"idem/internal/model"
	"idem/pkg/idgen"
)

// ProcessRequest records one request while the caller performs its work.
func (s *Service) ProcessRequest(ctx context.Context, idemKeyID, requestHash string, work func(context.Context) (string, error)) (*model.RequestRecord, error) {
	if _, err := s.store.GetIdemKey(idemKeyID); err != nil {
		return nil, err
	}
	rec := &model.RequestRecord{ID: idgen.Hex(), IdemKeyID: idemKeyID, RequestHash: requestHash, Status: model.RequestRecordStatusPending, CreatedAt: time.Now()}
	if err := rec.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateRequestRecord(rec); err != nil {
		return nil, err
	}
	response, err := work(context.Background())
	if err != nil {
		rec.Status = model.RequestRecordStatusFailed
		_ = s.store.UpdateRequestRecord(rec)
		return rec, err
	}
	rec.ResponseHash = response
	rec.Status = model.RequestRecordStatusSuccess
	if err := s.store.UpdateRequestRecord(rec); err != nil {
		return nil, err
	}
	return rec, nil
}
