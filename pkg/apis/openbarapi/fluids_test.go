package openbarapi

import (
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"net/http"
)

func (s *testSuite) TestFluids() {
	req, err := http.NewRequest(http.MethodGet, "/fluids", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var fluids wire.Fluids
	err = json.Unmarshal(respWr.Body(), &fluids)
	s.Require().NoError(err)
	s.Require().Len(fluids, 0)

	/*req, err = http.NewRequest(http.MethodPost, "/fluids")
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusMethodNotAllowed, respWr.StatusCode())*/

}
