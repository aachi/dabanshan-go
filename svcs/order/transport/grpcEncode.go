package transport

import (
	"context"
	"errors"

	"github.com/laidingqing/dabanshan/pb"
	"github.com/laidingqing/dabanshan/svcs/order/model"
)

// CreateOrder encode/decode
func decodeGRPCCreateOrderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateOrderRequest)
	return model.CreateOrderRequest{
		Amount: float32(req.Amount),
	}, nil
}

func encodeGRPCCreateOrderResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.CreatedOrderResponse)
	return &pb.CreatedOrderResponse{Err: err2str(resp.Err)}, nil
}

// GetOrders encode/decode

func decodeGRPCGetOrdersRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetOrdersRequest)
	return model.GetOrdersRequest{
		UserID: req.Userid,
	}, nil
}

func encodeGRPCGetOrdersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.GetOrdersResponse)
	return &pb.GetOrdersResponse{
		Err: err2str(resp.Err),
	}, nil
}

// addCart encode/decode func

func encodeGRPCAddCartResponse(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.CreateCartRequest)
	return &pb.CreateCartRequest{
		Item: &pb.OrderItemRecord{
			Productid: req.ProductID,
			Userid:    req.UserID,
			Price:     float32(req.Price),
		},
	}, nil
}

func decodeGRPCAddCartRequest(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreatedCartResponse)
	return model.CreatedCartResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err)}, nil
}

// client encode and decode

func encodeGRPCCreateOrderRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.CreateOrderRequest)
	return &pb.CreateOrderRequest{
		Amount: float32(req.Amount),
	}, nil
}

func decodeGRPCCreateOrderResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreatedOrderResponse)
	return model.CreatedOrderResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err)}, nil
}

// getOrders encode/decode func

func encodeGRPCGetOrdersRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.GetOrdersRequest)
	return &pb.GetOrdersRequest{
		Userid: req.UserID,
	}, nil
}

func decodeGRPCGetOrdersResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetOrdersResponse)
	return model.GetOrdersResponse{
		Err: str2err(reply.Err)}, nil
}

func encodeGRPCAddCartRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.CreateCartRequest)
	return &pb.CreateCartRequest{
		Item: &pb.OrderItemRecord{
			Price:     req.Price,
			Productid: req.ProductID,
			Userid:    req.UserID,
		},
	}, nil
}

func decodeGRPCAddCartResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.CreatedCartResponse)
	return model.CreatedCartResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err)}, nil
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
