package repository

import (
	"context"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/pgx"
	"time"
)

type Operation interface {
	GetOperations(ctx context.Context, year int, month int, limit, offset uint64) ([]model.Operation, error)
}

type User interface {
	CreateUser(ctx context.Context, id int) error
	UserExists(ctx context.Context, id int) (bool, error)
	GetSample(ctx context.Context, percent int) ([]int, error)
	GetUserSegments(ctx context.Context, userId int, limit uint64, offset uint64) ([]string, error)
}

type Membership interface {
	CreateSegmentMembership(ctx context.Context, userId int, slug string, ttl time.Duration) (model.SegmentMembership, error)
	DeleteSegmentMembership(ctx context.Context, userId int, slug string) error
	GetExpiredMemberships(ctx context.Context) ([]model.SegmentMembership, error)
}

type Segment interface {
	CreateSegment(ctx context.Context, slug string, percent int) (model.Segment, error)
	DeleteSegment(ctx context.Context, slug string) error
	GetSegmentMembers(ctx context.Context, slug string) ([]int, error)
}

type Repositories struct {
	User
	Membership
	Segment
	Operation
}

func NewRepositories(pg *postgres.Postgres) *Repositories {
	return &Repositories{
		User:       pgx.NewUserRepo(pg),
		Membership: pgx.NewMembershipRepo(pg),
		Segment:    pgx.NewSegmentRepo(pg),
		Operation:  pgx.NewOperationRepository(pg),
	}
}
