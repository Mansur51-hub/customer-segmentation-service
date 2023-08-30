package repository

import (
	"context"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"time"
)

type Repository interface {
	CreateUser(ctx context.Context, id int) error
	UserExists(ctx context.Context, id int) (bool, error)
	CreateSegment(ctx context.Context, slug string) (model.Segment, error)
	DeleteSegment(ctx context.Context, slug string) error
	CreateSegmentMembership(ctx context.Context, userId int, slug string, ttl time.Duration) (model.SegmentMembership, error)
	DeleteSegmentMembership(ctx context.Context, userId int, slug string) error
	GetUserSegments(ctx context.Context, userId int, limit uint64, offset uint64) ([]string, error)
}
