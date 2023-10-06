package cocktailsdb

import (
	"context"
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/google/uuid"
	"strings"
)

const (
	recipeIngredientsTable = "cocktails.recipe_ingredients"
	recipesTable           = "cocktails.recipes"

	idCol           = "id"
	directionsCol   = "directions"
	recipeIdFkCol   = "recipe_id_fk"
	ingredientFkCol = "ingredient_fk"
	amountCol       = "amount"

	recipeIngredientsQueryFmt = `
		SELECT cocktails.recipe_ingredients.*
		FROM cocktails.recipe_ingredients
		WHERE recipe_id_fk IN (
			SELECT we_need.recipe_id_fk FROM (
				SELECT recipe_id_fk, COUNT(recipe_id_fk) AS cnt
				FROM cocktails.recipe_ingredients
				GROUP BY recipe_id_fk
			) we_need
			JOIN (
				SELECT recipe_id_fk, COUNT(recipe_id_fk) AS cnt
				FROM cocktails.recipe_ingredients ri
				WHERE ri.ingredient_fk IN (%s)
				GROUP BY recipe_id_fk
			) we_have
			ON we_have.recipe_id_fk = we_need.recipe_id_fk and we_have.cnt = we_need.cnt
		);`

	recipesQueryFmt = `
		SELECT *
		FROM cocktails.recipes
		WHERE id IN (%s);`
)

// RecipeIngredient represents a recipe ingredient.
type RecipeIngredient struct {
	RecipeIdFk   string  `db:"recipe_id_fk"`
	IngredientFk string  `db:"ingredient_fk"`
	Amount       float64 `db:"amount"`
}

type Recipe struct {
	Id          string  `db:"id"`
	CocktailFk  string  `db:"cocktail_fk"`
	DisplayName string  `db:"display_name"`
	Description *string `db:"description"`
	Directions  *string `db:"directions"`

	Ingredients []RecipeIngredient
}

// String returns a string representation of a RecipeIngredient.
func (r *RecipeIngredient) String() string {
	return strings.Join([]string{
		r.RecipeIdFk,
		r.IngredientFk,
		fmt.Sprintf("%f", r.Amount),
	}, ", ")
}

// CreateRecipe creates a recipe.
func CreateRecipe(ctx context.Context, tx *dbr.Tx, recipe *Recipe) error {
	if recipe.Id != "" {
		return fmt.Errorf("recipe id must be empty")
	} else if len(recipe.Ingredients) == 0 {
		return fmt.Errorf("recipe must have at least one ingredient")
	}

	recipe.Id = uuid.New().String()

	for i := range recipe.Ingredients {
		if recipe.Ingredients[i].RecipeIdFk != "" {
			return fmt.Errorf("recipe ingredient recipe id must be empty")
		} else if recipe.Ingredients[i].Amount <= 0 {
			return fmt.Errorf("recipe ingredient must have amount")
		}

		recipe.Ingredients[i].RecipeIdFk = recipe.Id
	}

	r, err := tx.InsertInto(recipesTable).Columns(idCol, displayNameCol, cocktailsFkCol, descriptionCol, directionsCol).Record(recipe).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert recipe: %w", err)
	}

	n, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if n != 1 {
		return fmt.Errorf("expected 1 row affected, got %d", n)
	}

	ins := tx.InsertInto(recipeIngredientsTable).Columns(recipeIdFkCol, ingredientFkCol, amountCol)
	for _, ingredient := range recipe.Ingredients {
		ins = ins.Record(&ingredient)
	}

	r, err = ins.ExecContext(ctx)
	if err != nil {
		return err
	}

	n, err = r.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	} else if n != int64(len(recipe.Ingredients)) {
		return fmt.Errorf("expected %d rows affected, got %d", len(recipe.Ingredients), n)
	}

	return nil
}

// mapRecipeIngredients maps a slice of RecipeIngredients to a map of recipe name to a slice of RecipeIngredients.
func mapRecipeIngredients(recipeIngredients []RecipeIngredient) map[string][]RecipeIngredient {
	recipeIdToIngredients := make(map[string][]RecipeIngredient)
	for _, ri := range recipeIngredients {
		recipeIdToIngredients[ri.RecipeIdFk] = append(recipeIdToIngredients[ri.RecipeIdFk], ri)
	}

	return recipeIdToIngredients
}

// ingListToMap takes a list of ingredients and turns it into a map of ingredient name to ingredient.
func ingListToMap(ingredients []RecipeIngredient) map[string]RecipeIngredient {
	ingMap := make(map[string]RecipeIngredient)
	for i := range ingredients {
		ingMap[ingredients[i].IngredientFk] = ingredients[i]
	}

	return ingMap
}

// GetRecipes returns all recipes, even if they are not available.
func GetRecipes(ctx context.Context, tx *dbr.Tx) ([]Recipe, error) {
	var recipes []Recipe
	_, err := tx.Select("*").From(recipesTable).LoadContext(ctx, &recipes)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes: %w", err)
	}

	var recipeIngredients []RecipeIngredient
	_, err = tx.Select("*").From(recipeIngredientsTable).LoadContext(ctx, &recipeIngredients)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe ingredients: %w", err)
	}

	nameToRI := mapRecipeIngredients(recipeIngredients)
	for i, recipe := range recipes {
		recipes[i].Ingredients = nameToRI[recipe.Id]
	}

	return recipes, nil
}

func quotedJoin(strs []string) string {
	quoted := make([]string, len(strs))
	for i, str := range strs {
		quoted[i] = fmt.Sprintf("'%s'", str)
	}

	return strings.Join(quoted, ",")
}

func quotedJoinRecipeIngMap(strSet map[string][]RecipeIngredient) string {
	strs := make([]string, len(strSet))
	i := 0
	for str := range strSet {
		strs[i] = fmt.Sprintf("'%s'", str)
		i++
	}

	return strings.Join(strs, ",")
}

// GetRecipesForIngredients returns all recipes that are available.
func GetRecipesForIngredients(ctx context.Context, tx *dbr.Tx, ingredients []string) ([]Recipe, error) {
	recipeIngredientsQuery := fmt.Sprintf(recipeIngredientsQueryFmt, quotedJoin(ingredients))

	var recipeIngredients []RecipeIngredient
	_, err := tx.SelectBySql(recipeIngredientsQuery).LoadContext(ctx, &recipeIngredients)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe ingredients: %w", err)
	}

	var recipes []Recipe
	if len(recipeIngredients) > 0 {
		recipeIdToIng := mapRecipeIngredients(recipeIngredients)
		recipesQuery := fmt.Sprintf(recipesQueryFmt, quotedJoinRecipeIngMap(recipeIdToIng))

		_, err = tx.SelectBySql(recipesQuery).LoadContext(ctx, &recipes)
		if err != nil {
			return nil, fmt.Errorf("failed to get recipes: %w", err)
		}

		for i := range recipes {
			recipes[i].Ingredients = recipeIdToIng[recipes[i].Id]
		}
	}

	return recipes, nil
}

// GetRecipesById returns all recipes with the given ids.
func GetRecipesById(ctx context.Context, tx *dbr.Tx, ids ...string) ([]Recipe, error) {
	var recipes []Recipe
	_, err := tx.Select("*").From(recipesTable).Where(dbr.Eq(idCol, ids)).LoadContext(ctx, &recipes)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes: %w", err)
	}

	var recipeIngredients []RecipeIngredient
	_, err = tx.Select("*").From(recipeIngredientsTable).Where(dbr.Eq(recipeIdFkCol, ids)).LoadContext(ctx, &recipeIngredients)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipe ingredients: %w", err)
	}

	recipeIdToIng := mapRecipeIngredients(recipeIngredients)
	for i := range recipes {
		recipes[i].Ingredients = recipeIdToIng[recipes[i].Id]
	}

	return recipes, nil
}

// DeleteRecipes deletes the recipes with the given names.
func DeleteRecipes(ctx context.Context, tx *dbr.Tx, name ...string) error {
	_, err := tx.DeleteFrom(recipeIngredientsTable).Where(dbr.Eq(recipeIdFkCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete recipe ingredients: %w", err)
	}

	_, err = tx.DeleteFrom(recipesTable).Where(dbr.Eq(idCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete recipes: %w", err)
	}

	return nil
}

// UpdateRecipe updates a recipe and the associated recipe ingredients.
func UpdateRecipe(ctx context.Context, tx *dbr.Tx, recipe *Recipe) error {
	_, err := tx.Update(recipesTable).SetMap(map[string]interface{}{
		displayNameCol: recipe.DisplayName,
		descriptionCol: recipe.Description,
		directionsCol:  recipe.Directions,
	}).Where(dbr.Eq(idCol, recipe.Id)).ExecContext(ctx)

	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}

	// get curreent recipe ingredients
	var currentRecipeIngredients []RecipeIngredient
	sel := tx.Select("*").From(recipeIngredientsTable).Where(dbr.Eq(recipeIdFkCol, recipe.Id))
	_, err = sel.LoadContext(ctx, &currentRecipeIngredients)
	if err != nil {
		return fmt.Errorf("failed to get current recipe ingredients: %w", err)
	}

	idToIngForNew := ingListToMap(recipe.Ingredients)
	idToIngForCurr := ingListToMap(currentRecipeIngredients)

	var changed []RecipeIngredient
	var added []RecipeIngredient

	// loop over the ingredients from the new recipe comparing them to the current recipe ingredients.
	// if the ingredient is not found in the current recipe ingredients, it is added.
	// if the ingredient is found in the current recipe ingredients, but the amount differs, it is changed.
	// if the ingredient is found in the current recipe ingredients, and the amount is the same, it is ignored.
	// if there is an ingredient in the current ingredients that is not found in the new recipe, it is deleted.
	for i := range recipe.Ingredients {
		name := strings.ToLower(recipe.Ingredients[i].IngredientFk)
		currIng, found := idToIngForCurr[name]

		if !found {
			added = append(added, recipe.Ingredients[i])
		} else if currIng.Amount != recipe.Ingredients[i].Amount {
			changed = append(changed, recipe.Ingredients[i])
		} // if found but doesn't differ ignore
	}

	// create lists of added, changed, and deleted ingredients
	deletedNames := make([]string, 0, len(currentRecipeIngredients))
	for i := range currentRecipeIngredients {
		name := strings.ToLower(currentRecipeIngredients[i].IngredientFk)
		_, found := idToIngForNew[name]
		if !found {
			deletedNames = append(deletedNames, name)
		}
	}

	// if there are items to delete, delete them in the db
	if len(deletedNames) > 0 {
		res, err := tx.DeleteFrom(recipeIngredientsTable).Where(
			dbr.And(dbr.Eq(recipeIdFkCol, recipe.Id), dbr.Eq(ingredientFkCol, deletedNames)),
		).ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to delete recipe ingredients: %w", err)
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		} else if affected != int64(len(deletedNames)) {
			return fmt.Errorf("expected %d rows to be deleted, but %d were", len(deletedNames), affected)
		}
	}

	// if there are items to add, add them in the db
	if len(added) > 0 {
		ins := tx.InsertInto(recipeIngredientsTable).Columns(recipeIdFkCol, ingredientFkCol, amountCol)
		for i := range added {
			ins = ins.Record(&added[i])
		}
		res, err := ins.ExecContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to insert recipe ingredients: %w", err)
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		} else if affected != int64(len(added)) {
			return fmt.Errorf("expected %d rows to be added, but %d were", len(added), affected)
		}
	}

	// if there are items to change, change them in the db
	if len(changed) > 0 {
		for i := range changed {
			_, err := tx.Update(recipeIngredientsTable).Set(amountCol, changed[i].Amount).Where(
				dbr.And(dbr.Eq(recipeIdFkCol, recipe.Id), dbr.Eq(ingredientFkCol, changed[i].IngredientFk)),
			).ExecContext(ctx)
			if err != nil {
				return fmt.Errorf("failed to update recipe ingredients: %w", err)
			}
		}
	}

	return nil
}
