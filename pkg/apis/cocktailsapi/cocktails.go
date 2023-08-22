package cocktailsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db/cocktailsdb"
	"github.com/gocraft/dbr/v2"
	"net/http"
)

// CocktailsHandler handles requests to /cocktails.
func (api *CocktailsAPI) CocktailsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.ListCocktailsHandler(ctx, w, r)
	case http.MethodPost:
		api.PostCocktailsHandler(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ListCocktailsHandler handles requests to GET /cocktails.
func (api *CocktailsAPI) ListCocktailsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var cocktails []cocktailsdb.Cocktail
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		cocktails, err = cocktailsdb.GetCocktails(ctx, tx)
		if err != nil {
			return fmt.Errorf("error getting cocktails from db: %w", err)
		}

		return nil
	})

	api.Respond(w, r, wire.FromDbCocktails(cocktails), err)
}

// PostCocktailsHandler handles requests to POST /cocktails.
func (api *CocktailsAPI) PostCocktailsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var reqCocktails wire.Cocktails
	err := json.NewDecoder(r.Body).Decode(&reqCocktails)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = reqCocktails.Validate()
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	cocktails := reqCocktails.ToDbCocktails()
	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.AddCocktails(ctx, tx, cocktails...)
		if err != nil {
			return fmt.Errorf("error updating cocktail: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

// CocktailHandler handles requests to /cocktails/{cocktailName}.
func (api *CocktailsAPI) CocktailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodGet {
		api.GetCocktailsHandler(ctx, w, r)
	} else if r.Method == http.MethodPatch {
		api.UpdateCocktailHandler(ctx, w, r)
	} else if r.Method == http.MethodDelete {
		api.DeleteCocktailHandler(ctx, w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetCocktailsHandler handles requests to GET /cocktails/{cocktailName}.
func (api *CocktailsAPI) GetCocktailsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	cocktailName := pathTokens[len(pathTokens)-1]

	var cocktails []cocktailsdb.Cocktail
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		cocktails, err = cocktailsdb.GetCocktailsWithNames(ctx, tx, cocktailName)
		if err != nil {
			return fmt.Errorf("error getting cocktails from db: %w", err)
		}

		return nil
	})

	var respObj any
	if err == nil && len(cocktails) > 0 {
		respObj = wire.FromDbCocktails(cocktails)[0]
	} else if err == nil {
		err = apis.ErrNotFound
	}

	api.Respond(w, r, respObj, err)
}

// UpdateCocktailHandler handles requests to PATCH /cocktails/{cocktailName}.
func (api *CocktailsAPI) UpdateCocktailHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}
	cocktailName := pathTokens[len(pathTokens)-1]

	var cocktail wire.Cocktail
	err := json.NewDecoder(r.Body).Decode(&cocktail)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	if cocktail.Name != cocktailName {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		dbCocktail := wire.Cocktails{cocktail}.ToDbCocktails()[0]
		existing, err := cocktailsdb.GetCocktailsWithNames(ctx, tx, cocktailName)
		if err != nil {
			return fmt.Errorf("error getting cocktails from db: %w", err)
		} else if len(existing) == 0 {
			return apis.ErrNotFound
		}

		err = cocktailsdb.UpdateCocktail(ctx, tx, &dbCocktail)
		if err != nil {
			return fmt.Errorf("error updating cocktail: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

// DeleteCocktailHandler handles requests to DELETE /cocktails/{cocktailName}.
func (api *CocktailsAPI) DeleteCocktailHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}
	cocktailName := pathTokens[len(pathTokens)-1]

	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.DeleteCocktails(ctx, tx, cocktailName)
		if err != nil {
			if errors.Is(err, dbr.ErrNotFound) {
				return apis.ErrNotFound
			}

			return fmt.Errorf("error deleting cocktail: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}
