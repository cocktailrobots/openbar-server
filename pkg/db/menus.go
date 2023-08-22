package db

import (
	"context"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/util/functional"
	"github.com/gocraft/dbr/v2"
)

const (
	MenusTable           = "menus"
	MenuItemsTable       = "menu_items"
	MenuIngredientsTable = "menu_ingredients"

	ingredientNameCol = "ingredient_name"
	recipeIdCol       = "recipe_id"
	menuNameFkCol     = "menu_name_fk"
)

type MenuItem struct {
	MenuNameFk string `db:"menu_name_fk"`
	RecipeIdFk string `db:"recipe_id"`
}

type MenuIngredient struct {
	MenuNameFk     string `db:"menu_name_fk"`
	IngredientName string `db:"ingredient_name"`
}

type Menu struct {
	Name        string
	RecipeIds   []string
	Ingredients []string
}

func GetMenuNames(ctx context.Context, tx *dbr.Tx) ([]string, error) {
	var menus []string
	_, err := tx.Select(nameCol).From(MenusTable).LoadContext(ctx, &menus)
	if err != nil {
		return nil, err
	}

	return menus, nil
}

func GetMenu(ctx context.Context, tx *dbr.Tx, name string) (*Menu, error) {
	var menuName string
	err := tx.Select(nameCol).From(MenusTable).Where(dbr.Eq(nameCol, name)).LoadOneContext(ctx, &menuName)
	if err != nil {
		return nil, fmt.Errorf("failed to load menu '%s': %w", name, err)
	} else if menuName != name {
		return nil, fmt.Errorf("failed to load menu '%s': not found", name)
	}

	var ingredients []MenuIngredient
	_, err = tx.Select("*").From(MenuIngredientsTable).Where(dbr.Eq(menuNameFkCol, name)).OrderBy(ingredientNameCol).LoadContext(ctx, &ingredients)
	if err != nil {
		return nil, fmt.Errorf("failed to load ingredients for menu '%s': %w", name, err)
	}

	var items []MenuItem
	_, err = tx.Select("*").From(MenuItemsTable).Where(dbr.Eq(menuNameFkCol, name)).LoadContext(ctx, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to load menu items for menu '%s': %w", name, err)
	}

	ingredientNames := functional.Map(func(item MenuIngredient) string { return item.IngredientName }, ingredients)
	recipeIds := functional.Map(func(item MenuItem) string { return item.RecipeIdFk }, items)
	return &Menu{
		Name:        name,
		RecipeIds:   recipeIds,
		Ingredients: ingredientNames,
	}, nil
}

func CreateMenu(ctx context.Context, tx *dbr.Tx, name string, ingredients []string) error {
	_, err := tx.InsertInto(MenusTable).Columns(nameCol).Values(name).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert menu '%s': %w", name, err)
	}

	ins := tx.InsertInto(MenuIngredientsTable).Columns(menuNameFkCol, ingredientNameCol)
	for i := range ingredients {
		item := MenuIngredient{
			MenuNameFk:     name,
			IngredientName: ingredients[i],
		}

		ins = ins.Record(&item)
	}

	_, err = ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert ingredients for menu '%s': %w", name, err)
	}

	return nil
}

func DeleteMenu(ctx context.Context, tx *dbr.Tx, name string) error {
	_, err := tx.DeleteFrom(MenusTable).Where(dbr.Eq(nameCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete menu '%s': %w", name, err)
	}

	_, err = tx.DeleteFrom(MenuIngredientsTable).Where(dbr.Eq(menuNameFkCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete ingredients for menu '%s': %w", name, err)
	}

	_, err = tx.DeleteFrom(MenuItemsTable).Where(dbr.Eq(menuNameFkCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete menu items for menu '%s': %w", name, err)
	}

	return nil
}

func AddMenuItem(ctx context.Context, tx *dbr.Tx, name, recipeId string) error {
	_, err := tx.InsertInto(MenuItemsTable).Columns(menuNameFkCol, "recipe_id").Values(name, recipeId).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert menu item '%s' into menu '%s': %w", recipeId, name, err)
	}

	return nil
}

func RemoveMenuItem(ctx context.Context, tx *dbr.Tx, name, recipeId string) error {
	res, err := tx.DeleteFrom(MenuItemsTable).Where(dbr.And(dbr.Eq(menuNameFkCol, name), dbr.Eq(recipeIdCol, recipeId))).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove menu item '%s' from menu '%s': %w", recipeId, name, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows for menu item '%s' from menu '%s': %w", recipeId, name, err)
	} else if n == 0 {
		return dbr.ErrNotFound
	}

	return nil
}

func UpdateMenu(ctx context.Context, tx *dbr.Tx, name string, ingredients, recipeIds []string) error {
	_, err := tx.DeleteFrom(MenuIngredientsTable).Where(dbr.Eq(menuNameFkCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete ingredients for menu '%s': %w", name, err)
	}

	_, err = tx.DeleteFrom(MenuItemsTable).Where(dbr.Eq(menuNameFkCol, name)).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete menu items for menu '%s': %w", name, err)
	}

	ins := tx.InsertInto(MenuIngredientsTable).Columns(menuNameFkCol, ingredientNameCol)
	for i := range ingredients {
		item := MenuIngredient{
			MenuNameFk:     name,
			IngredientName: ingredients[i],
		}

		ins = ins.Record(&item)
	}

	_, err = ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert ingredients for menu '%s': %w", name, err)
	}

	ins = tx.InsertInto(MenuItemsTable).Columns(menuNameFkCol, recipeIdCol)
	for i := range recipeIds {
		item := MenuItem{
			MenuNameFk: name,
			RecipeIdFk: recipeIds[i],
		}

		ins = ins.Record(&item)
	}

	_, err = ins.ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to insert menu items for menu '%s': %w", name, err)
	}

	return nil
}
