package wire

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/db/cocktailsdb"
)

// RecipeIngredient is a single ingredient in a recipe.
type RecipeIngredient struct {
	Name   string  `db:"name"`
	Amount float64 `db:"amount"`
}

// Recipe is a recipe.
type Recipe struct {
	Id          string             `json:"id"`
	CocktailId  string             `json:"cocktail_id"`
	DisplayName string             `json:"display_name"`
	Description string             `json:"description"`
	Directions  string             `json:"directions"`
	Ingredients []RecipeIngredient `json:"ingredients"`
}

type Recipes []Recipe

func (r Recipes) ToDbRecipes() []cocktailsdb.Recipe {
	recipes := make([]cocktailsdb.Recipe, 0, len(r))
	for _, recipe := range r {
		ingredients := make([]cocktailsdb.RecipeIngredient, len(recipe.Ingredients))
		for i, ingredient := range recipe.Ingredients {
			ingredients[i] = cocktailsdb.RecipeIngredient{
				RecipeIdFk:   recipe.Id,
				IngredientFk: ingredient.Name,
				Amount:       ingredient.Amount,
			}
		}

		var description *string
		if recipe.Description != "" {
			description = &recipe.Description
		}

		var directions *string
		if recipe.Directions != "" {
			directions = &recipe.Directions
		}

		recipes = append(recipes, cocktailsdb.Recipe{
			Id:          recipe.Id,
			DisplayName: recipe.DisplayName,
			CocktailFk:  recipe.CocktailId,
			Description: description,
			Directions:  directions,
			Ingredients: ingredients,
		})
	}

	return recipes
}

func (r Recipes) Validate() error {
	for _, recipe := range r {
		if len(recipe.Id) == 0 {
			return fmt.Errorf("recipe id cannot be empty")
		}

		if len(recipe.DisplayName) == 0 {
			return fmt.Errorf("recipe display name cannot be empty")
		}

		for _, ingredient := range recipe.Ingredients {
			if len(ingredient.Name) == 0 {
				return fmt.Errorf("ingredient name cannot be empty")
			} else if ingredient.Amount <= 0 {
				return fmt.Errorf("ingredient amount must be greater than 0")
			}
		}
	}

	return nil
}

func FromDbRecipes(recipes []cocktailsdb.Recipe) Recipes {
	r := make(Recipes, 0, len(recipes))
	for _, recipe := range recipes {
		ingredients := make([]RecipeIngredient, len(recipe.Ingredients))
		for j, ingredient := range recipe.Ingredients {
			ingredients[j] = RecipeIngredient{
				Name:   ingredient.IngredientFk,
				Amount: ingredient.Amount,
			}
		}

		var description string
		if recipe.Description != nil {
			description = *recipe.Description
		}

		var directions string
		if recipe.Directions != nil {
			directions = *recipe.Directions
		}

		r = append(r, Recipe{
			Id:          recipe.Id,
			CocktailId:  recipe.CocktailFk,
			DisplayName: recipe.DisplayName,
			Description: description,
			Directions:  directions,
			Ingredients: ingredients,
		})
	}

	return r
}
