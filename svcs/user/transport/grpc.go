package transport

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	jujuratelimit "github.com/juju/ratelimit"
	"github.com/laidingqing/dabanshan/pb"
	u_endpoint "github.com/laidingqing/dabanshan/svcs/user/endpoint"
	m_user "github.com/laidingqing/dabanshan/svcs/user/model"
	"github.com/laidingqing/dabanshan/svcs/user/service"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	oldcontext "golang.org/x/net/context"
)

type grpcServer struct {
	getuser grpctransport.Handler
}

// NewGRPCServer ...
func NewGRPCServer(endpoints u_endpoint.Set, tracer stdopentracing.Tracer, logger log.Logger) pb.UserRpcServiceServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		getuser: grpctransport.NewServer(
			endpoints.GetUserEndpoint,
			decodeGRPCGetUserRequest,
			encodeGRPCGetUserResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "GetUser", logger)))...,
		),
	}
}

func (s *grpcServer) GetUser(ctx oldcontext.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	_, rep, err := s.getuser.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.GetUserResponse)
	return res, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetUserRequest)
	return m_user.GetUserRequest{A: req.Userid}, nil
}

func encodeGRPCGetUserResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(m_user.GetUserResponse)
	return &pb.GetUserResponse{V: &pb.UserRecord{
		Firstname: resp.V.FirstName,
		Lastname:  resp.V.LastName,
		Email:     resp.V.Email,
		Username:  resp.V.Username,
		Password:  "",
		Salt:      "",
		Userid:    resp.V.UserID,
	}, Err: err2str(resp.Err)}, nil
}

// NewGRPCClient ...
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {
	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))
	var getUserEndpoint endpoint.Endpoint
	{
		getUserEndpoint = grpctransport.NewClient(
			conn,
			"pb.UserRpcService",
			"GetUser",
			encodeGRPCGetUserRequest,
			decodeGRPCGetUserResponse,
			pb.GetUserResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		getUserEndpoint = opentracing.TraceClient(tracer, "GetUser")(getUserEndpoint)
		getUserEndpoint = limiter(getUserEndpoint)
		getUserEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetUser",
			Timeout: 30 * time.Second,
		}))(getUserEndpoint)
	}
	return u_endpoint.Set{
		GetUserEndpoint: getUserEndpoint,
	}
}

func encodeGRPCGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(m_user.GetUserRequest)
	return &pb.GetUserRequest{Userid: req.A}, nil
}

func decodeGRPCGetUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetUserResponse)
	return m_user.GetUserResponse{V: m_user.User{
		FirstName: reply.V.Firstname,
		LastName:  reply.V.Lastname,
		Email:     reply.V.Email,
		Username:  reply.V.Username,
		Password:  "",
		Salt:      "",
		UserID:    reply.V.Userid,
	}, Err: str2err(reply.Err)}, nil
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
