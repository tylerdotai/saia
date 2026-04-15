package health

import (
	"fmt"
	"net/http"
)

type Server struct {
	port int
	dbOK bool
}

func NewServer(port int) *Server {
	return &Server{port: port}
}

func (s *Server) SetDBOK(ok bool) {
	s.dbOK = ok
}

func (s *Server) Start() error {
	http.HandleFunc("/health", s.handleHealth)
	http.HandleFunc("/health/ready", s.handleReady)
	addr := fmt.Sprintf(":%d", s.port)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	if !s.dbOK {
		status = "degraded"
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"%s","db":"%s"}`, status, boolStr(s.dbOK))))
}

func boolStr(b bool) string {
	if b {
		return "ok"
	}
	return "fail"
}
