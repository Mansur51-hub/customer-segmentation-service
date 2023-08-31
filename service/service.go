package service

import (
	"context"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"github.com/Mansur51-hub/customer-segmentation-service/service/my"
	"time"
)

type OperationService interface {
	GetOperations(ctx context.Context, year int, month int, limit, offset uint64) ([]byte, error)
}

type UserService interface {
	CreateUser(ctx context.Context, userId int) error
	CreateUserSegments(ctx context.Context, userId int, slugsToAdd []string, ttl []time.Duration, slugsToRemove []string) ([]model.SegmentMembership, error)
	GetUserSegments(ctx context.Context, userId int, limit, offset uint64) ([]string, error)
}

type SegmentService interface {
	CreateSegment(ctx context.Context, slug string, percent int) (model.Segment, error)
	DeleteSegment(ctx context.Context, slug string) error
}

type TtlService interface {
	Exec(ctx context.Context)
}

type Services struct {
	UserService
	OperationService
	SegmentService
	TtlService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		UserService:      my.NewUserService(repos),
		SegmentService:   my.NewSegmentService(repos),
		OperationService: my.NewOperationService(repos),
		TtlService:       my.NewTtlService(repos),
	}
}
