package health

import (
	"fmt"
	"net/http"
	"time"
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
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	if !s.dbOK {
		status = "degraded"
	}
	body := fmt.Sprintf(`{"status":"%s","db":"%s"}`, status, boolStr(s.dbOK))
	_, _ = w.Write([]byte(body))
}

func boolStr(b bool) string {
	if b {
		return "ok"
	}
	return "fail"
}
