package pgx

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
	"github.com/Mansur51-hub/customer-segmentation-service/repository/repoerrs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type UserRepo struct {
	pg *postgres.Postgres
}

func (r *UserRepo) GetUserSegments(ctx context.Context, userId int, limit uint64, offset uint64) ([]string, error) {
	t := time.Now()

	sql, args, _ := r.pg.Sq.Select("segment_slug").
		From("memberships").
		Where("user_id = ? and (expired_at is null or expired_at >= ?)", userId, t).
		Limit(limit).
		Offset(offset).
		ToSql()

	rows, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("error get user: %d segments: %w", userId, err)
	}

	var slugs []string

	for rows.Next() {
		var slug string

		err = rows.Scan(&slug)

		if err != nil {
			return nil, fmt.Errorf("error load {user: %d}: %w", userId, err)
		}

		slugs = append(slugs, slug)
	}

	return slugs, nil
}

func (r *UserRepo) UserExists(ctx context.Context, id int) (bool, error) {
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

func (r *UserRepo) CreateUser(ctx context.Context, id int) error {
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

func (r *UserRepo) GetSample(ctx context.Context, percent int) ([]int, error) {
	sql, args, _ := r.pg.Sq.Select("COUNT(*)").
		From("users").
		ToSql()

	var count int

	err := r.pg.Pool.QueryRow(ctx, sql, args...).Scan(&count)

	if err != nil {
		return nil, err
	}

	val := count * percent / 100

	sql, args, _ = r.pg.Sq.Select("*").
		From("users").
		Where("random() < 0.5").
		Limit(uint64(val)).
		ToSql()

	res, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	users := make([]int, 0, val)

	for res.Next() {
		var id int

		err = res.Scan(&id)

		if err != nil {
			return nil, err
		}

		users = append(users, id)
	}

	return users, nil
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg: pg}
}
