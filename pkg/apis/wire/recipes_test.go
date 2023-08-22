package wire

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRecipes(t *testing.T) {

	recipes := []db.Recipe{
		{
			Id:          "00000000-0000-0000-0000-000000000000",
			CocktailFk:  "negroni",
			DisplayName: "Negroni",
			Description: ptr("Campari, Gin, and Sweet Vermouth"),
			Directions:  ptr("Mix equal parts of Gin, Campari, and Sweet Vermouth in a glass with ice. Stir. Garnish with an orange peel."),
			Ingredients: []db.RecipeIngredient{
				{
					RecipeIdFk:   "00000000-0000-0000-0000-000000000000",
					IngredientFk: "gin",
					Amount:       1,
				},
				{
					RecipeIdFk:   "00000000-0000-0000-0000-000000000000",
					IngredientFk: "campari",
					Amount:       1,
				},
				{
					RecipeIdFk:   "00000000-0000-0000-0000-000000000000",
					IngredientFk: "sweet_vermouth",
					Amount:       1,
				},
			},
		},
		{
			Id:          "00000000-0000-0000-0000-000000000001",
			CocktailFk:  "boulevardier",
			DisplayName: "Boulevardier",
			Description: ptr("Campari, Bourbon, and Sweet Vermouth"),
			Directions:  ptr("Mix equal parts of Bourbon, Campari, and Sweet Vermouth in a glass with ice. Stir. Garnish with an orange peel."),
			Ingredients: []db.RecipeIngredient{
				{
					RecipeIdFk:   "00000000-0000-0000-0000-000000000001",
					IngredientFk: "bourbon",
					Amount:       1,
				},
				{
					RecipeIdFk:   "00000000-0000-0000-0000-000000000001",
					IngredientFk: "campari",
					Amount:       1,
				},
				{
					RecipeIdFk:   "00000000-0000-0000-0000-000000000001",
					IngredientFk: "sweet_vermouth",
					Amount:       1,
				},
			},
		},
	}

	recipesWire := FromDbRecipes(recipes)
	data, err := json.Marshal(recipesWire)
	require.NoError(t, err)
	require.NoError(t, recipesWire.Validate())

	var decoded Recipes
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	recipes2 := decoded.ToDbRecipes()
	require.True(t, repipesArEqual(recipes, recipes2))
}

func repipesArEqual(recipes []db.Recipe, recipes2 []db.Recipe) bool {
	if len(recipes) != len(recipes2) {
		return false
	}

	repcipesById := make(map[string]db.Recipe)
	for i := range recipes {
		repcipesById[recipes[i].Id] = recipes[i]
	}

	for i := range recipes2 {
		recipe2 := recipes2[i]
		recipe, ok := repcipesById[recipe2.Id]
		if !ok {
			return false
		}

		if !recipeIngredientsAreEqual(recipe.Ingredients, recipe2.Ingredients) {
			return false
		}
	}

	return true
}

func recipeIngredientsAreEqual(ingredients, ingredients2 []db.RecipeIngredient) bool {
	if len(ingredients) != len(ingredients2) {
		return false
	}

	ingredientsById := make(map[string]db.RecipeIngredient)
	for i := range ingredients {
		ingredientsById[ingredients[i].IngredientFk] = ingredients[i]
	}

	for i := range ingredients2 {
		ingredient2 := ingredients2[i]
		ingredient, ok := ingredientsById[ingredient2.IngredientFk]
		if !ok {
			return false
		}

		if ingredient2.Amount != ingredient.Amount {
			return false
		} else if ingredient2.RecipeIdFk != ingredient.RecipeIdFk {
			return false
		}
	}

	return true
}
