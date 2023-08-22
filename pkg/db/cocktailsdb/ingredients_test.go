package cocktailsdb

import (
	"context"
	"github.com/cocktailrobots/openbar-server/pkg/util"
)

const (
	initialIngredientsInTestDB = 11
)

func (s *testSuite) TestIngredients() {
	s.Run("GetIngredients", func() {
		s.Run("Get All Ingredients", s.testGetAllIngredients)
		s.Run("Get Ingredients With Names", s.testGetIngredientsWithNames)
	})
	s.Run("CreateIngredient", func() {
		s.Run("Create Valid Ingredients", s.testCreateIngredient)
	})
	s.Run("DeleteIngredient", s.testDeleteIngredients)
	s.Run("UpdateIngredient", s.testUpdateIngredient)
}

func (s *testSuite) testGetAllIngredients() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	ingredients, err := GetIngredients(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(ingredients, initialIngredientsInTestDB)
}

func (s *testSuite) testGetIngredientsWithNames() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	ingredients, err := GetIngredientsWithNames(ctx, tx, "noexist")
	s.Require().NoError(err)
	s.Require().Len(ingredients, 0)

	ingredients, err = GetIngredientsWithNames(ctx, tx, "gin")
	s.Require().NoError(err)
	s.Require().Len(ingredients, 1)

	// test dupblicate names
	ingredients, err = GetIngredientsWithNames(ctx, tx, "gin", "gin")
	s.Require().NoError(err)
	s.Require().Len(ingredients, 1)

	ingredients, err = GetIngredientsWithNames(ctx, tx, "gin", "rye")
	s.Require().NoError(err)
	s.Require().Len(ingredients, 2)
}

func (s *testSuite) testCreateIngredient() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	ingredient := Ingredient{
		Name: "new_ingredient",
	}

	err = CreateIngredients(ctx, tx, ingredient)
	s.Require().NoError(err)

	ingredients, err := GetIngredientsWithNames(ctx, tx, ingredient.Name)
	s.Require().NoError(err)
	s.Require().Len(ingredients, 1)
}

func (s *testSuite) testDeleteIngredients() {
	const unusedIngredient = "dry_vermouth"

	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	ingredients, err := GetIngredients(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(ingredients, initialIngredientsInTestDB)

	err = DeleteIngredients(ctx, tx, unusedIngredient)
	s.Require().NoError(err)

	ingredients, err = GetIngredients(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(ingredients, initialIngredientsInTestDB-1)
}

func (s *testSuite) testUpdateIngredient() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	ingredients, err := GetIngredientsWithNames(ctx, tx, "gin")
	s.Require().NoError(err)
	s.Require().Len(ingredients, 1)

	ingredients[0].DisplayName = "new_display_name"
	ingredients[0].Description = util.Ptr("new_description")

	err = UpdateIngredient(ctx, tx, ingredients[0])
	s.Require().NoError(err)

	updated, err := GetIngredientsWithNames(ctx, tx, ingredients[0].Name)
	s.Require().NoError(err)
	s.Require().Len(ingredients, 1)
	s.Require().Equal(ingredients[0].DisplayName, updated[0].DisplayName)
	s.Require().Equal(ingredients[0].Description, updated[0].Description)
}
