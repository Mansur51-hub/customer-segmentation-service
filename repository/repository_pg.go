package repository

import (
	"context"
	dbsql "database/sql"
	"errors"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/repoerrs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

const DurationNilValue = "0s"

type PgRepo struct {
	pg *postgres.Postgres
}

func (r *PgRepo) GetUserSegments(ctx context.Context, userId int, limit uint64, offset uint64) ([]string, error) {
	t := time.Now()

	sql, args, _ := r.pg.Sq.Select("segment_slug").
		From("memberships").
		Where("user_id = ? and (expired_at is null or expired_at >= ?)", userId, t).
		Limit(limit).
		Offset(offset).
		ToSql()

	rows, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("error get user segments: %w", err)
	}

	var slugs []string

	for rows.Next() {
		var slug string

		err = rows.Scan(&slug)

		if err != nil {
			return nil, fmt.Errorf("error load user slug: %w", err)
		}

		slugs = append(slugs, slug)
	}

	return slugs, nil
}

func (r *PgRepo) DeleteSegmentMembership(ctx context.Context, userId int, slug string) error {
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
		return fmt.Errorf("error delete membership: %w", err)
	}

	if res.RowsAffected() == 0 {
		return repoerrs.ErrNotExists
	}

	sql, args, _ = r.pg.Sq.Insert("operations").
		Columns("user_id", "segment_slug", "created_at", "operation_type").
		Values(userId, slug, time.Now(), model.CustomerRemoveFromSegmentOperationType).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error add delete to operations: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error tx commit: %w", err)
	}

	return nil
}

func (r *PgRepo) UserExists(ctx context.Context, id int) (bool, error) {
	sql, args, _ := r.pg.Sq.Select("id").
		From("users").
		Where("id = ?", id).
		ToSql()

	var userId int

	err := r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&userId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("error user exists : %w", err)
	}

	return true, nil
}

func (r *PgRepo) CreateSegmentMembership(ctx context.Context, userId int, slug string, ttl time.Duration) (model.SegmentMembership, error) {
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
				fmt.Println(pgErr.Detail)
				return model.SegmentMembership{}, fmt.Errorf(repoerrs.ErrNotExists.Error() + " " + pgErr.Detail)
			}
		}
		return model.SegmentMembership{}, fmt.Errorf("error create membership: %w", err)
	}

	sql, args, _ = r.pg.Sq.Insert("operations").
		Columns("user_id", "segment_slug", "created_at", "operation_type").
		Values(userId, slug, createdAt, model.CustomerAddToSegmentOperationType).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)

	if err != nil {
		return model.SegmentMembership{}, fmt.Errorf("error add membership to operations: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.SegmentMembership{}, fmt.Errorf("error tx commit: %w", err)
	}

	return model.SegmentMembership{Id: id, UserId: userId, SegmentSlug: slug, CreatedAt: createdAt, ExpiredAt: expiredAt.Time}, nil
}

func (r *PgRepo) CreateUser(ctx context.Context, id int) error {
	sql, args, _ := r.pg.Sq.Insert("users").
		Columns("id").
		Values(id).
		Suffix("RETURNING id").
		ToSql()

	var userId int

	err := r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&userId)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return repoerrs.ErrAlreadyExists
			}
		}

		return fmt.Errorf("error create user: %w", err)
	}

	return nil
}

func (r *PgRepo) CreateSegment(ctx context.Context, slug string) (model.Segment, error) {
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

	return model.Segment{Id: id, Slug: slug}, nil
}

func (r *PgRepo) DeleteSegment(ctx context.Context, slug string) error {
	sql, args, _ := r.pg.Sq.Select("slug").
		From("segments").
		Where("slug = ?", slug).
		ToSql()

	res, err := r.pg.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error get segment: %w", err)
	}

	if res.RowsAffected() == 0 {
		return repoerrs.ErrNotExists
	}

	var users []int

	sql, args, _ = r.pg.Sq.Select("user_id").
		From("memberships").
		Where("segment_slug = ?", slug).
		ToSql()

	rows, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error get slug members: %w", err)
	}

	for rows.Next() {
		var id int

		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("error get user_id: %w", err)
		}

		users = append(users, id)
	}

	for _, id := range users {
		err = r.DeleteSegmentMembership(ctx, id, slug)

		if err != nil {
			return err
		}
	}

	sql, args, _ = r.pg.Sq.Delete("segments").
		Where("slug = ?", slug).
		ToSql()

	_, err = r.pg.Pool.Exec(ctx, sql, args...)

	if err != nil {
		return fmt.Errorf("error delete segment: %w", err)
	}

	return nil
}

func NewPostgresRepository(pg *postgres.Postgres) *PgRepo {
	return &PgRepo{pg: pg}
}
