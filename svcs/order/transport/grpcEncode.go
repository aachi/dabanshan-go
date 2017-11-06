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

func decodeGRPCAddCartRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateCartRequest)
	return model.CreateCartRequest{
		UserID:    req.Item.Userid,
		Price:     req.Item.Price,
		ProductID: req.Item.Productid,
	}, nil
}

func encodeGRPCAddCartResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.CreatedCartResponse)
	return &pb.CreatedCartResponse{
		Id:  resp.ID,
		Err: err2str(resp.Err),
	}, nil
}

// GetCartItems encode/decode

func decodeGRPCGetCartItemsRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetCartItemsRequest)
	return model.GetCartItemsRequest{
		UserID: req.Userid,
	}, nil
}

func encodeGRPCGetCartItemsResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.GetCartItemsResponse)
	return &pb.GetCartItemsResponse{
		Items: modelCartItem2Pb(resp.Items),
		Err:   err2str(resp.Err),
	}, nil
}

// removeCartItem
func decodeGRPCRemoveCartItemRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RemoveCartItemRequest)
	return model.RemoveCartItemRequest{
		CartID: req.Cartid,
	}, nil
}

func encodeGRPCRemoveCartItemResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(model.RemoveCartItemResponse)
	return &pb.RemoveCartItemResponse{
		Err: err2str(resp.Err),
	}, nil
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

func encodeGRPCCartItemsRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.GetCartItemsRequest)
	return &pb.GetCartItemsRequest{
		Userid: req.UserID,
	}, nil
}

func decodeGRPCCartItemsResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetCartItemsResponse)
	return model.GetCartItemsResponse{
		Items: pbCartItem2Model(reply.Items),
		Err:   str2err(reply.Err)}, nil
}

func encodeGRPCRemoveCartItemRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(model.RemoveCartItemRequest)
	return &pb.RemoveCartItemRequest{
		Cartid: req.CartID,
	}, nil
}

func decodeGRPCRemoveCartItemResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.RemoveCartItemResponse)
	return model.RemoveCartItemResponse{
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

func pbCartItem2Model(records []*pb.OrderItemRecord) []model.Cart {
	var models []model.Cart
	for _, record := range records {
		models = append(models, model.Cart{
			UserID:    record.Userid,
			Price:     record.Price,
			ProductID: record.Productid,
			CartID:    record.Cartid,
		})
	}
	return models
}

func modelCartItem2Pb(models []model.Cart) []*pb.OrderItemRecord {
	var records []*pb.OrderItemRecord
	for _, model := range models {
		records = append(records, &pb.OrderItemRecord{
			Price:     model.Price,
			Productid: model.ProductID,
			Userid:    model.UserID,
			Cartid:    model.CartID,
		})
	}

	return records
}
