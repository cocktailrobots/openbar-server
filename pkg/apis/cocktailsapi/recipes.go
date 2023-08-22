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
	"strings"
)

// RecipesHandler handles requests to /recipes.
func (api *CocktailsAPI) RecipesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.ListRecipesHandler(ctx, w, r)
	case http.MethodPost:
		api.PostRecipesHandler(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ListRecipesHandler handles requests to GET /recipes.
func (api *CocktailsAPI) ListRecipesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var fluids []string
	if r.URL.Query().Has("fluids") {
		fluidsStr := r.URL.Query().Get("fluids")
		fluids = strings.Split(fluidsStr, ",")

		for i, fluid := range fluids {
			fluids[i] = strings.ToUpper(strings.TrimSpace(fluid))

			if fluids[i] == "" {
				fluids = append(fluids[:i], fluids[i+1:]...)
			}
		}
	}

	var recipes wire.Recipes
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var recipeToIngredients []cocktailsdb.Recipe
		var err error
		if len(fluids) > 0 {
			recipeToIngredients, err = cocktailsdb.GetRecipesForIngredients(ctx, tx, fluids)
		} else {
			recipeToIngredients, err = cocktailsdb.GetRecipes(ctx, tx)
		}

		if err != nil {
			return fmt.Errorf("error getting recipes from db: %w", err)
		}

		recipes = wire.FromDbRecipes(recipeToIngredients)
		return nil
	})

	api.Respond(w, r, recipes, err)
}

// PostRecipesHandler handles requests to POST /recipes.
func (api *CocktailsAPI) PostRecipesHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var reqRecipes wire.Recipes
	err := json.NewDecoder(r.Body).Decode(&reqRecipes)
	if err != nil {
		return
	}

	recipes := reqRecipes.ToDbRecipes()
	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		for _, recipe := range recipes {
			err := cocktailsdb.CreateRecipe(ctx, tx, &recipe)
			if err != nil {
				return fmt.Errorf("error updating recipe '%s': %w", recipe.Id, err)
			}
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

// RecipeHandler handles requests to /recipes/{recipe}.
func (api *CocktailsAPI) RecipeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetRecipe(ctx, w, r)
	case http.MethodDelete:
		api.DeleteRecipe(ctx, w, r)
	case http.MethodPatch:
		api.PatchRecipe(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// GetRecipe handles requests to GET /recipes/{recipe}.
func (api *CocktailsAPI) GetRecipe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	recipeId := pathTokens[len(pathTokens)-1]
	var recipe wire.Recipe
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		recipes, err := cocktailsdb.GetRecipesById(ctx, tx, recipeId)
		if err != nil {
			return fmt.Errorf("error getting recipes from db: %w", err)
		}

		if len(recipes) != 1 {
			return fmt.Errorf("recipe %s not found: %w", recipeId, apis.ErrNotFound)
		}

		recipe = wire.FromDbRecipes(recipes)[0]
		return nil
	})

	api.Respond(w, r, recipe, err)
}

// DeleteRecipe handles requests to DELETE /recipes/{recipe}.
func (api *CocktailsAPI) DeleteRecipe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	recipeId := pathTokens[len(pathTokens)-1]
	api.Logger().Info("Deleting recipe", zap.String("recipe", recipeId))
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.DeleteRecipes(ctx, tx, recipeId)
		if err != nil {
			return fmt.Errorf("error deleting recipe '%s': %w", recipeId, err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

// PatchRecipe handles requests to PATCH /recipes/{recipe}.
func (api *CocktailsAPI) PatchRecipe(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	pathTokens := apis.GetPathTokens(r)
	if len(pathTokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	recipeId := pathTokens[len(pathTokens)-1]
	var reqRecipe wire.Recipes
	err := json.NewDecoder(r.Body).Decode(&reqRecipe)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	recipe := reqRecipe.ToDbRecipes()[0]

	if recipe.Id != recipeId {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := cocktailsdb.UpdateRecipe(ctx, tx, &recipe)
		if err != nil {
			return fmt.Errorf("error updating recipe '%s': %w", recipeId, err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}
