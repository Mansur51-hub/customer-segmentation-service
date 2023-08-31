package pgx

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/repoerrs"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type MembershipRepo struct {
	pg *postgres.Postgres
}

func NewMembershipRepo(pg *postgres.Postgres) *MembershipRepo {
	return &MembershipRepo{pg: pg}
}

func (r *MembershipRepo) DeleteSegmentMembership(ctx context.Context, userId int, slug string) error {
	tx, err := r.pg.Pool.Begin(ctx)

	if err != nil {
		return fmt.Errorf("error tx begin %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	sql, args, _ := r.pg.Sq.Delete("memberships").
		Where("user_id = ? and segment_slug = ?", userId, slug).
		ToSql()

	res, err := tx.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error delete membership {user_id: %d, slug: %s}: %w", userId, slug, err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("error delete membership {user_id: %d, slug: %s}: %w", userId, slug, repoerrs.ErrNotExists)
	}

	sql, args, _ = r.pg.Sq.Insert("operations").
		Columns("user_id", "segment_slug", "created_at", "operation_type").
		Values(userId, slug, time.Now(), model.CustomerRemoveFromSegmentOperationType).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error add delete to operations {user_id: %d, slug: %s}: %w", userId, slug, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error tx commit: %w", err)
	}

	return nil
}

func (r *MembershipRepo) CreateSegmentMembership(ctx context.Context, userId int, slug string, ttl time.Duration) (model.SegmentMembership, error) {
	tx, err := r.pg.Pool.Begin(ctx)

	if err != nil {
		return model.SegmentMembership{}, fmt.Errorf("error tx begin %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	createdAt := time.Now()

	nilTime, _ := time.ParseDuration(DurationNilValue)

	var exTime interface{}

	if ttl == nilTime {
		exTime, _ = dbsql.NullTime{}.Value()
	} else {
		exTime = createdAt.Add(ttl)
	}

	sql, args, _ := r.pg.Sq.Insert("memberships").
		Columns("user_id", "segment_slug", "created_at", "expired_at").
		Values(userId, slug, createdAt, exTime).
		Suffix("RETURNING id, expired_at").
		ToSql()

	var id int
	var expiredAt dbsql.NullTime

	err = tx.QueryRow(ctx, sql, args...).Scan(&id, &expiredAt)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return model.SegmentMembership{}, fmt.Errorf("error create {user_id: %d, slug: %s} slug %w",
					userId, slug, repoerrs.ErrNotExists)
			}

			if pgErr.Code == "23505" {
				return model.SegmentMembership{}, fmt.Errorf("error create {user_id: %d, slug: %s} %w",
					userId, slug, repoerrs.ErrAlreadyExists)
			}
		}
		return model.SegmentMembership{},
			fmt.Errorf("error create membership {user_id: %d slug: %s}: %w", userId, slug, err)
	}

	sql, args, _ = r.pg.Sq.Insert("operations").
		Columns("user_id", "segment_slug", "created_at", "operation_type").
		Values(userId, slug, createdAt, model.CustomerAddToSegmentOperationType).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		return model.SegmentMembership{}, fmt.Errorf("error add membership {user_id: %d slug: %s} to operations: %w",
			userId, slug, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.SegmentMembership{}, fmt.Errorf("error tx commit: %w", err)
	}

	return model.SegmentMembership{
		Id:          id,
		UserId:      userId,
		SegmentSlug: slug,
		CreatedAt:   createdAt,
		ExpiredAt:   expiredAt.Time}, nil
}

func (r *MembershipRepo) GetExpiredMemberships(ctx context.Context) ([]model.SegmentMembership, error) {

	t := time.Now()

	sql, args, _ := r.pg.Sq.Select("*").
		From("memberships").
		Where("expired_at is not null and expired_at <= ?", t).
		ToSql()

	rows, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("error get expired memberships: %w", err)
	}

	memberships := make([]model.SegmentMembership, 0, 0)

	for rows.Next() {
		var mmbr model.SegmentMembership

		err = rows.Scan(&mmbr.Id, &mmbr.UserId, &mmbr.SegmentSlug, &mmbr.CreatedAt, &mmbr.ExpiredAt)

		if err != nil {
			return nil, fmt.Errorf("error load {membership: %d}: %w", mmbr.UserId, err)
		}

		memberships = append(memberships, mmbr)
	}

	return memberships, nil
}
