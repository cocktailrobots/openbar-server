package openbardb

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"strconv"
)

const (
	pumpsTable  = "pumps"
	configTable = "config"
	keyCol      = "key"

	NumPumpsConfigKey = "num_pumps"
)

type RequiredKey struct {
	name         string
	defaultValue string
}

var requiredKeys = map[string]string{
	NumPumpsConfigKey: "0",
}

func GetConfig(ctx context.Context, tx *dbr.Tx) (map[string]string, error) {
	type kv struct {
		Key   string `db:"key"`
		Value string `db:"value"`
	}

	var configValues []kv
	_, err := tx.Select("*").From(configTable).LoadContext(ctx, &configValues)
	if err != nil {
		return nil, err
	}

	config := make(map[string]string)
	for i := range configValues {
		config[configValues[i].Key] = configValues[i].Value
	}

	return config, nil
}

func SetConfig(ctx context.Context, tx *dbr.Tx, config map[string]string) error {
	for k, defVal := range requiredKeys {
		if _, ok := config[k]; !ok {
			config[k] = defVal
		}
	}

	_, err := tx.DeleteFrom(configTable).ExecContext(ctx)
	if err != nil {
		return err
	}

	query := "INSERT INTO Config (key, value) VALUES" + "\n"
	first := true
	for k := range config {
		if !first {
			query += ",\n"
		}

		query += fmt.Sprintf("('%s', '%s')", k, config[k])
		first = false
	}

	ins := tx.InsertBySql(query)
	_, err = ins.ExecContext(ctx)
	if err != nil {
		return err
	}

	return numPumpsUpdated(ctx, tx, config[NumPumpsConfigKey])
}

func DeleteConfigValues(ctx context.Context, tx *dbr.Tx, keys ...string) error {
	updatedReqKeys := make(map[string]string)
	for i, key := range keys {
		if defaultVal, ok := requiredKeys[key]; ok {
			updatedReqKeys[key] = defaultVal
			keys = append(keys[:i], keys[i+1:]...)
		}
	}

	if len(keys) > 0 {
		_, err := tx.DeleteFrom(configTable).Where(dbr.Eq(keyCol, keys)).ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	if len(updatedReqKeys) > 0 {
		for k, v := range updatedReqKeys {
			_, err := tx.Update(configTable).Set("value", v).Where(dbr.Eq(keyCol, k)).ExecContext(ctx)
			if err != nil {
				return err
			}
		}
	}

	if _, ok := updatedReqKeys[NumPumpsConfigKey]; ok {
		return numPumpsUpdated(ctx, tx, "0")
	}

	return nil
}

func numPumpsUpdated(ctx context.Context, tx *dbr.Tx, numPumpsStr string) error {
	numPumps, err := strconv.ParseInt(numPumpsStr, 10, 64)
	if err != nil {
		return err
	} else if numPumps < 0 {
		return fmt.Errorf("numPumps must be >= 0")
	}

	_, err = tx.DeleteFrom(pumpsTable).Where(dbr.Gte(idxCol, numPumps)).ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom(FluidsTable).Where(dbr.Gte(idxCol, numPumps)).ExecContext(ctx)
	if err != nil {
		return err
	}

	if numPumps > 0 {
		ins := tx.InsertInto(pumpsTable).Ignore().Columns(idxCol)
		for i := int64(0); i < numPumps; i++ {
			ins = ins.Values(i)
		}

		_, err = ins.ExecContext(ctx)
		if err != nil {
			return err
		}

		ins = tx.InsertInto(FluidsTable).Ignore().Columns(idxCol)
		for i := int64(0); i < numPumps; i++ {
			ins = ins.Values(i)
		}

		_, err = ins.ExecContext(ctx)
		return err
	}

	return nil
}
