package my

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/Mansur51-hub/customer-segmentation-service/model"
	"github.com/Mansur51-hub/customer-segmentation-service/repository"
	"strconv"
	"time"
)

type OperationService struct {
	repos *repository.Repositories
}

func NewOperationService(repos *repository.Repositories) *OperationService {
	return &OperationService{repos: repos}
}

func (s *OperationService) GetOperations(ctx context.Context, year int, month int, limit, offset uint64) ([]byte, error) {

	ops, err := s.repos.GetOperations(ctx, year, month, limit, offset)

	if err != nil {
		return nil, err
	}

	recs := getRecords(ops)

	var b bytes.Buffer

	err = csv.NewWriter(&b).WriteAll(recs)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func getRecords(ops []model.Operation) (records [][]string) {
	records = make([][]string, 0, 0)

	for _, op := range ops {
		rec := []string{
			strconv.Itoa(op.UserId),
			op.SegmentSlug,
			op.Type,
			op.CreatedAt.Format(time.ANSIC)}

		records = append(records, rec)
	}

	return records
}
