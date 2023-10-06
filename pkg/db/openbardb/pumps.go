package openbardb

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
)

const (
	PumpsTable = "pumps"
)

type Pump struct {
	Idx      int     `db:"idx"`
	MlPerSec float64 `db:"ml_per_sec"`
}

func CountPumpRows(ctx context.Context, tx *dbr.Tx) (int, error) {
	var count int
	_, err := tx.Select("COUNT(*)").From(PumpsTable).LoadContext(ctx, &count)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func ListPumps(ctx context.Context, tx *dbr.Tx) ([]Pump, error) {
	var pumps []Pump
	_, err := tx.Select("*").From(PumpsTable).OrderBy("idx").LoadContext(ctx, &pumps)
	if err != nil {
		return nil, err
	}

	return pumps, nil
}

func UpdatePumps(ctx context.Context, tx *dbr.Tx, pumps []Pump) error {
	_, err := tx.DeleteFrom(PumpsTable).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to clear pumps: %w", err)
	}

	ins := tx.InsertInto(PumpsTable).Columns("idx", "ml_per_sec")
	for i := range pumps {
		ins.Record(&pumps[i])
	}

	res, err := ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert pumps: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if rowsAffected != int64(len(pumps)) {
		return fmt.Errorf("failed to insert all pumps: %w", err)
	}

	return nil
}
