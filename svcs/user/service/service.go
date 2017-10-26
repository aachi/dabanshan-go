package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/svcs/user/db"
	"github.com/laidingqing/dabanshan/svcs/user/model"
)

var (
	// ErrUserNotFound ..
	ErrUserNotFound = errors.New("not found user")
)

// Service describes a service that adds things together.
type Service interface {
	GetUser(ctx context.Context, id string) (model.GetUserResponse, error)
	Register(ctx context.Context, RegisterRequest model.RegisterRequest) (model.RegisterUserResponse, error)
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
	u := model.New()
	u.Username = req.Username
	u.Password = calculatePassHash(req.Password, u.Salt)
	u.Email = req.Email
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	id, err := db.CreateUser(&u)
	return model.RegisterUserResponse{ID: id}, err
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
