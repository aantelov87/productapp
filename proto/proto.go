//go:generate protoc product.proto --go_out=.
package proto

import (
	"github.com/golang/protobuf/ptypes"

	app "github.com/productapp"
)

// AppProduct converts a proto.Product struct to a app.Product.
func AppProduct(dst *app.Product, src *Product) {
	if dst == nil || src == nil {
		return
	}
	if src.Id != nil {
		dst.ID = src.GetId()
	}
	if src.Name != nil {
		dst.Name = src.GetName()
	}
	if src.Description != nil {
		dst.Desc = src.GetDescription()
	}
	if src.Attributes != nil {
		dst.Attributes = src.GetAttributes()
	}
	if src.Images != nil {
		dst.Images = src.GetImages()
	}
	if src.CreateTime != nil {
		t, _ := ptypes.Timestamp(src.CreateTime)
		dst.Created = t
	}
}

// ProductProto converts an app.Product struct to a proto.Product.
func ProductProto(p *app.Product) *Product {
	if p == nil {
		return nil
	}
	t, _ := ptypes.TimestampProto(p.Created)
	return &Product{
		Id:          &p.ID,
		Name:        &p.Name,
		Description: &p.Desc,
		Attributes:  p.Attributes,
		Images:      p.Images,
		CreateTime:  t,
	}
}

// Products converts from slices of app.Product to proto's.
func Products(pl []*app.Product) []*Product {
	if len(pl) == 0 {
		return nil
	}
	var rpl = make([]*Product, len(pl))
	for i := range rpl {
		rpl[i] = ProductProto(pl[i])
	}
	return rpl
}
