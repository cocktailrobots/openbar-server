package cocktailsdb

import (
	"context"
	"github.com/cocktailrobots/openbar-server/pkg/util"
	"github.com/gocraft/dbr/v2"
	"log"
	"strings"
)

const (
	ingredientsTable = "cocktails.ingredients"
)

// Ingredient represents an ingredient in the database.
type Ingredient struct {
	Name        string  `db:"name"`
	DisplayName string  `db:"display_name"`
	Description *string `db:"description"`
}

// GetIngredients returns all ingredients.
func GetIngredients(ctx context.Context, tx *dbr.Tx) ([]Ingredient, error) {
	var ingredients []Ingredient
	_, err := tx.Select("*").From(ingredientsTable).LoadContext(ctx, &ingredients)
	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

func normalizeIngredient(ingredient *Ingredient) {
	ingredient.Name = strings.TrimSpace(ingredient.Name)
	ingredient.DisplayName = strings.TrimSpace(ingredient.DisplayName)
	if len(ingredient.DisplayName) == 0 {
		ingredient.DisplayName = util.TitleCase(ingredient.Name)
	}
}

// GetIngredientsWithNames returns the ingredients with the given names.
func GetIngredientsWithNames(ctx context.Context, tx *dbr.Tx, names ...string) ([]Ingredient, error) {
	var ingredients []Ingredient
	sel := tx.Select("*").From(ingredientsTable).Where(dbr.Eq(nameCol, names))
	_, err := sel.LoadContext(ctx, &ingredients)

	if err != nil {
		return nil, err
	}

	return ingredients, nil
}

// CreateIngredients adds ingredients to the database.
func CreateIngredients(ctx context.Context, tx *dbr.Tx, ingredients ...Ingredient) error {
	for i := range ingredients {
		normalizeIngredient(&ingredients[i])
	}

	ins := tx.InsertInto(ingredientsTable).Columns(nameCol, displayNameCol, descriptionCol)
	for _, ingredient := range ingredients {
		ins = ins.Record(&ingredient)
	}

	_, err := ins.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

// DeleteIngredients deletes ingredients from the database. If an ingredient is used in a recipe this will fail.
func DeleteIngredients(ctx context.Context, tx *dbr.Tx, ingredients ...string) error {
	_, err := tx.DeleteFrom(ingredientsTable).Where(dbr.Eq(nameCol, ingredients)).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

// UpdateIngredient updates an ingredient in the database.
func UpdateIngredient(ctx context.Context, tx *dbr.Tx, ingredient Ingredient) error {
	normalizeIngredient(&ingredient)

	log.Println(ingredient)
	upd := tx.Update(ingredientsTable).SetMap(map[string]interface{}{
		descriptionCol: ingredient.Description,
		displayNameCol: ingredient.DisplayName,
	}).Where(dbr.Eq(nameCol, ingredient.Name))

	_, err := upd.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
