package service

import (
	"context"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"time"
)

type Service interface {
	CreateUser(ctx context.Context, userId int) error
	CreateSegment(ctx context.Context, slug string) (model.Segment, error)
	DeleteSegment(ctx context.Context, slug string) error
	GetUserSegments(ctx context.Context, userId int, limit, offset uint64) ([]string, error)
	CreateUserSegments(ctx context.Context, userId int, slugsToAdd []string, ttl []time.Duration, slugsToRemove []string) ([]model.SegmentMembership, error)
}
