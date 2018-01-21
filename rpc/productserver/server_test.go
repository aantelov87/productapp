package productserver_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	pb "github.com/golang/protobuf/proto"
	app "github.com/productapp"

	"github.com/productapp/proto"
	"github.com/productapp/rpc/productserver"
)

var _ app.ProductStore = &MockStorage{}

type MockStorage struct {
	OnCreate func(p *app.Product) error
	OnList   func() ([]*app.Product, error)
	OnLookup func(id string) (*app.Product, error)
}

func (ms *MockStorage) Create(p *app.Product) error {
	return ms.OnCreate(p)
}

func (ms *MockStorage) Lookup(id string) (*app.Product, error) {
	return ms.OnLookup(id)
}

func (ms *MockStorage) List() ([]*app.Product, error) {
	return ms.OnList()
}

func TestCreateProduct(t *testing.T) {
	storage := MockStorage{
		OnCreate: func(p *app.Product) error {
			var t time.Time
			p.ID = "1"
			p.Created = t
			if p.Name == "Product2" {
				return errors.New("could not create a new product with this name")
			}
			return nil
		},
	}
	tc := []struct {
		reqProto     pb.Message
		expectedRes  app.Product
		expectedCode int
		expectedErr  string
	}{
		{
			&proto.CreateProductRequest{
				Product: &proto.Product{
					Name: pb.String("Product1"),
				},
			},
			app.Product{
				ID:   "1",
				Name: "Product1",
			},
			http.StatusOK,
			"",
		},
		{
			&proto.CreateProductRequest{
				Product: &proto.Product{
					Name: pb.String("Product2"),
				},
			},
			app.Product{},
			http.StatusInternalServerError,
			"could not create a new product with this name",
		},
	}
	h := productserver.New(&storage)
	for _, c := range tc {
		reqBytes, _ := pb.Marshal(c.reqProto)
		buf := bytes.NewBuffer(reqBytes)
		req, _ := http.NewRequest("POST", "/api/Product/Create", buf)
		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		rr.Result()
		if rr.Code != c.expectedCode {
			t.Errorf("wrong http status code: got %+v, want %+v", rr.Code, c.expectedCode)
		}
		if c.expectedErr != "" && !strings.Contains(rr.Body.String(), c.expectedErr) {
			t.Errorf("wrong body error description: got %s, want %s", rr.Body.String(), c.expectedErr)
		}
		if rr.Code == http.StatusOK {
			var p proto.Product
			if err := pb.Unmarshal(rr.Body.Bytes(), &p); err != nil {
				t.Fatal("wrong protobuf send by handler expected a proto.Product")
			}
			var p1 app.Product
			proto.AppProduct(&p1, &p)
			if !reflect.DeepEqual(p1, c.expectedRes) {
				t.Errorf("wrong create product response: got %+v, want %+v", p1, c.expectedRes)
			}
		}
	}
}
