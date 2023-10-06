package openbarapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db/openbardb"
	"github.com/gocraft/dbr/v2"
	"go.uber.org/zap"
	"net/http"
)

// FluidsHandler handles requests to /fluids
func (api *OpenBarAPI) FluidsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetFluidsHandler(ctx, w, r)
	case http.MethodPost:
		api.PostFluidsHandler(ctx, w, r)
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodPost}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

// PostFluidsHandler handles POST requests to /fluids
func (api *OpenBarAPI) PostFluidsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var fluidsReq wire.Fluids
	err := json.NewDecoder(r.Body).Decode(&fluidsReq)
	if err != nil {
		api.Logger().Info("Error decoding request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := openbardb.UpdateFluids(ctx, tx, fluidsReq.ToDbFluids())
		if err != nil {
			return fmt.Errorf("error updating fluid: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

// GetFluidsHandler handles GET requests to /fluids
func (api *OpenBarAPI) GetFluidsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var fluidsResp wire.Fluids
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		fluids, err := openbardb.ListFluids(ctx, tx)
		if err != nil {
			return fmt.Errorf("error getting fluids from db: %w", err)
		}

		fluidsResp = wire.FromDbFluids(fluids)
		return nil
	})

	api.Respond(w, r, fluidsResp, err)
}
