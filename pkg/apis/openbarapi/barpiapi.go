package openbarapi

import (
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/util/dbutils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type OpenBarAPI struct {
	*apis.API
}

func New(logger *zap.Logger, txp dbutils.TxProvider, rtr *mux.Router) *OpenBarAPI {
	api := &OpenBarAPI{
		API: apis.NewAPI(logger, txp, rtr),
	}

	rtr.HandleFunc("/fluids", api.FluidsHandler)
	rtr.HandleFunc("/config", api.ConfigHandler)
	rtr.HandleFunc("/config/{key}", api.ConfigValueHandler)
	rtr.HandleFunc("/menus", api.MenusHandler)
	rtr.HandleFunc("/menus/{name}", api.MenuHandler)
	rtr.HandleFunc("/menus/{name}/recipes", api.MenuRecipesHandler)
	rtr.HandleFunc("/menus/{name}/recipes/{id}", api.MenuRecipeHandler)

	return api
}
