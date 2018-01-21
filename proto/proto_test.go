package proto_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	app "github.com/productapp"
	pb "github.com/productapp/proto"
)

func TestAppProduct(t *testing.T) {
	tc := []struct {
		dest     *app.Product
		src      *pb.Product
		expected *app.Product
	}{
		{nil, nil, nil},
		{
			&app.Product{},
			&pb.Product{
				Name:        proto.String("Product1"),
				Description: proto.String("Description1"),
				Attributes:  []string{"attr1", "attr2"},
				Images:      []string{"images/image1.png", "images/image2.png"},
			},
			&app.Product{
				Name:       "Product1",
				Desc:       "Description1",
				Attributes: []string{"attr1", "attr2"},
				Images:     []string{"images/image1.png", "images/image2.png"},
			},
		},
		{
			&app.Product{
				Name: "Product1",
			},
			&pb.Product{
				Name:        proto.String("Product2"),
				Description: proto.String("Description1"),
				Attributes:  []string{"attr1", "attr2"},
				Images:      []string{"images/image1.png", "images/image2.png"},
			},
			&app.Product{
				Name:       "Product2",
				Desc:       "Description1",
				Attributes: []string{"attr1", "attr2"},
				Images:     []string{"images/image1.png", "images/image2.png"},
			},
		}, {
			&app.Product{
				Name: "Product3",
			},
			&pb.Product{
				Description: proto.String("Description1"),
				Attributes:  []string{"attr1", "attr2"},
				Images:      []string{"images/image1.png", "images/image2.png"},
			},
			&app.Product{
				Name:       "Product3",
				Desc:       "Description1",
				Attributes: []string{"attr1", "attr2"},
				Images:     []string{"images/image1.png", "images/image2.png"},
			},
		},
	}

	for _, c := range tc {
		pb.AppProduct(c.dest, c.src)
		if !reflect.DeepEqual(c.dest, c.expected) {
			t.Errorf("wrong product from proto.Product: got %+v, want %+v\n", c.dest, c.expected)
		}
	}
}

func TestProductProto(t *testing.T) {
	var tm time.Time
	tmpb, _ := ptypes.TimestampProto(tm)
	tc := []struct {
		in  *app.Product
		out *pb.Product
	}{
		{nil, nil},
		{
			&app.Product{
				Name:       "Product1",
				Desc:       "Description1",
				Attributes: []string{"attr1", "attr2"},
				Images:     []string{"images/image1.png", "images/image2.png"},
			},
			&pb.Product{
				Id:          proto.String(""),
				Name:        proto.String("Product1"),
				Description: proto.String("Description1"),
				Attributes:  []string{"attr1", "attr2"},
				Images:      []string{"images/image1.png", "images/image2.png"},
				CreateTime:  tmpb,
			},
		},
		{
			&app.Product{
				Name: "Product2",
			},
			&pb.Product{
				Id:          proto.String(""),
				Name:        proto.String("Product2"),
				Description: proto.String(""),
				CreateTime:  tmpb,
			},
		}, {
			&app.Product{},
			&pb.Product{
				Id:          proto.String(""),
				Name:        proto.String(""),
				Description: proto.String(""),
				CreateTime:  tmpb,
			},
		},
	}
	for _, c := range tc {
		got := pb.ProductProto(c.in)
		if !reflect.DeepEqual(got, c.out) {
			t.Errorf("wrong proto product from app.Product: got %+v, want %+v\n", got, c.out)
		}
	}

}
