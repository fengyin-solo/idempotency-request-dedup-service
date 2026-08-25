package service

import "time"

// EnqueueRequestHashes hands request hashes to a delayed indexer.
func (s *Service) EnqueueRequestHashes(hashes []string) <-chan []string {
	out := make(chan []string, 1)
	go func() {
		time.Sleep(10 * time.Millisecond)
		out <- hashes
	}()
	return out
}
