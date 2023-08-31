package pgx

import (
	"context"
	"fmt"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/pkg/postgres"
)

type OperationRepo struct {
	pg *postgres.Postgres
}

func NewOperationRepository(pg *postgres.Postgres) *OperationRepo {
	return &OperationRepo{pg: pg}
}

func (r *OperationRepo) GetOperations(ctx context.Context, year int, month int, limit, offset uint64) ([]model.Operation, error) {
	sql, args, _ := r.pg.Sq.Select("*").
		From("operations").
		Where("date_part('year', created_at) = ? and date_part('month', created_at) = ?", year, month).
		Limit(limit).
		Offset(offset).
		ToSql()

	rows, err := r.pg.Pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, fmt.Errorf("error get operations: %w", err)
	}

	operations := make([]model.Operation, 0, 0)

	for rows.Next() {
		var op model.Operation

		err = rows.Scan(&op.Id, &op.UserId, &op.SegmentSlug, &op.CreatedAt, &op.Type)

		if err != nil {
			return nil, fmt.Errorf("error load operation: %w", err)
		}

		operations = append(operations, op)
	}

	return operations, nil
}
