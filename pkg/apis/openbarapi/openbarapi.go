package openbarapi

import (
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/hardware"
	"github.com/cocktailrobots/openbar-server/pkg/util/dbutils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type OpenBarAPI struct {
	*apis.API
	hw    hardware.Hardware
	ashwr *hardware.AsyncHWRunner
}

func New(logger *zap.Logger, txp dbutils.TxProvider, rtr *mux.Router, hw hardware.Hardware) *OpenBarAPI {
	api := &OpenBarAPI{
		API:   apis.NewAPI(logger, txp, rtr),
		hw:    hw,
		ashwr: hardware.NewAsyncHWRunner(hw),
	}

	rtr.HandleFunc("/", api.DefaultHandler)
	rtr.HandleFunc("/fluids", api.FluidsHandler)
	rtr.HandleFunc("/config", api.ConfigHandler)
	rtr.HandleFunc("/config/{key}", api.ConfigValueHandler)
	rtr.HandleFunc("/menus", api.MenusHandler)
	rtr.HandleFunc("/menus/{name}", api.MenuHandler)
	rtr.HandleFunc("/menus/{name}/recipes", api.MenuRecipesHandler)
	rtr.HandleFunc("/menus/{name}/recipes/{id}", api.MenuRecipeHandler)
	rtr.HandleFunc("/make", api.MakeHandler)
	rtr.HandleFunc("/buttons", api.ButtonsHandler)
	rtr.HandleFunc("/networking", api.NetworkingHandler)
	rtr.HandleFunc("/shutdown", api.ShutdownHandler)

	return api
}
