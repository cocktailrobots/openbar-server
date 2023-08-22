package openbarapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/gocraft/dbr/v2"
	"go.uber.org/zap"
	"net/http"
)

// FluidsHandler handles requests to /fluids
func (api *OpenBarAPI) FluidsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodGet {
		api.GetFluidsHandler(ctx, w, r)
	} else if r.Method == http.MethodPost {
		api.PostFluidsHandler(ctx, w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// PostFluidsHandler handles POST requests to /fluids
func (api *OpenBarAPI) PostFluidsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var fluidsReq wire.Fluids
	err := json.NewDecoder(r.Body).Decode(&fluidsReq)
	if err != nil {
		api.Logger().Info("Error decoding request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := db.UpdateFluids(ctx, tx, fluidsReq.ToDbFluids())
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
		fluids, err := db.ListFluids(ctx, tx)
		if err != nil {
			return fmt.Errorf("error getting fluids from db: %w", err)
		}

		fluidsResp = wire.FromDbFluids(fluids)
		return nil
	})

	api.Respond(w, r, fluidsResp, err)
}
