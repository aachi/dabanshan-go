package transport

import (
	"context"
	"errors"
	"time"

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
	"google.golang.org/grpc"
)

type grpcServer struct {
	getuser  grpctransport.Handler
	register grpctransport.Handler
	login    grpctransport.Handler
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
		register: grpctransport.NewServer(
			endpoints.RegisterEndpoint,
			decodeGRPCRegisterRequest,
			encodeGRPCRegisterResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "Register", logger)))...,
		),
		login: grpctransport.NewServer(
			endpoints.LoginEndpoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
			append(options, grpctransport.ServerBefore(opentracing.GRPCToContext(tracer, "Login", logger)))...,
		),
	}
}

// GetUser RPC
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

// Register RPC
func (s *grpcServer) Register(ctx oldcontext.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	_, rep, err := s.register.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.RegisterResponse)
	return res, nil
}

func decodeGRPCRegisterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.RegisterRequest)
	return m_user.RegisterRequest{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		FirstName: req.Firstname,
		LastName:  req.Lastname,
	}, nil
}

func encodeGRPCRegisterResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(m_user.RegisterUserResponse)
	return &pb.RegisterResponse{
		Id:  resp.ID,
		Err: err2str(resp.Err),
	}, nil
}

func (s *grpcServer) Login(ctx oldcontext.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	res := rep.(*pb.LoginResponse)
	return res, nil
}

func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.LoginRequest)
	return m_user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func encodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(m_user.LoginResponse)
	return &pb.LoginResponse{
		User:  modelUser2PbUser(*resp.User),
		Token: resp.Token,
		Err:   err2str(resp.Err),
	}, nil
}

// NewGRPCClient ...
func NewGRPCClient(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) service.Service {
	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))
	var getUserEndpoint endpoint.Endpoint
	var registerEndpoint endpoint.Endpoint
	var loginEndPoint endpoint.Endpoint
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

		registerEndpoint = grpctransport.NewClient(
			conn,
			"pb.UserRpcService",
			"Register",
			encodeGRPCRegisterRequest,
			decodeGRPCRegisterResponse,
			pb.RegisterResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		registerEndpoint = opentracing.TraceClient(tracer, "Register")(registerEndpoint)
		registerEndpoint = limiter(registerEndpoint)
		registerEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "GetUser",
			Timeout: 30 * time.Second,
		}))(registerEndpoint)

		loginEndPoint = grpctransport.NewClient(
			conn,
			"pb.UserRpcService",
			"Login",
			encodeGRPCLoginRequest,
			decodeGRPCLoginResponse,
			pb.LoginResponse{},
			grpctransport.ClientBefore(opentracing.ContextToGRPC(tracer, logger)),
		).Endpoint()
		loginEndPoint = opentracing.TraceClient(tracer, "Login")(loginEndPoint)
		loginEndPoint = limiter(loginEndPoint)
		loginEndPoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Login",
			Timeout: 30 * time.Second,
		}))(loginEndPoint)
	}
	return u_endpoint.Set{
		GetUserEndpoint:  getUserEndpoint,
		RegisterEndpoint: registerEndpoint,
		LoginEndpoint:    loginEndPoint,
	}
}

func encodeGRPCGetUserRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(m_user.GetUserRequest)
	return &pb.GetUserRequest{Userid: req.A}, nil
}

func encodeGRPCRegisterRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(m_user.RegisterRequest)
	return &pb.RegisterRequest{
		Username:  req.Username,
		Firstname: req.FirstName,
		Lastname:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	}, nil
}

func encodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(m_user.LoginRequest)
	return &pb.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

func decodeGRPCGetUserResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GetUserResponse)
	return m_user.GetUserResponse{V: m_user.User{
		FirstName: reply.V.Firstname,
		LastName:  reply.V.Lastname,
		Email:     reply.V.Email,
		Username:  reply.V.Username,
		Password:  reply.V.Password,
		Salt:      "",
		UserID:    reply.V.Userid,
	}, Err: str2err(reply.Err)}, nil
}

func decodeGRPCRegisterResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.RegisterResponse)
	return m_user.RegisterUserResponse{
		ID:  reply.Id,
		Err: str2err(reply.Err),
	}, nil
}

func decodeGRPCLoginResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply, _ := grpcReply.(*pb.LoginResponse)
	return m_user.LoginResponse{
		User:  pbUser2ModelUser(*reply.User),
		Token: reply.Token,
	}, nil
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

func pbUser2ModelUser(record pb.UserRecord) *m_user.User {
	return &m_user.User{
		Username:  record.Username,
		Email:     record.Email,
		Password:  record.Password,
		FirstName: record.Firstname,
		LastName:  record.Lastname,
		UserID:    record.Userid,
	}
}

func modelUser2PbUser(model m_user.User) *pb.UserRecord {
	return &pb.UserRecord{
		Username:  model.Username,
		Email:     model.Email,
		Password:  model.Password,
		Firstname: model.FirstName,
		Lastname:  model.LastName,
		Userid:    model.UserID,
	}
}
