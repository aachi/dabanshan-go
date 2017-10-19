package service

import (
	"context"
	"errors"
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/laidingqing/dabanshan/pb"
)

// Storage
var (
	mem map[int64]map[int64]*pb.ProductRecord
	mu  sync.RWMutex
)

func init() {
	mem = make(map[int64]map[int64]*pb.ProductRecord)
}

// Service describes a service that adds things together.
type Service interface {
	GetProducts(ctx context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error)
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
	ErrUserNotFound = errors.New("user not found")
)

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() Service {
	return basicService{}
}

type basicService struct{}

func (s basicService) GetProducts(_ context.Context, req *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	userID := req.GetCreatorid()
	size := req.GetSize()
	products := []*pb.ProductRecord{}
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := mem[userID]; !ok {
		return nil, ErrUserNotFound
	} else {
		for _, f := range v {
			if size <= 0 {
				break
			}
			products = append(products, f)
			size--
		}
	}
	return &pb.GetProductsResponse{Products: products}, nil
}
