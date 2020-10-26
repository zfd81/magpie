package server

import (
	"context"

	"github.com/zfd81/magpie/rpc/grpc"

	"github.com/zfd81/magpie/meta"
)

type Server struct{}

func (s *Server) CreateTable(ctx context.Context, tbl *grpc.Table) (*grpc.Response, error) {
	info := meta.TableInfo{
		Name:    tbl.Name,
		Text:    tbl.Text,
		Comment: tbl.Comment,
	}
	for _, v := range tbl.Columns {
		info.CreateColumn(v.Name, v.DataType.String())
	}
	info.Keys = tbl.Keys
	info.Indexes = tbl.Indexes
	for _, v := range tbl.DerivedCols {
		info.CreateDerivedColumn(v.Name, v.Expression)
	}
	db.CreateTable(info)
	return &grpc.Response{
		Code:    200,
		Message: "c",
	}, nil
}
