package wire

import (
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/db/cocktailsdb"
	"strings"
)

// Ingredient is an ingredient as it will be written to the wire in HTTP responses.
type Ingredient struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Description *string `json:"description"`
}

// Ingredients is a slice of ingredients.
type Ingredients []Ingredient

// ToDbIngredients converts a list of Ingredients to a list of cocktailsdb.Ingredients.
func (ing Ingredients) ToDbIngredients() []cocktailsdb.Ingredient {
	ingredients := make([]cocktailsdb.Ingredient, len(ing))
	for i := range ing {
		ingredients[i] = cocktailsdb.Ingredient{
			Name:        ing[i].Name,
			DisplayName: ing[i].DisplayName,
			Description: ing[i].Description,
		}
	}

	return ingredients
}

func (ing Ingredients) Validate() error {
	for i := range ing {
		ing[i].Name = strings.TrimSpace(ing[i].Name)
		if len(ing[i].Name) == 0 {
			return fmt.Errorf("ingredient name cannot be empty")
		}

		ing[i].DisplayName = strings.TrimSpace(ing[i].DisplayName)
		if len(ing[i].DisplayName) == 0 {
			return fmt.Errorf("ingredient display name cannot be empty")
		}
	}

	return nil
}

// FromDbIngredients converts a list of cocktailsdb.Ingredients to a list of Ingredients.
func FromDbIngredients(ingredients []cocktailsdb.Ingredient) Ingredients {
	ing := make(Ingredients, len(ingredients))
	for i := range ingredients {
		ing[i] = Ingredient{
			Name:        ingredients[i].Name,
			DisplayName: ingredients[i].DisplayName,
			Description: ingredients[i].Description,
		}
	}

	return ing
}
