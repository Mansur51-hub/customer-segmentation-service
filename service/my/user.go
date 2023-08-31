package my

import (
	"context"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type UserService struct {
	repos *repository.Repositories
}

func (m *UserService) CreateUser(ctx context.Context, userId int) error {
	return m.repos.User.CreateUser(ctx, userId)
}

func (m *UserService) CreateUserSegments(ctx context.Context, userId int, slugsToAdd []string, ttl []time.Duration, slugsToRemove []string) ([]model.SegmentMembership, error) {
	log.Info().Int("user_id", userId).Msg("creating user segments...")

	if len(slugsToAdd) != len(ttl) {
		return nil, fmt.Errorf("error create user segments slugs and exp time len are not same")
	}

	ok, err := m.repos.UserExists(ctx, userId)

	if err != nil {
		log.Debug().Err(err).Int("user_id", userId).Msg("create user")
		return nil, err
	}

	if !ok {
		if err := m.repos.User.CreateUser(ctx, userId); err != nil {
			log.Debug().Err(err).Msg("create user")
			return nil, err
		}
	}

	memberships := make([]model.SegmentMembership, 0, len(slugsToAdd))

	for i, val := range slugsToAdd {
		if mmbr, err := m.repos.Membership.CreateSegmentMembership(ctx, userId, val, ttl[i]); err != nil {
			log.Debug().Err(err).Int("user", userId).Str("slug", val).Msg("create membership")
			return nil, err
		} else {
			memberships = append(memberships, mmbr)
		}
	}

	for _, val := range slugsToRemove {
		if err := m.repos.Membership.DeleteSegmentMembership(ctx, userId, val); err != nil {
			log.Debug().Err(err).Int("user", userId).Str("slug", val).Msg("delete membership")
			return nil, err
		}
	}

	return memberships, nil
}

func (m *UserService) GetUserSegments(ctx context.Context, userId int, limit, offset uint64) ([]string, error) {
	return m.repos.User.GetUserSegments(ctx, userId, limit, offset)
}

func NewUserService(repos *repository.Repositories) *UserService {
	return &UserService{repos: repos}
}
