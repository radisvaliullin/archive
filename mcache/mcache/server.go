package mcache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
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

	jsonCmd, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body error", 400)
		s.srvErr <- fmt.Errorf("command handler, body read err %v", err)
		return
	}

	fmt.Println(string(jsonCmd))

	cmd := &Command{}
	err = json.Unmarshal(jsonCmd, cmd)
	if err != nil {
		http.Error(w, "wrong command json", 400)
		s.srvErr <- fmt.Errorf("command handler, wrong command json %v", err)
		return
	}

	rr := s.commandExecut(cmd)
	respBody, err := json.Marshal(rr)
	if err != nil {
		http.Error(w, "response marshal error", 500)
		s.srvErr <- fmt.Errorf("command handler, response json marshal err %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}

//
func (s *Server) commandExecut(cmd *Command) *ReqResp {

	switch cmd.Name {
	case "set":
		switch *cmd.Type {
		case "str":
			fmt.Printf("SET STR - %+v\n", cmd)
			err := s.store.Set(*cmd.Key, *cmd.Str, time.Second*time.Duration(*cmd.TTL))
			if err != nil {
				return &ReqResp{Success: false, Error: PStr(err.Error())}
			}
			return &ReqResp{Success: true}

		case "slice":
			err := s.store.Set(*cmd.Key, cmd.Slice, time.Second*time.Duration(*cmd.TTL))
			if err != nil {
				return &ReqResp{Success: false, Error: PStr(err.Error())}
			}
			return &ReqResp{Success: true}

		case "map":
			err := s.store.Set(*cmd.Key, cmd.Map, time.Second*time.Duration(*cmd.TTL))
			if err != nil {
				return &ReqResp{Success: false, Error: PStr(err.Error())}
			}
			return &ReqResp{Success: true}

		default:
			return &ReqResp{Success: false, Error: PStr("unknown command")}
		}

	case "get":
		sv := s.store.Get(*cmd.Key)
		if sv == nil {
			return &ReqResp{Success: true}
		}

		// string
		str, ok := sv.GetString()
		fmt.Printf("GET STR %v %v %v %+v\n", str, ok, cmd.Key, cmd)
		if ok {
			return &ReqResp{Success: true, Result: PStr(str)}
		}

		// slice
		sl, ok := sv.GetSlice()
		if ok {
			if cmd.IdxKey != nil {
				idx, err := strconv.Atoi(*cmd.IdxKey)
				if err != nil {
					return &ReqResp{Success: false, Error: PStr("get slice by index, index must be number")}
				}
				item, _, err := sv.GetSliceItem(idx)
				if err != nil {
					return &ReqResp{Success: false, Error: PStr("get slice by index, err " + err.Error())}
				}
				return &ReqResp{Success: true, Result: PStr(item)}
			}
			return &ReqResp{Success: true, Result: PStr(fmt.Sprint(sl))}
		}

		// map
		m, ok := sv.GetMap()
		if ok {
			if cmd.IdxKey != nil {
				item, _, mok := sv.GetMapValByKey(*cmd.IdxKey)
				if !mok {
					return &ReqResp{Success: false, Error: PStr("get map by key, not exist")}
				}
				return &ReqResp{Success: true, Result: PStr(item)}
			}
			return &ReqResp{Success: true, Result: PStr(fmt.Sprint(m))}
		}

		return &ReqResp{Success: false, Error: PStr("get by key, unknown result")}

	case "remove":
		s.store.Remove(*cmd.Key)
		return &ReqResp{Success: true}

	case "keys":
		keys := s.store.Keys()
		return &ReqResp{Success: true, Result: PStr(fmt.Sprint(keys))}

	default:
		return &ReqResp{Success: false, Error: PStr("unknown command")}

	}

	return nil
}
