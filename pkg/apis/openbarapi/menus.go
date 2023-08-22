package openbarapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/gocraft/dbr/v2"
	"go.uber.org/zap"
	"net/http"
)

func (api *OpenBarAPI) MenusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetMenus(ctx, w, r)
	case http.MethodPost:
		api.PostMenus(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (api *OpenBarAPI) GetMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var menus []string
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		menus, err = db.GetMenuNames(ctx, tx)
		return err
	})

	api.Respond(w, r, menus, err)
}

func (api *OpenBarAPI) PostMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var menus wire.Menus
	err := json.NewDecoder(r.Body).Decode(&menus)
	if err != nil {
		api.Logger().Info("Error decoding request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		for _, menu := range menus {
			err := db.CreateMenu(ctx, tx, menu.Name, menu.Ingredients)
			if err != nil {
				return fmt.Errorf("error creating menu: %w", err)
			}
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) MenuHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetMenu(ctx, w, r)
	case http.MethodDelete:
		api.DeleteMenu(ctx, w, r)
	case http.MethodPatch:
		api.PatchMenu(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (api *OpenBarAPI) GetMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokens := apis.GetPathTokens(r)
	if len(tokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	menuName := tokens[len(tokens)-1]

	var menus wire.Menu
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		menu, err := db.GetMenu(ctx, tx, menuName)
		if err != nil {
			return fmt.Errorf("error getting menu: %w", err)
		}

		menus = wire.FromDbMenus([]*db.Menu{menu})[0]
		return nil
	})

	api.Respond(w, r, menus, err)
}

func (api *OpenBarAPI) DeleteMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokens := apis.GetPathTokens(r)
	if len(tokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	menuName := tokens[len(tokens)-1]

	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := db.DeleteMenu(ctx, tx, menuName)
		if err != nil {
			return fmt.Errorf("error deleting menu: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) PatchMenu(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokens := apis.GetPathTokens(r)
	if len(tokens) != 2 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	menuName := tokens[len(tokens)-1]

	var menu wire.Menu
	err := json.NewDecoder(r.Body).Decode(&menu)
	if err != nil {
		api.Logger().Info("Error decoding request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		_, err := db.GetMenu(ctx, tx, menuName)
		if err != nil {
			return fmt.Errorf("error getting menu '%s': %w", menuName, err)
		}

		err = db.UpdateMenu(ctx, tx, menuName, menu.Ingredients, menu.RecipeIds)
		if err != nil {
			return fmt.Errorf("error updating menu: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) MenuRecipesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetMenuRecipes(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (api *OpenBarAPI) GetMenuRecipes(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokens := apis.GetPathTokens(r)
	if len(tokens) != 3 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	menuName := tokens[len(tokens)-2]

	var recipes []string
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		menu, err := db.GetMenu(ctx, tx, menuName)
		if err != nil {
			return fmt.Errorf("error getting menu: %w", err)
		}

		recipes = menu.RecipeIds
		return nil
	})

	api.Respond(w, r, recipes, err)
}

func (api *OpenBarAPI) MenuRecipeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodDelete:
		api.DeleteRecipeFromMenuHandler(ctx, w, r)
	case http.MethodPost:
		api.AddRecipeToMenuHandler(ctx, w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (api *OpenBarAPI) DeleteRecipeFromMenuHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokens := apis.GetPathTokens(r)
	if len(tokens) != 4 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	menuName := tokens[len(tokens)-3]
	recipeId := tokens[len(tokens)-1]

	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := db.RemoveMenuItem(ctx, tx, menuName, recipeId)
		if err != nil {
			return fmt.Errorf("error removing menu item: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}

func (api *OpenBarAPI) AddRecipeToMenuHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokens := apis.GetPathTokens(r)
	if len(tokens) != 4 {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	menuName := tokens[len(tokens)-3]
	recipeId := tokens[len(tokens)-1]

	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		err := db.AddMenuItem(ctx, tx, menuName, recipeId)
		if err != nil {
			return fmt.Errorf("error adding menu item: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}
