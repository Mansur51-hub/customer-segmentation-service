package service

import (
	"context"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type MyService struct {
	repo *repository.PgRepo
}

func (m *MyService) CreateUser(ctx context.Context, userId int) error {
	return m.repo.CreateUser(ctx, userId)
}

func (m *MyService) CreateSegment(ctx context.Context, slug string) (model.Segment, error) {
	return m.repo.CreateSegment(ctx, slug)
}

func (m *MyService) DeleteSegment(ctx context.Context, slug string) error {
	return m.repo.DeleteSegment(ctx, slug)
}

func (m *MyService) GetUserSegments(ctx context.Context, userId int, limit, offset uint64) ([]string, error) {
	return m.repo.GetUserSegments(ctx, userId, limit, offset)
}

func (m *MyService) CreateUserSegments(ctx context.Context, userId int, slugsToAdd []string, ttl []time.Duration, slugsToRemove []string) ([]model.SegmentMembership, error) {
	log.Info().Int("user_id", userId).Msg("creating user segments...")

	if len(slugsToAdd) != len(ttl) {
		return nil, fmt.Errorf("error create user segments slugs and exp time len are not same")
	}

	ok, err := m.repo.UserExists(ctx, userId)

	if err != nil {
		log.Debug().Err(err).Int("user_id", userId).Msg("create user")
		return nil, err
	}

	if !ok {
		if err := m.repo.CreateUser(ctx, userId); err != nil {
			log.Debug().Err(err).Msg("create user")
			return nil, err
		}
	}

	memberships := make([]model.SegmentMembership, 0, len(slugsToAdd))

	for i, val := range slugsToAdd {
		if mmbr, err := m.repo.CreateSegmentMembership(ctx, userId, val, ttl[i]); err != nil {
			log.Debug().Err(err).Int("user", userId).Str("slug", val).Msg("create membership")
			return nil, err
		} else {
			memberships = append(memberships, mmbr)
		}
	}

	for _, val := range slugsToRemove {
		if err := m.repo.DeleteSegmentMembership(ctx, userId, val); err != nil {
			log.Debug().Err(err).Int("user", userId).Str("slug", val).Msg("delete membership")
			return nil, err
		}
	}

	return memberships, nil
}

func NewMyService(repo *repository.PgRepo) *MyService {
	return &MyService{repo: repo}
}
