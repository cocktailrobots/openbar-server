package openbardb

import (
	"context"
	"errors"
	"github.com/gocraft/dbr/v2"
)

func (s *testSuite) TestMenus() {
	s.Run("GetMenuNames", s.testGetMenuNames)
	s.Run("MenuIngredients", s.testMenuIngredients)
	s.Run("DeleteMenu", s.testDeleteMenu)
	s.Run("MenuItems", s.testMenuItems)
}

func (s *testSuite) testGetMenuNames() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	menuNames, err := GetMenuNames(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(menuNames, 0)

	err = CreateMenu(ctx, tx, "test_menu", []string{"gin", "rye"})
	s.Require().NoError(err)

	menuNames, err = GetMenuNames(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(menuNames, 1)
	s.Require().Equal("test_menu", menuNames[0])
}

func (s *testSuite) testMenuIngredients() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	menu, err := GetMenu(ctx, tx, "test_menu")
	s.Require().True(errors.Is(err, dbr.ErrNotFound))
	s.Require().Nil(menu)

	err = CreateMenu(ctx, tx, "test_menu", []string{"gin", "rye"})
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().Equal("test_menu", menu.Name)
	s.Require().Len(menu.Ingredients, 2)
	s.Require().Equal("gin", menu.Ingredients[0])
	s.Require().Equal("rye", menu.Ingredients[1])
}

func (s *testSuite) testDeleteMenu() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	err = CreateMenu(ctx, tx, "test_menu", []string{"gin", "rye"})
	s.Require().NoError(err)

	menu, err := GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().NotNil(menu)
	s.Require().Equal("test_menu", menu.Name)
	s.Require().Len(menu.Ingredients, 2)
	s.Require().Len(menu.RecipeIds, 0)

	err = DeleteMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().True(errors.Is(err, dbr.ErrNotFound))
	s.Require().Nil(menu)
}

func (s *testSuite) testMenuItems() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	err = CreateMenu(ctx, tx, "test_menu", []string{"gin", "rye"})
	s.Require().NoError(err)

	menu, err := GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().NotNil(menu)
	s.Require().Equal("test_menu", menu.Name)
	s.Require().Len(menu.Ingredients, 2)
	s.Require().Len(menu.RecipeIds, 0)

	err = AddMenuItem(ctx, tx, "test_menu", "1")
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().Len(menu.RecipeIds, 1)

	err = AddMenuItem(ctx, tx, "test_menu", "2")
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().Len(menu.RecipeIds, 2)

	err = RemoveMenuItem(ctx, tx, "test_menu", "noexist")
	s.Require().Error(err)

	err = RemoveMenuItem(ctx, tx, "test_menu", "1")
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().Len(menu.RecipeIds, 1)

	err = DeleteMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().True(errors.Is(err, dbr.ErrNotFound))
	s.Require().Nil(menu)

	err = AddMenuItem(ctx, tx, "test_menu", "1")
	s.Require().Error(err)

	err = RemoveMenuItem(ctx, tx, "test_menu", "1")
	s.Require().Error(err)

	err = CreateMenu(ctx, tx, "test_menu", []string{"gin", "rye", "compari", "vermouth"})
	s.Require().NoError(err)

	menu, err = GetMenu(ctx, tx, "test_menu")
	s.Require().NoError(err)
	s.Require().Len(menu.Ingredients, 4)
	s.Require().Len(menu.RecipeIds, 0)
}
