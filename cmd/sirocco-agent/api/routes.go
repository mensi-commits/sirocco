package api

import (
	"net/http"
	"sirocco-agent/node"
	"sirocco-agent/shard"
)

func (s *Server) routes() {

	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	s.mux.HandleFunc("/command", s.auth(shard.HandleCommand(s.cfg)))

	s.mux.HandleFunc("/node/register", node.RegisterHandler(s.cfg))
}