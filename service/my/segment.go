package my

import (
	"context"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/pgx"
	"github.com/rs/zerolog/log"
	"time"
)

type SegmentService struct {
	repos *repository.Repositories
}

func (m *SegmentService) CreateSegment(ctx context.Context, slug string, percent int) (model.Segment, error) {
	seg, err := m.repos.Segment.CreateSegment(ctx, slug, percent)

	if err != nil {
		return model.Segment{}, err
	}

	if percent > 0 && percent <= 100 {
		users, err := m.repos.User.GetSample(ctx, percent)

		if err != nil {
			return model.Segment{}, fmt.Errorf("error get users sample: %w", err)
		}

		for _, id := range users {
			d, _ := time.ParseDuration(pgx.DurationNilValue)
			_, err := m.repos.CreateSegmentMembership(ctx, id, slug, d)

			if err != nil {
				return model.Segment{}, fmt.Errorf("error create membreship: %w", err)
			}
		}
	}

	return seg, nil
}

func (m *SegmentService) DeleteSegment(ctx context.Context, slug string) error {
	log.Info().Str("slug", slug).Msg("getting segment memberships")

	users, err := m.repos.Segment.GetSegmentMembers(ctx, slug)

	if err != nil {
		return err
	}

	log.Info().Str("slug", slug).Msg("deleting segment memberships")

	for _, user := range users {
		err = m.repos.DeleteSegmentMembership(ctx, user, slug)
		if err != nil {
			return err
		}
	}

	return m.repos.Segment.DeleteSegment(ctx, slug)
}

func NewSegmentService(repos *repository.Repositories) *SegmentService {
	return &SegmentService{repos: repos}
}
