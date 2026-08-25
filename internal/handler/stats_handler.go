package handler

import (
	"net/http"

	"idem/internal/model"
	"idem/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	idemKeys, _, err := s.svc.ListIdemKeys(model.IdemKeyFilter{}, 1, 100000)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	idemStatus := map[string]int{}
	for _, ik := range idemKeys {
		idemStatus[ik.Status]++
	}

	dedupHits, _, err := s.svc.ListDedupHits(model.DedupHitFilter{}, 1, 100000)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	dedupTotal := len(dedupHits)

	cleanupJobs, _, err := s.svc.ListCleanupJobs(model.CleanupJobFilter{}, 1, 100000)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	cleanupStatus := map[string]int{}
	for _, cj := range cleanupJobs {
		cleanupStatus[cj.Status]++
	}

	scopes, _, err := s.svc.ListScopes(model.ScopeFilter{}, 1, 100000)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	scopeKeyCount := map[string]int{}
	for _, sc := range scopes {
		scopeKeyCount[sc.Name] = 0
	}
	for _, ik := range idemKeys {
		for _, sc := range scopes {
			if ik.ScopeID == sc.ID {
				scopeKeyCount[sc.Name]++
				break
			}
		}
	}

	httpx.OK(w, map[string]interface{}{
		"idem_key_status":    idemStatus,
		"dedup_hit_total":    dedupTotal,
		"cleanup_job_status": cleanupStatus,
		"scope_key_count":    scopeKeyCount,
	})
}
