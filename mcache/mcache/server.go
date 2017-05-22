package mcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

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

	http.HandleFunc("/cmd", s.commandHandler)

	go s.run()
}

//
func (s *Server) run() {
	if err := http.ListenAndServe(s.srvAddr, nil); err != nil {
		s.srvErr <- err
	}
}

// GetSerErrChan -
func (s *Server) GetSerErrChan() <-chan error {
	return s.srvErr
}

//
func (s *Server) commandHandler(w http.ResponseWriter, r *http.Request) {

	json, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("setHandler read boby err ", err)
		return
	}
	fmt.Println(string(json))
}
