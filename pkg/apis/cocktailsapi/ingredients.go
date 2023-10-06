package cocktailsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db/cocktailsdb"
	"github.com/gocraft/dbr/v2"
	"go.uber.org/zap"
	"net/http"
)

// IngredientsHandler handles requests to /ingredients.
func (api *CocktailsAPI) IngredientsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.ListIngredientsHandler(ctx, w, r)
	case http.MethodPost:
		api.PostIngredientsHandler(ctx, w, r)
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodPost}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

// ListIngredientsHandler handles requests to GET /ingredients.
func (api *CocktailsAPI) ListIngredientsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var ingredients []cocktailsdb.Ingredient
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		ingredients, err = cocktailsdb.GetIngredients(ctx, tx)
		if err != nil {
			return fmt.Errorf("error getting ingredients from db: %w", err)
		}

		return nil
	})

	api.Respond(w, r, wire.FromDbIngredients(ingredients), err)
}

// PostIngredientsHandler handles requests to POST /ingredients.
func (api *CocktailsAPI) PostIngredientsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var reqIngredients wire.Ingredients
	err := json.NewDecoder(r.Body).Decode(&reqIngredients)
	if err != nil {
		api.Logger().Info("bad request. failed to deserialize", zap.String("path", r.URL.Path), zap.Error(err))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	ingredients := reqIngredients.ToDbIngredients()
	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.CreateIngredients(ctx, tx, ingredients...)
		if err != nil {
			return fmt.Errorf("error creating ingredient in db: %w", err)
		}

		return tx.Commit()
	})

	api.Logger().Info("attempted to create ingredients", zap.Any("ingredients", ingredients), zap.Error(err))

	api.Respond(w, r, nil, err)
}

// IngredientHandler handles requests to /ingredients/{ingredient_name}.
func (api *CocktailsAPI) IngredientHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetIngredientHandler(ctx, w, r)
	case http.MethodPatch:
		api.UpdateIngredientHandler(ctx, w, r)
	case http.MethodDelete:
		api.DeleteIngredientHandler(ctx, w, r)
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodPatch, http.MethodDelete}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

// GetIngredientHandler handles requests to GET /ingredients/{ingredient_name}.
func (api *CocktailsAPI) GetIngredientHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	api.Logger().Info("get ingredient handler")

	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Logger().Info("bad request", zap.String("path", r.URL.Path), zap.Strings("tokens", pathTokens))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	ingredientName := pathTokens[len(pathTokens)-1]

	var ingredients []cocktailsdb.Ingredient
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		ingredients, err = cocktailsdb.GetIngredientsWithNames(ctx, tx, ingredientName)
		if err != nil {
			return fmt.Errorf("error getting ingredients from db: %w", err)
		}

		return nil
	})

	var respObj any
	if err == nil && len(ingredients) > 0 {
		respObj = wire.FromDbIngredients(ingredients)[0]
	} else if err == nil {
		err = apis.ErrNotFound
	}

	api.Respond(w, r, respObj, err)
}

// UpdateIngredientHandler handles requests to PATCH /ingredients/{ingredient_name}.
func (api *CocktailsAPI) UpdateIngredientHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Logger().Info("bad request", zap.String("path", r.URL.Path), zap.Strings("tokens", pathTokens))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	ingredientName := pathTokens[len(pathTokens)-1]

	var reqIngredient wire.Ingredient
	err := json.NewDecoder(r.Body).Decode(&reqIngredient)
	if err != nil {
		api.Logger().Info("bad request", zap.String("path", r.URL.Path), zap.Strings("tokens", pathTokens), zap.Error(err))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	api.Logger().Info("update ingredient", zap.String("name", ingredientName), zap.Any("ingredient", reqIngredient))
	ingredient := wire.Ingredients{reqIngredient}.ToDbIngredients()[0]
	if ingredient.Name != ingredientName {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.UpdateIngredient(ctx, tx, ingredient)
		if err != nil {
			return fmt.Errorf("error updating ingredient in db: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

// DeleteIngredientHandler handles requests to DELETE /ingredients/{ingredient_name}.
func (api *CocktailsAPI) DeleteIngredientHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	ingredientName := pathTokens[len(pathTokens)-1]

	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.DeleteIngredients(ctx, tx, ingredientName)
		if err != nil {
			return fmt.Errorf("error deleting ingredient in db: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}
