package openbarapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db/openbardb"
	"github.com/gocraft/dbr/v2"
	"net/http"
)

func (api *OpenBarAPI) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.getConfig(ctx, w, r)
	case http.MethodPost:
		api.setConfig(ctx, w, r)
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodPost}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

func (api *OpenBarAPI) getConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var cfg wire.Config
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		cfg, err = openbardb.GetConfig(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get config from db: %w", err)
		}

		return nil
	})

	api.Respond(w, r, cfg, err)
}

func (api *OpenBarAPI) setConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var reqCfg wire.Config
	err := json.NewDecoder(r.Body).Decode(&reqCfg)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		err = openbardb.SetConfig(ctx, tx, reqCfg)
		if err != nil {
			return fmt.Errorf("failed to set config in db: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) ConfigValueHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.getConfigValue(ctx, w, r)
	case http.MethodPatch:
		api.updateConfigValue(ctx, w, r)
	case http.MethodPost:
		api.setConfigValue(ctx, w, r)
	case http.MethodDelete:
		api.deleteConfigValue(ctx, w, r)
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

func (api *OpenBarAPI) getConfigValue(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}
	configKey := pathTokens[len(pathTokens)-1]

	var value string
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		cfg, err := openbardb.GetConfig(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get config from db: %w", err)
		}

		var ok bool
		value, ok = cfg[configKey]
		if !ok {
			return fmt.Errorf("config key %s not found: %w", configKey, dbr.ErrNotFound)
		}

		return nil
	})

	resp := wire.Config{
		configKey: value,
	}
	api.Respond(w, r, resp, err)
}

func (api *OpenBarAPI) updateConfigValue(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}
	configKey := pathTokens[len(pathTokens)-1]

	var reqCfg wire.Config
	err := json.NewDecoder(r.Body).Decode(&reqCfg)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	if len(reqCfg) != 1 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	} else if _, ok := reqCfg[configKey]; !ok {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		cfg, err := openbardb.GetConfig(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get config from db: %w", err)
		}

		var ok bool
		_, ok = cfg[configKey]
		if !ok {
			return fmt.Errorf("config key %s not found: %w", configKey, dbr.ErrNotFound)
		}

		cfg[configKey] = reqCfg[configKey]
		err = openbardb.SetConfig(ctx, tx, cfg)
		if err != nil {
			return fmt.Errorf("failed to set config in db: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) setConfigValue(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}
	configKey := pathTokens[len(pathTokens)-1]

	var reqCfg wire.Config
	err := json.NewDecoder(r.Body).Decode(&reqCfg)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	if len(reqCfg) != 1 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	} else if _, ok := reqCfg[configKey]; !ok {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		cfg, err := openbardb.GetConfig(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to get config from db: %w", err)
		}

		_, ok := cfg[configKey]
		if ok {
			return fmt.Errorf("config key %s already exists: %w", configKey, apis.ErrAlreadyExists)
		}

		cfg[configKey] = reqCfg[configKey]
		err = openbardb.SetConfig(ctx, tx, cfg)
		if err != nil {
			return fmt.Errorf("failed to set config in db: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) deleteConfigValue(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}
	configKey := pathTokens[len(pathTokens)-1]

	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := openbardb.DeleteConfigValues(ctx, tx, configKey)
		if err != nil {
			return fmt.Errorf("failed to delete config values from db: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}
