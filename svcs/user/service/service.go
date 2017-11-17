package service

import (
	"context"
	"errors"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	auth "github.com/laidingqing/dabanshan/svcs/authorize"
	"github.com/laidingqing/dabanshan/svcs/user/db"
	"github.com/laidingqing/dabanshan/svcs/user/model"
)

var (
	// ErrUserNotFound 用户未发现
	ErrUserNotFound = errors.New("not found user")
	// ErrUserAlreadyExisting 用户名已存在
	ErrUserAlreadyExisting = errors.New("username already existing")
)

// Service describes a service that adds things together.
type Service interface {
	GetUser(ctx context.Context, id string) (model.GetUserResponse, error)
	Register(ctx context.Context, RegisterRequest model.RegisterRequest) (model.RegisterUserResponse, error)
	Login(ctx context.Context, login model.LoginRequest) (model.LoginResponse, error)
	Upload(ctx context.Context, manifestName string, manifest io.Reader, fileName string, file io.Reader) (string, error)
}

// New returns a basic Service with all of the expected middlewares wired in.
func New(logger log.Logger, ints, chars metrics.Counter) Service {
	var svc Service
	{
		svc = NewBasicService()
		svc = LoggingMiddleware(logger)(svc)
		svc = InstrumentingMiddleware(ints, chars)(svc)
	}
	return svc
}

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

const ()

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

// GetUser get user by id
func (s basicService) GetUser(_ context.Context, id string) (model.GetUserResponse, error) {
	us, err := db.GetUser(id)
	if err != nil {
		return model.GetUserResponse{V: model.New(), Err: nil}, ErrUserNotFound
	}
	return model.GetUserResponse{
		V:   us,
		Err: nil,
	}, nil
}

// Register user
func (s basicService) Register(ctx context.Context, req model.RegisterRequest) (model.RegisterUserResponse, error) {
	us, err := db.GetUserByName(req.Username)
	if us.Username != "" && err == nil {
		return model.RegisterUserResponse{
			Err: ErrUserAlreadyExisting,
		}, nil
	}

	u := model.New()
	u.Username = req.Username
	u.Password = auth.CalculatePassHash(req.Password, u.Salt)
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	id, err := db.CreateUser(&u)
	return model.RegisterUserResponse{ID: id}, err
}

func (s basicService) Login(ctx context.Context, login model.LoginRequest) (model.LoginResponse, error) {
	u, err := db.GetUserByName(login.Username)
	if err != nil {
		return model.LoginResponse{
			Err: err,
		}, err
	}
	if u.Password != auth.CalculatePassHash(login.Password, u.Salt) {
		return model.LoginResponse{
			Err: ErrUnauthorized,
		}, ErrUnauthorized
	}
	t, err := auth.CreateJWT()
	if err != nil {
		return model.LoginResponse{
			Err: err,
		}, err
	}

	return model.LoginResponse{
		User:  &u,
		Token: t,
	}, nil
}

func (s basicService) Upload(ctx context.Context, manifestName string, manifest io.Reader, fileName string, file io.Reader) (string, error) {

	return "", nil
}

// private func
