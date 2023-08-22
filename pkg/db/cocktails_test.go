package db

import (
	"context"
	"github.com/cocktailrobots/openbar-server/pkg/util"
)

const (
	initialCocktailsInTestDB = 6
)

func (s *testSuite) TestCocktails() {
	s.Run("GetCocktails", func() {
		s.Run("Get All Cocktails", s.testGetAllCocktails)
		s.Run("Get Cocktails With Names", s.testGetCocktailsWithNames)
	})
	s.Run("CreateCocktail", func() {
		s.Run("Create Valid Cocktails", s.testCreateCocktail)
		s.Run("Test Cocktail Uniqueness", s.testCocktailUniqueness)
	})
	s.Run("NormalizeCocktail", s.testNormalizeCocktail)
	s.Run("DeleteCocktail", s.testDeleteCocktail)
	s.Run("UpdateCocktail", s.testUpdateCocktail)
}

func (s *testSuite) testGetAllCocktails() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx, CocktailsDB)
	s.Require().NoError(err)

	cocktails, err := GetCocktails(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(cocktails, initialCocktailsInTestDB)
}

func (s *testSuite) testGetCocktailsWithNames() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx, CocktailsDB)
	s.Require().NoError(err)

	cocktails, err := GetCocktailsWithNames(ctx, tx, "noexist")
	s.Require().NoError(err)
	s.Require().Len(cocktails, 0)

	cocktails, err = GetCocktailsWithNames(ctx, tx, "boulevardier", "noexist")
	s.Require().NoError(err)
	s.Require().Len(cocktails, 1)

	// test dupblicate names
	cocktails, err = GetCocktailsWithNames(ctx, tx, "boulevardier", "boulevardier")
	s.Require().NoError(err)
	s.Require().Len(cocktails, 1)

	cocktails, err = GetCocktailsWithNames(ctx, tx, "boulevardier", "americano")
	s.Require().NoError(err)
	s.Require().Len(cocktails, 2)
}

func (s *testSuite) testCreateCocktail() {
	ctx := context.Background()

	cocktails := []Cocktail{
		{
			Name:        "cocktail1",
			DisplayName: "Cocktail One",
			Description: util.Ptr("A cocktail"),
		},
		{
			Name:        "cocktail2",
			DisplayName: "Cocktail Two",
		},
	}

	tx, err := s.BeginTx(ctx, CocktailsDB)
	s.Require().NoError(err)

	err = AddCocktails(ctx, tx, cocktails...)
	s.Require().NoError(err)

	err = tx.Commit()
	s.Require().NoError(err)

	tx, err = s.BeginTx(ctx, CocktailsDB)
	s.Require().NoError(err)

	retrieved, err := GetCocktailsWithNames(ctx, tx, "cocktail1", "cocktail2")
	s.Require().NoError(err)

	s.Require().Equal(cocktails, retrieved)
}

func (s *testSuite) testCocktailUniqueness() {
	ctx := context.Background()

	cocktails := []Cocktail{
		{
			Name:        "cocktail1",
			DisplayName: "Cocktail One",
		},
		{
			Name:        "CoCkTaIl1",
			DisplayName: "Cocktail One Again",
		},
	}

	tx, err := s.BeginTx(ctx, CocktailsDB)
	s.Require().NoError(err)

	err = AddCocktails(ctx, tx, cocktails...)
	s.Require().Error(err)
}

func (s *testSuite) testNormalizeCocktail() {
	tests := []struct {
		name string
		in   Cocktail
		out  Cocktail
	}{
		{
			name: "No Display Name",
			in:   Cocktail{Name: ""},
			out:  Cocktail{Name: "", DisplayName: ""},
		},
		{
			name: "Display Name",
			in:   Cocktail{Name: "test", DisplayName: "Test"},
			out:  Cocktail{Name: "test", DisplayName: "Test"},
		},
		{
			name: "Name with dashes and underscores",
			in:   Cocktail{Name: "part1-part2_part3"},
			out:  Cocktail{Name: "part1-part2_part3", DisplayName: "Part1 Part2 Part3"},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			in := test.in
			normalizeCocktail(&in)
			s.Require().Equal(test.out, in)
		})
	}
}

func (s *testSuite) testDeleteCocktail() {
	tests := []struct {
		name                string
		toDelete            []string
		expectedCountBefore int
		expecteCountAfter   int
		expectErr           bool
	}{
		{
			name:                "Delete Nonexistent Cocktail",
			toDelete:            []string{"noexist"},
			expectedCountBefore: initialCocktailsInTestDB,
			expecteCountAfter:   initialCocktailsInTestDB,
			expectErr:           true,
		},
		{
			name:                "Delete Single Cocktail",
			toDelete:            []string{"boulevardier"},
			expectedCountBefore: initialCocktailsInTestDB,
			expecteCountAfter:   initialCocktailsInTestDB - 1,
		},
		{
			name:                "Delete Multiple Cocktails",
			toDelete:            []string{"boulevardier", "americano"},
			expectedCountBefore: initialCocktailsInTestDB,
			expecteCountAfter:   initialCocktailsInTestDB - 2,
		},
		{
			name:                "Delete Same Cocktail Twice",
			toDelete:            []string{"boulevardier", "boulevardier"},
			expectedCountBefore: initialCocktailsInTestDB,
			expecteCountAfter:   initialCocktailsInTestDB - 1,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			ctx := context.Background()
			tx, err := s.BeginTx(ctx, CocktailsDB)
			s.Require().NoError(err)

			cocktails, err := GetCocktails(ctx, tx)
			s.Require().NoError(err)
			initialLen := len(cocktails)
			s.Require().Equal(test.expectedCountBefore, initialLen)

			err = DeleteCocktails(ctx, tx, test.toDelete...)
			if test.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}

			cocktails, err = GetCocktails(ctx, tx)
			s.Require().NoError(err)
			s.Require().Len(cocktails, test.expecteCountAfter)
		})
	}
}

func (s *testSuite) testUpdateCocktail() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx, CocktailsDB)
	s.Require().NoError(err)

	cocktails, err := GetCocktailsWithNames(ctx, tx, "boulevardier")
	s.Require().NoError(err)

	cocktail := cocktails[0]
	cocktail.DisplayName = "New Name"
	cocktail.Description = util.Ptr("New Description")

	err = UpdateCocktail(ctx, tx, &cocktail)
	s.Require().NoError(err)

	cocktails, err = GetCocktailsWithNames(ctx, tx, "boulevardier")
	s.Require().NoError(err)
	s.Require().Equal(cocktail, cocktails[0])
}
