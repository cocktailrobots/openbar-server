package cocktailsdb

import (
	"context"
	"github.com/cocktailrobots/openbar-server/pkg/util"
)

const (
	initialRecipesInTestDB = 6
)

func (s *testSuite) TestRecipes() {
	s.Run("GetRecipes", func() {
		s.Run("Get All Recipes", s.testGetAllRecipes)
		s.Run("Get Recipes With Names", s.testGetRecipesById)
	})
	s.Run("CreateRecipe", func() {
		s.Run("Create Valid Recipes", s.testCreateRecipe)
	})
	s.Run("DeleteRecipe", s.testDeleteRecipes)
	s.Run("GetRecipesForIngredients", s.testGetRecipesForIngredients)
	s.Run("UpdateRecipe", s.testUpdateRecipe)
}

func (s *testSuite) testGetAllRecipes() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	recipes, err := GetRecipes(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(recipes, initialRecipesInTestDB)
}

func (s *testSuite) testGetRecipesById() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	recipes, err := GetRecipesById(ctx, tx, "noexist")
	s.Require().NoError(err)
	s.Require().Len(recipes, 0)

	recipes, err = GetRecipes(ctx, tx)
	s.Require().NoError(err)

	ids := make([]string, len(recipes))
	for i := range recipes {
		ids[i] = recipes[i].Id
	}

	recipes, err = GetRecipesById(ctx, tx, ids[0], "noexist")
	s.Require().NoError(err)
	s.Require().Len(recipes, 1)

	// test dupblicate names
	recipes, err = GetRecipesById(ctx, tx, ids[0], ids[0])
	s.Require().NoError(err)
	s.Require().Len(recipes, 1)

	recipes, err = GetRecipesById(ctx, tx, ids...)
	s.Require().NoError(err)
	s.Require().Len(recipes, len(ids))
}

func (s *testSuite) testCreateRecipe() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	cocktail := Cocktail{
		Name: "martini",
	}

	err = AddCocktails(ctx, tx, cocktail)
	s.Require().NoError(err)

	recipe := Recipe{
		DisplayName: "Test",
		CocktailFk:  cocktail.Name,
		Ingredients: []RecipeIngredient{
			{
				IngredientFk: "gin",
				Amount:       1,
			},
			{
				IngredientFk: "dry_vermouth",
				Amount:       1,
			},
		},
	}

	err = CreateRecipe(ctx, tx, &recipe)
	s.Require().NoError(err)

	recipes, err := GetRecipesById(ctx, tx, recipe.Id)
	s.Require().NoError(err)
	s.Require().Len(recipes, 1)
}

func (s *testSuite) testDeleteRecipes() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	recipes, err := GetRecipes(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(recipes, initialRecipesInTestDB)

	err = DeleteRecipes(ctx, tx, recipes[0].Id)
	s.Require().NoError(err)

	recipes, err = GetRecipes(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(recipes, initialRecipesInTestDB-1)
}

func (s *testSuite) testGetRecipesForIngredients() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	recipes, err := GetRecipesForIngredients(ctx, tx, []string{"bourbon", "vodka"})
	s.Require().NoError(err)
	s.Require().Len(recipes, 0)

	recipes, err = GetRecipesForIngredients(ctx, tx, []string{"bourbon", "campari", "sweet_vermouth"})
	s.Require().NoError(err)
	s.Require().Len(recipes, 3)

	displayNames := util.NewSet[string]()
	for i := range recipes {
		displayNames.Add(recipes[i].DisplayName)
	}

	s.Require().True(displayNames.HasOnly("Boulevardier", "Americano", "Manhattan"))
}

func (s *testSuite) testUpdateRecipe() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	recipes, err := GetRecipes(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(recipes, initialRecipesInTestDB)

	var boulevardierIdx int
	for i := range recipes {
		if recipes[i].DisplayName == "Boulevardier" {
			boulevardierIdx = i
			break
		}
	}

	recipe := recipes[boulevardierIdx]
	recipe.DisplayName = "NewName"
	recipe.Directions = util.Ptr("NewDirections")
	recipe.Description = util.Ptr("NewDescription")

	// change 1 ingredient amount, delete 2 ingredients, add 3
	recipe.Ingredients = []RecipeIngredient{
		{
			RecipeIdFk:   recipe.Id,
			IngredientFk: "bourbon",
			Amount:       3,
		},
		{
			RecipeIdFk:   recipe.Id,
			IngredientFk: "lemon_juice",
			Amount:       3,
		},
		{
			RecipeIdFk:   recipe.Id,
			IngredientFk: "amaro_nonino",
			Amount:       3,
		},
		{
			RecipeIdFk:   recipe.Id,
			IngredientFk: "aperol",
			Amount:       3,
		},
	}

	err = UpdateRecipe(ctx, tx, &recipe)
	s.Require().NoError(err)

	updatedRecipes, err := GetRecipesById(ctx, tx, recipe.Id)
	s.Require().NoError(err)

	updated := updatedRecipes[0]
	s.Require().Equal("NewName", updated.DisplayName)
	s.Require().Equal("NewDirections", *updated.Directions)
	s.Require().Equal("NewDescription", *updated.Description)
	s.Require().True(ingredientsEqual(recipe.Ingredients, updated.Ingredients))
}

func ingredientsEqual(expected, actual []RecipeIngredient) bool {
	if len(expected) != len(actual) {
		return false
	}

	expectedMap := make(map[string]RecipeIngredient)
	for i := range expected {
		expectedMap[expected[i].IngredientFk] = expected[i]
	}

	for i := range actual {
		if expected, ok := expectedMap[actual[i].IngredientFk]; !ok {
			return false
		} else if expected.Amount != actual[i].Amount {
			return false
		}
	}

	return true
}
