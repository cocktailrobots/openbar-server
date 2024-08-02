package openbardb

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
)

const (
	PumpsTable = "pumps"

	mlPerSecCol = "ml_per_sec"
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
	ins := tx.InsertInto(PumpsTable).Ignore().Columns("idx", "ml_per_sec")
	for i := range pumps {
		ins.Record(&pumps[i])
	}

	_, err := ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert pumps: %w", err)
	}

	for i := range pumps {
		_, err := tx.Update(PumpsTable).Set("ml_per_sec", pumps[i].MlPerSec).Where(dbr.Eq(idxCol, pumps[i].Idx)).ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
