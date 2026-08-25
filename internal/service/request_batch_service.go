package service

import (
	"time"
)

// EnqueueRequestHashes hands request hashes to a delayed indexer.
//
// hashes 是调用方持有的切片，其底层可能被后续请求原地改写或追加复用，
// 因此必须在入队时立即快照拷贝一份独立副本交给后台，避免延迟消费期间
// 跨请求的切片修改污染后台最终拿到的内容。
func (s *Service) EnqueueRequestHashes(hashes []string) <-chan []string {
	snapshot := make([]string, len(hashes))
	copy(snapshot, hashes)
	out := make(chan []string, 1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		out <- snapshot
	}()
	return out
}
