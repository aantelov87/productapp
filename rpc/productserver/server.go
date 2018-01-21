package productserver

import (
	"net/http"
	"time"

	pb "github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"

	app "github.com/productapp"
	"github.com/productapp/proto"
	"github.com/productapp/rpc"
)

type server struct {
	storage app.ProductStore
}

func New(store app.ProductStore) http.Handler {
	s := server{
		storage: store,
	}
	return rpc.NewServer(rpc.Service{
		Name: "Product",
		Methods: map[string]rpc.Method{
			"List":   s.List,
			"Create": s.Create,
			"Lookup": s.Lookup,
		},
	})
}

func (s *server) List(_ []byte) (pb.Message, error) {
	pl, err := s.storage.List()
	if err != nil {
		return nil, err
	}
	var res proto.ListProductsResponse
	res.Products = proto.Products(pl)
	return &res, nil
}

func (s *server) Create(reqBytes []byte) (pb.Message, error) {
	var req proto.CreateProductRequest
	if err := pb.Unmarshal(reqBytes, &req); err != nil {
		return nil, err
	}
	var p app.Product
	proto.AppProduct(&p, req.Product)
	id := uuid.NewV4()
	p.ID = string(id.Bytes())
	p.Created = time.Now()
	err := s.storage.Create(&p)
	if err != nil {
		return nil, err
	}
	return proto.ProductProto(&p), nil
}

func (s *server) Lookup(reqBytes []byte) (pb.Message, error) {
	var req proto.GetProductRequest
	if err := pb.Unmarshal(reqBytes, &req); err != nil {
		return nil, err
	}
	p, err := s.storage.Lookup(req.GetId())
	return proto.ProductProto(p), err
}
