package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/repoerrs"
	"github.com/jackc/pgx/v5/pgconn"
)

const DurationNilValue = "0s"

type SegmentRepo struct {
	pg *postgres.Postgres
}

func (r *SegmentRepo) CreateSegment(ctx context.Context, slug string, percent int) (model.Segment, error) {
	sql, args, _ := r.pg.Sq.Insert("segments").
		Columns("slug").
		Values(slug).
		Suffix("RETURNING id").
		ToSql()

	var id int

	err := r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return model.Segment{}, repoerrs.ErrAlreadyExists
			}
		}

		return model.Segment{}, fmt.Errorf("error create segment: %w", err)
	}

	return model.Segment{Id: id, Slug: slug, Percent: percent}, nil
}

func (r *SegmentRepo) GetSegmentMembers(ctx context.Context, slug string) ([]int, error) {
	sql, args, _ := r.pg.Sq.Select("slug").
		From("segments").
		Where("slug = ?", slug).
		ToSql()

	res, err := r.pg.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("error get segment: %w", err)
	}

	if res.RowsAffected() == 0 {
		return nil, repoerrs.ErrNotExists
	}

	sql, args, _ = r.pg.Sq.Select("user_id").
		From("memberships").
		Where("segment_slug = ?", slug).
		ToSql()

	rows, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("error get slug members: %w", err)
	}

	users := make([]int, 0, res.RowsAffected())

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error get user_id: %w", err)
		}

		users = append(users, id)
	}
	return users, nil
}

func (r *SegmentRepo) DeleteSegment(ctx context.Context, slug string) error {

	sql, args, _ := r.pg.Sq.Delete("segments").
		Where("slug = ?", slug).
		ToSql()

	res, err := r.pg.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error delete segment: %w", err)
	}

	if res.RowsAffected() == 0 {
		return repoerrs.ErrNotExists
	}

	return nil
}

func NewSegmentRepo(pg *postgres.Postgres) *SegmentRepo {
	return &SegmentRepo{pg: pg}
}
