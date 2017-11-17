package transport

import (
	"context"
	"errors"
	"fmt"

	"github.com/laidingqing/dabanshan/pb"
	"github.com/laidingqing/dabanshan/svcs/product/model"
	"github.com/laidingqing/dabanshan/utils"
)

// server

// get products encode/decode
func decodeGRPCGetProductsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetProductsRequest)
	return model.GetProductsRequest{A: int64(req.Creatorid), B: int64(req.Size)}, nil
}

func encodeGRPCGetProductsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.GetProductsResponse)
	return &pb.GetProductsResponse{V: int64(resp.V), Err: err2str(resp.Err)}, nil
}

// Upload ...
func decodeGRPCUploadRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ProductUploadRequest)
	return model.UploadProductRequest{
		Name: req.Name,
		Md5:  req.Md5,
		Body: req.B,
	}, nil
}

func encodeGRPCUploadResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.UploadProductResponse)
	return &pb.ProductUploadResponse{Name: resp.ID}, nil
}

// create products encode/decode
func decodeGRPCCreateProductRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateProductRequest)
	logger := utils.NewLogger()
	logger.Log("create name", req.Name)
	fmt.Println("create name fmt", req.Name)
	return model.CreateProductRequest{
		Product: model.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			UserID:      req.UserID,
			CatalogID:   req.CatalogID,
			Status:      req.Status,
			Thumbnails:  req.Thumbnails,
		},
	}, nil
}
func encodeGRPCCreateProductResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.CreateProductResponse)
	return &pb.CreateProductResponse{
		Id:  resp.ID,
		Err: err2str(resp.Err),
	}, nil
}

// client

// create products encode/decode
func encodeGRPCCreateProductRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.CreateProductRequest)
	return &pb.CreateProductRequest{
		Name:        req.Product.Name,
		Description: req.Product.Description,
		Price:       req.Product.Price,
		UserID:      req.Product.UserID,
		CatalogID:   req.Product.CatalogID,
		Status:      req.Product.Status,
		Thumbnails:  req.Product.Thumbnails,
	}, nil
}

func decodeGRPCCreateProductResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreateProductResponse)
	return model.CreateProductResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err)}, nil
}

// get products encode/decode
func encodeGRPCGetProductsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.GetProductsRequest)
	return &pb.GetProductsRequest{Creatorid: int64(req.A), Size: int64(req.B)}, nil
}

func decodeGRPCGetProductsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetProductsResponse)
	return model.GetProductsResponse{V: int64(reply.V), Err: str2err(reply.Err)}, nil
}

// upload
func encodeGRPCUploadRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.UploadProductRequest)
	return &pb.ProductUploadRequest{
		Name: req.Name,
		B:    req.Body,
		Md5:  req.Md5,
	}, nil
}

func decodeGRPCUploadResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ProductUploadResponse)
	return model.UploadProductResponse{ID: reply.Name}, nil
}

func str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
