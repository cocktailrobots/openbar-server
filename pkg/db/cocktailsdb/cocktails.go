package cocktailsdb

import (
	"context"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/util"
	"github.com/gocraft/dbr/v2"
	"strings"
)

const (
	cocktailsTable = "cocktails.cocktails"

	nameCol        = "name"
	displayNameCol = "display_name"
	descriptionCol = "description"

	cocktailsFkCol = "cocktail_fk"
)

// Cocktail represents a cocktail in the database.
type Cocktail struct {
	Name        string  `db:"name"`
	DisplayName string  `db:"display_name"`
	Description *string `db:"description"`
}

// GetCocktails returns all cocktails.
func GetCocktails(ctx context.Context, tx *dbr.Tx) ([]Cocktail, error) {
	var cocktails []Cocktail
	_, err := tx.Select("*").From(cocktailsTable).OrderBy(nameCol).LoadContext(ctx, &cocktails)
	if err != nil {
		return nil, err
	}

	return cocktails, nil
}

// GetCocktailsWithNames returns the cocktails with the given names.
func GetCocktailsWithNames(ctx context.Context, tx *dbr.Tx, names ...string) ([]Cocktail, error) {
	var cocktails []Cocktail
	sel := tx.Select("*").From(cocktailsTable).Where(dbr.Eq(nameCol, names)).OrderBy(nameCol)
	_, err := sel.LoadContext(ctx, &cocktails)

	if err != nil {
		return nil, fmt.Errorf("failed to get cocktails with names %v: %w", names, err)
	}

	return cocktails, nil
}

// normalizeCocktail normalizes the cocktail's fields. This should be called before adding a cocktail to the database.
// It trims whitespace from the name and display name, and sets the display name to the title case of the name if it is
// empty.
func normalizeCocktail(cocktail *Cocktail) {
	cocktail.Name = strings.TrimSpace(cocktail.Name)
	cocktail.DisplayName = strings.TrimSpace(cocktail.DisplayName)
	if len(cocktail.DisplayName) == 0 {
		cocktail.DisplayName = util.ReplaceChars(cocktail.Name, map[rune]rune{'_': ' ', '-': ' '})
		cocktail.DisplayName = util.TitleCase(cocktail.DisplayName)
	}
}

// AddCocktails adds cocktails to the database.
func AddCocktails(ctx context.Context, tx *dbr.Tx, cocktails ...Cocktail) error {
	for i := range cocktails {
		normalizeCocktail(&cocktails[i])
	}

	ins := tx.InsertInto(cocktailsTable).Columns(nameCol, displayNameCol, descriptionCol)
	for _, cocktail := range cocktails {
		ins = ins.Record(&cocktail)
	}

	res, err := ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to add cocktails: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if rowsAffected != int64(len(cocktails)) {
		return fmt.Errorf("expected %d rows to be affected, but got %d", len(cocktails), rowsAffected)
	}

	return nil
}

// DeleteCocktails deletes cocktails from the database. This will fail if recipes referencing this cocktail aren't deleted first.
func DeleteCocktails(ctx context.Context, tx *dbr.Tx, names ...string) error {
	unique := util.NewSet(names...)
	names = unique.Items()

	var cocktails []Cocktail
	_, err := tx.Select("*").From(cocktailsTable).Where(dbr.Eq(nameCol, names)).LoadContext(ctx, &cocktails)
	if err != nil {
		return fmt.Errorf("failed to get cocktails: %w", err)
	} else if len(names) != len(cocktails) {
		found := util.NewSet[string]()
		for _, cocktail := range cocktails {
			found.Add(cocktail.Name)
		}

		for _, name := range names {
			if !found.Contains(name) {
				return fmt.Errorf("%s does not exist: %w", name, dbr.ErrNotFound)
			}
		}
	}

	var recipes []Recipe
	_, err = tx.Select("*").From(recipesTable).Where(dbr.Eq(cocktailsFkCol, names)).LoadContext(ctx, &recipes)
	if err != nil {
		return fmt.Errorf("failed to get recipes: %w", err)
	}

	if len(recipes) > 0 {
		var recipeIds []string
		for _, recipe := range recipes {
			recipeIds = append(recipeIds, recipe.Id)
		}

		_, err := tx.DeleteFrom(recipeIngredientsTable).Where(dbr.Eq(recipeIdFkCol, recipeIds)).ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete recipe ingredients: %w", err)
		}

		_, err = tx.DeleteFrom(recipesTable).Where(dbr.Eq(cocktailsFkCol, names)).ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete recipes: %w", err)
		}
	}

	res, err := tx.DeleteFrom(cocktailsTable).Where(dbr.Eq(nameCol, names)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete cocktails: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if rowsAffected != int64(unique.Len()) {
		return fmt.Errorf("expected 1 row to be affected, but got %d", rowsAffected)
	}

	return nil
}

// UpdateCocktail updates a cocktail in the database.
func UpdateCocktail(ctx context.Context, tx *dbr.Tx, cocktail *Cocktail) error {
	normalizeCocktail(cocktail)

	_, err := tx.Update(cocktailsTable).SetMap(map[string]interface{}{
		descriptionCol: cocktail.Description,
		displayNameCol: cocktail.DisplayName,
	}).Where(dbr.Eq(nameCol, cocktail.Name)).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("failed to update cocktail: %w", err)
	}

	return nil
}
