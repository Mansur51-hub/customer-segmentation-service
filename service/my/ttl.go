package my

import (
	"context"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"github.com/rs/zerolog/log"
	"time"
)

type TtlService struct {
	repos *repository.Repositories
}

func (s *TtlService) Exec(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)

	tickerChan := make(chan bool)

	go func() {
		for {
			select {
			case <-tickerChan:
				return
			// interval task
			case tm := <-ticker.C:
				log.Info().Str("time", tm.Format(time.ANSIC)).Msg("deleting expired rows")
				err := s.deleteExpiredMemberships(ctx)

				if err != nil {
					log.Info().Err(err).Msg("delete expired rows")
					ticker.Stop()
					tickerChan <- true
				}
			}
		}
	}()
}

func NewTtlService(repos *repository.Repositories) *TtlService {
	return &TtlService{repos: repos}
}

func (s *TtlService) deleteExpiredMemberships(ctx context.Context) error {
	mmbrs, err := s.repos.Membership.GetExpiredMemberships(ctx)

	if err != nil {
		return err
	}

	for _, mmbr := range mmbrs {
		err = s.repos.Membership.DeleteSegmentMembership(ctx, mmbr.UserId, mmbr.SegmentSlug)

		if err != nil {
			return err
		}
	}

	return nil
}
