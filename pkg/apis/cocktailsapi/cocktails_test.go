package cocktailsapi

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"net/http"
	"path"
)

func ptr[T any](t T) *T {
	return &t
}

func (s *testSuite) TestCocktails() {
	s.Run("ListCocktails", s.testListCocktails)
	s.Run("GetCocktail", func() {
		s.Run("americano", s.testGetCocktail)
		s.Run("noexist", s.testGetNoExistCocktail)
	})
	s.Run("TestPostCocktails", func() {
		s.Run("Success", s.testPostCocktailsSuccess)
		s.Run("Bad Request", s.testPostCocktailsBadRequest)
		s.Run("Duplicate", s.testDuplicateCocktail)
	})
	s.Run("TestPatchCocktails", func() {
		s.Run("Success", s.testPatchCocktailsSuccess)
		s.Run("noexist", s.testPatchUnknownCocktails)
	})
	s.Run("TestDeleteCocktails", func() {
		s.Run("Success", s.testDeleteCocktailsSuccess)
		s.Run("noexist", s.testDeleteUnknownCocktails)
	})
}

func (s *testSuite) testListCocktails() {
	req, err := http.NewRequest(http.MethodGet, "/cocktails", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var cocktails []wire.Cocktail
	err = json.Unmarshal(respWr.Body(), &cocktails)
	s.Require().NoError(err)
	s.Require().Len(cocktails, 6)
}

func (s *testSuite) testGetNoExistCocktail() {
	req, err := http.NewRequest(http.MethodGet, "/cocktails/noexist", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())
}

func (s *testSuite) testGetCocktail() {
	req, err := http.NewRequest(http.MethodGet, "/cocktails/americano", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var cocktail wire.Cocktail
	respBody := respWr.Body()
	err = json.Unmarshal(respBody, &cocktail)
	s.Require().NoError(err)
	s.Require().Equal("americano", cocktail.Name)
	s.Require().Equal("Americano", cocktail.DisplayName)
	s.Require().NotEmpty(cocktail.Description)
}

func (s *testSuite) testPostCocktailsSuccess() {
	cocktail := []wire.Cocktail{
		{
			Name:        "test",
			DisplayName: "Test Display Name",
			Description: ptr("Test Description"),
		},
	}
	req, err := http.NewRequest(http.MethodPost, "/cocktails", test.JsonReaderForObject(cocktail))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())
}

func (s *testSuite) testPostCocktailsBadRequest() {
	cocktail := []wire.Cocktail{
		{
			DisplayName: "Test Display Name",
			Description: ptr("Test Description"),
		},
	}
	req, err := http.NewRequest(http.MethodPost, "/cocktails", test.JsonReaderForObject(cocktail))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusBadRequest, respWr.StatusCode())
}

func (s *testSuite) testDuplicateCocktail() {
	cocktail := []wire.Cocktail{
		{
			Name:        "americano", // already exists
			DisplayName: "Americano",
		},
	}

	req, err := http.NewRequest(http.MethodPost, "/cocktails", test.JsonReaderForObject(cocktail))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusConflict, respWr.StatusCode())
}

func (s *testSuite) testPatchCocktailsSuccess() {
	const (
		name               = "americano"
		updatedDisplayName = "Americano Updated"
		updatedDescription = "Americano Description Updated"
	)

	p := path.Join("/cocktails", name)
	cocktail := wire.Cocktail{
		Name:        name,
		DisplayName: updatedDisplayName,
		Description: ptr(updatedDescription),
	}

	req, err := http.NewRequest(http.MethodPatch, p, test.JsonReaderForObject(cocktail))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodGet, p, nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var updated wire.Cocktail
	err = json.Unmarshal(respWr.Body(), &updated)
	s.Require().NoError(err)
	s.Require().Equal(updatedDisplayName, updated.DisplayName)
	s.Require().Equal(updatedDescription, *updated.Description)
}

func (s *testSuite) testPatchUnknownCocktails() {
	const (
		name = "noexist"
	)

	p := path.Join("/cocktails", name)
	cocktail := wire.Cocktail{
		Name:        name,
		DisplayName: name,
		Description: ptr(name),
	}

	req, err := http.NewRequest(http.MethodPatch, p, test.JsonReaderForObject(cocktail))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())
}

func (s *testSuite) testDeleteCocktailsSuccess() {
	const (
		name = "americano"
	)

	p := path.Join("/cocktails", name)
	req, err := http.NewRequest(http.MethodDelete, p, nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	req, err = http.NewRequest(http.MethodGet, p, nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())
}

func (s *testSuite) testDeleteUnknownCocktails() {
	const (
		name = "noexist"
	)

	p := path.Join("/cocktails", name)
	req, err := http.NewRequest(http.MethodDelete, p, nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())
}
