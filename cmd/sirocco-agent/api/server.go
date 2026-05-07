package api

import (
	"net/http"
	"sirocco-agent/config"
)

type Server struct {
	cfg config.Config
	mux *http.ServeMux
}

func NewServer(cfg config.Config) *Server {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}

	s.routes()
	return s
}

func (s *Server) Start() {
	http.ListenAndServe(":"+s.cfg.Port, s.mux)
}