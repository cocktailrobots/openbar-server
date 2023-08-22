package db

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
)

const (
	FluidsTable = "fluids"

	idxCol   = "idx"
	fluidCol = "fluid"
)

// Fluid represents a fluid in the database where each Idx is the index of the pump that will dispense it.
type Fluid struct {
	Idx   int     `db:"idx"`
	Fluid *string `db:"fluid"`
}

// ListFluids returns all fluids from the database.
func ListFluids(ctx context.Context, tx *dbr.Tx) ([]Fluid, error) {
	var fluids []Fluid
	_, err := tx.Select("*").From(FluidsTable).OrderBy(idxCol).LoadContext(ctx, &fluids)
	if err != nil {
		return nil, err
	}

	return fluids, nil
}

// UpdateFluids updates the fluids in the database.
func UpdateFluids(ctx context.Context, tx *dbr.Tx, fluids []Fluid) error {
	_, err := tx.Update(FluidsTable).Set(fluidCol, nil).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to clear fluids: %w", err)
	}

	query := fmt.Sprintf("REPLACE INTO %s (%s, %s) VALUES\n", FluidsTable, idxCol, fluidCol)
	for i, fluid := range fluids {
		if i != 0 {
			query += ",\n"
		}

		fluidVal := "NULL"
		if fluid.Fluid != nil {
			fluidVal = "'" + *fluid.Fluid + "'"
		}

		query += fmt.Sprintf("(%d, %s)", fluid.Idx, fluidVal)
	}

	ins := tx.InsertBySql(query)
	res, err := ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to update fluids: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected != int64(len(fluids)) {
		//return fmt.Errorf("expected %d rows to be affected, but got %d", len(fluids), rowsAffected)
	}

	return nil
}

// InitFluids initializes the fluids in the database.
func InitFluids(ctx context.Context, tx *dbr.Tx, numFluids int) error {
	_, err := tx.DeleteFrom(FluidsTable).Where(dbr.Gte(idxCol, numFluids)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete fluids: %w", err)
	}

	ins := tx.InsertInto(FluidsTable).Ignore().Columns(idxCol, fluidCol)
	for i := 0; i < numFluids; i++ {
		ins = ins.Record(&Fluid{Idx: i})
	}

	_, err = ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert fluids: %w", err)
	}

	return nil
}
