package openbarapi

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"net/http"
)

var menus = wire.Menus{
	{
		Name:        "test1",
		Ingredients: []string{"ing1", "ing2"},
		RecipeIds:   []string{},
	},
	{
		Name:        "test2",
		Ingredients: []string{"ing1", "ing2"},
		RecipeIds:   []string{},
	},
}

func (s *testSuite) TestMenus() {
	s.Run("create menus", s.testCreateMenus)
	s.Run("delete menus", s.testDeleteMenus)
	s.Run("get menu", s.testGetMenu)
	s.Run("add/remove recipes", s.testAddRemoveRecipes)
}

func (s *testSuite) testCreateMenus() {
	// starts with no menus

	// GET /menus should return empty list
	req, err := http.NewRequest(http.MethodGet, "/menus", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var menuNames []string
	err = json.Unmarshal(respWr.Body(), &menuNames)
	s.Require().NoError(err)
	s.Require().Len(menuNames, 0)

	// POST /menus should create two menus
	req, err = http.NewRequest(http.MethodPost, "/menus", test.JsonReaderForObject(menus))
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	// GET /menus should return the two menus
	req, err = http.NewRequest(http.MethodGet, "/menus", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var results []string
	err = json.Unmarshal(respWr.Body(), &results)
	s.Require().NoError(err)
	s.Require().Len(results, 2)
	s.Require().Equal([]string{menus[0].Name, menus[1].Name}, results)
}

func (s *testSuite) testDeleteMenus() {
	// POST /menus should create two menus
	req, err := http.NewRequest(http.MethodPost, "/menus", test.JsonReaderForObject(menus))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	// DELETE /menus/<name> should delete the menu
	req, err = http.NewRequest(http.MethodDelete, "/menus/"+menus[0].Name, nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	// GET /menus should return the one remaining menu
	req, err = http.NewRequest(http.MethodGet, "/menus", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var results []string
	err = json.Unmarshal(respWr.Body(), &results)
	s.Require().NoError(err)
	s.Require().Len(results, 1)
	s.Require().Equal(menus[1].Name, results[0])
}

func (s *testSuite) testGetMenu() {
	// POST /menus should create two menus
	req, err := http.NewRequest(http.MethodPost, "/menus", test.JsonReaderForObject(menus))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	// GET /menus/<name> should return the menu
	req, err = http.NewRequest(http.MethodGet, "/menus/"+menus[0].Name, nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var menu wire.Menu
	err = json.Unmarshal(respWr.Body(), &menu)
	s.Require().NoError(err)
	s.Require().Equal(menus[0], menu)
}

func (s *testSuite) testAddRemoveRecipes() {
	// POST /menus should create two menus
	req, err := http.NewRequest(http.MethodPost, "/menus", test.JsonReaderForObject(menus))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	// GET /menus/<name>/recipes should return empty list
	req, err = http.NewRequest(http.MethodGet, "/menus/"+menus[0].Name+"/recipes", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var recipes []string
	err = json.Unmarshal(respWr.Body(), &recipes)
	s.Require().NoError(err)
	s.Require().Len(recipes, 0)

	// POST /menus/<name>/recipes/<recipe_id> should add recipe
	recipes = []string{"recipe1", "recipe2", "recipe3"}
	for i := range recipes {
		req, err = http.NewRequest(http.MethodPost, "/menus/"+menus[0].Name+"/recipes/"+recipes[i], nil)
		s.Require().NoError(err)

		respWr = test.NewResponseWriter()
		s.Api.Handle(respWr, req)
		s.Require().Equal(http.StatusOK, respWr.StatusCode())
	}

	// POST to /menus/<noexist>/recipes/<recipe_id> should return 500
	req, err = http.NewRequest(http.MethodPost, "/menus/noexist/recipes/"+recipes[0], nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusInternalServerError, respWr.StatusCode())

	// DELETE /menus/<name>/recipes/<recipe_name> should remove recipe
	req, err = http.NewRequest(http.MethodDelete, "/menus/"+menus[0].Name+"/recipes/"+recipes[1], nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	// /GET /menus/<name>/recipes should return recipe
	req, err = http.NewRequest(http.MethodGet, "/menus/"+menus[0].Name+"/recipes", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var updated []string
	err = json.Unmarshal(respWr.Body(), &updated)
	s.Require().NoError(err)
	s.Require().Len(updated, 2)
	s.Require().Equal([]string{recipes[0], recipes[2]}, updated)
}
