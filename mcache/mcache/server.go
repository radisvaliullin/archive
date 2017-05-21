package mcache

import "net/http"

// Server - implements memory cache server with REST Api for clients
type Server struct {
	store *Storage

	srvAddr string
	srvErr  chan error
}

// NewMCacheServer -
func NewMCacheServer(addr string) *Server {
	s := &Server{
		store:   NewStorage(),
		srvAddr: addr,
		srvErr:  make(chan error, 100),
	}
	return s
}

// Start - start server
func (s *Server) Start() {

	http.HandleFunc("/set", s.setHandler)

	go s.run()
}

//
func (s *Server) run() {
	if err := http.ListenAndServe(":8080", nil); err != nil {
		s.srvErr <- err
	}
}

// GetSerErrChan -
func (s *Server) GetSerErrChan() <-chan error {
	return s.srvErr
}

//
func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {

}
