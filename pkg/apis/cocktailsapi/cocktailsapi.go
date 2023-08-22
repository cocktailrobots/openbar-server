package cocktailsapi

import (
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/util/dbutils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type CocktailsAPI struct {
	*apis.API
}

func New(logger *zap.Logger, txp dbutils.TxProvider, rtr *mux.Router) *CocktailsAPI {
	api := &CocktailsAPI{
		API: apis.NewAPI(logger, txp, rtr),
	}

	rtr.HandleFunc("/cocktails", api.CocktailsHandler)
	rtr.HandleFunc("/cocktails/{name}", api.CocktailHandler)
	rtr.HandleFunc("/recipes", api.RecipesHandler)
	rtr.HandleFunc("/recipes/{id}", api.RecipeHandler)
	rtr.HandleFunc("/ingredients", api.IngredientsHandler)
	rtr.HandleFunc("/ingredients/{name}", api.IngredientHandler)

	return api
}
