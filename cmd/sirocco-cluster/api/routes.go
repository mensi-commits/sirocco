package api

import "net/http"

func (s *Server) routes() {
	s.mux.HandleFunc("/register", registerNode)
	s.mux.HandleFunc("/heartbeat", heartbeat)
	s.mux.HandleFunc("/shard/create", createShard)
}