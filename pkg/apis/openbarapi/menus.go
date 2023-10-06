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

func (api *OpenBarAPI) MenusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		api.GetMenus(ctx, w, r)
	case http.MethodPost:
		api.PostMenus(ctx, w, r)
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodPost}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

func (api *OpenBarAPI) GetMenus(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var menus []string
	err := api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		menus, err = openbardb.GetMenuNames(ctx, tx)
		return err
	})

	if err != nil {
		api.Respond(w, r, nil, err)
	}

	if menus == nil {
		menus = []string{}
	}

	api.Respond(w, r, menus, nil)
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
			err := openbardb.CreateMenu(ctx, tx, menu.Name, menu.Ingredients)
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
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet, http.MethodDelete, http.MethodPatch}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
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
		menu, err := openbardb.GetMenu(ctx, tx, menuName)
		if err != nil {
			return fmt.Errorf("error getting menu: %w", err)
		}

		menus = wire.FromDbMenus([]*openbardb.Menu{menu})[0]
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
		err := openbardb.DeleteMenu(ctx, tx, menuName)
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
		_, err := openbardb.GetMenu(ctx, tx, menuName)
		if err != nil {
			return fmt.Errorf("error getting menu '%s': %w", menuName, err)
		}

		err = openbardb.UpdateMenu(ctx, tx, menuName, menu.Ingredients, menu.RecipeIds)
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
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
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
		menu, err := openbardb.GetMenu(ctx, tx, menuName)
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
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodDelete, http.MethodPost}, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
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
		err := openbardb.RemoveMenuItem(ctx, tx, menuName, recipeId)
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
		err := openbardb.AddMenuItem(ctx, tx, menuName, recipeId)
		if err != nil {
			return fmt.Errorf("error adding menu item: %w", err)
		}

		return tx.Commit()
	})

	api.Respond(w, r, nil, err)
}
