package cocktailsapi

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"net/http"
)

const (
	ingredientsInTestDB = 11
)

// ingredients
// GET /ingredients
// GET /ingredients/{id}
// POST /ingredients
// PUT /ingredients/{id}
// DELETE /ingredients/{id}

func (s *testSuite) TestIngredients() {
	s.Run("ListIngredients", s.testListIngredients)
	s.Run("GetIngredient", s.testGetIngredient)
	s.Run("PostIngredient", s.testPostIngredient)
	s.Run("PatchIngredient", s.testPatchIngredient)
	s.Run("DeleteIngredient", s.testDeleteIngredient)
}

func (s *testSuite) testListIngredients() {
	req, err := http.NewRequest(http.MethodGet, "/ingredients", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var ingredients []wire.Ingredient
	err = json.Unmarshal(respWr.Body(), &ingredients)
	s.Require().NoError(err)
	s.Require().Len(ingredients, ingredientsInTestDB)
}

func (s *testSuite) testGetIngredient() {
	req, err := http.NewRequest(http.MethodGet, "/ingredients/noexist", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodGet, "/ingredients/gin", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var ingredient wire.Ingredient
	err = json.Unmarshal(respWr.Body(), &ingredient)
	s.Require().NoError(err)
	s.Require().Equal("gin", ingredient.Name)
}

func (s *testSuite) testPostIngredient() {
	ingredients := []wire.Ingredient{
		{
			Name:        "test1",
			DisplayName: "Test",
			Description: ptr("Test ingredient description"),
		}, {
			Name:        "test2",
			DisplayName: "Test",
			Description: ptr("Test2 ingredient description"),
		},
	}

	req, err := http.NewRequest(http.MethodPost, "/ingredients", test.JsonReaderForObject(ingredients))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodGet, "/ingredients", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var ingredientsResp []wire.Ingredient
	err = json.Unmarshal(respWr.Body(), &ingredientsResp)
	s.Require().NoError(err)
	s.Require().Len(ingredientsResp, ingredientsInTestDB+len(ingredients))
}

func (s *testSuite) testPatchIngredient() {
	ingredient := wire.Ingredient{
		Name:        "gin",
		DisplayName: "Test Patch",
		Description: ptr("Test patch ingredient description"),
	}

	req, err := http.NewRequest(http.MethodPatch, "/ingredients/gin", test.JsonReaderForObject(ingredient))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodGet, "/ingredients/gin", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var ingredientResp wire.Ingredient
	err = json.Unmarshal(respWr.Body(), &ingredientResp)
	s.Require().NoError(err)
	s.Require().Equal(ingredient, ingredientResp)
}

func (s *testSuite) testDeleteIngredient() {
	req, err := http.NewRequest(http.MethodGet, "/ingredients/dry_vermouth", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodDelete, "/ingredients/dry_vermouth", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodGet, "/ingredients/dry_vermouth", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())
}
