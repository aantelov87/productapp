package rpc

import (
	"io/ioutil"
	"net/http"
	"strings"

	pb "github.com/golang/protobuf/proto"
)

// Service describes an RPC service.
type Service struct {
	// The name of the service, which forms the first path component of any
	// HTTP request.
	Name string

	// The RPC methods to serve.
	Methods map[string]Method
}

// Method describes an RPC method.
type Method func(reqBytes []byte) (pb.Message, error)

// NewServer returns a new Server that uses the given ServerConfig.
func NewServer(svc Service) http.Handler {
	// Validate Service.
	if svc.Name == "" {
		panic("ServerConfig provided with empty Name")
	}

	return &serverImpl{
		service: svc,
	}
}

type serverImpl struct {
	service Service
}

// ServeHTTP exposes the configured Service as an HTTP API.
func (s *serverImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d := &s.service
	prefix := "/api/" + d.Name + "/"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		http.NotFound(w, r)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, prefix)
	method, ok := d.Methods[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := method(body)
	sendResponse(w, resp, err)
}

func sendResponse(w http.ResponseWriter, resp pb.Message, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	payload, err := pb.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h := w.Header()
	h.Set("Content-type", "application/octet-stream")
	w.Write(payload)
}
