package openbarapi

import (
	"context"
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"net/http"
	"strconv"
)

func (s *testSuite) TestConfig() {
	cfg := wire.Config{db.NumPumpsConfigKey: "0"}
	getCfgAndTest(s, cfg)
	testPumpsAndFluids(s, 0)

	// Set to an empty config and test. (Will still contain required keys)
	setConfig(s, wire.Config{})
	getCfgAndTest(s, cfg)
	testPumpsAndFluids(s, 0)

	// Post to a non-empty config and test
	cfg = wire.Config{
		"one":   "1",
		"two":   "2",
		"three": "3",
	}
	setConfig(s, cfg)
	cfg[db.NumPumpsConfigKey] = "0" // Required key
	getCfgAndTest(s, cfg)

	// Test patching a single value
	reqObj := map[string]string{"one": "uno"}
	req, err := http.NewRequest(http.MethodPatch, "/config/one", test.JsonReaderForObject(reqObj))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	expected := mapCopy(cfg)
	expected["one"] = "uno"
	getCfgAndTest(s, expected)

	// Test get a single value
	req, err = http.NewRequest(http.MethodGet, "/config/one", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var cfgResp wire.Config
	err = json.Unmarshal(respWr.Body(), &cfgResp)
	s.Require().NoError(err)
	s.Require().Len(cfgResp, 1)
	s.Require().Equal("uno", cfgResp["one"])

	// Test patching a value that doesn't exist fails
	reqObj = map[string]string{"four": "4"}
	req, err = http.NewRequest(http.MethodPatch, "/config/four", test.JsonReaderForObject(reqObj))
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusNotFound, respWr.StatusCode())

	// Test posting a value that already exists fails
	reqObj = map[string]string{"one": "1"}
	req, err = http.NewRequest(http.MethodPost, "/config/one", test.JsonReaderForObject(reqObj))
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusConflict, respWr.StatusCode())

	// Test posting a value that doesn't exist succeeds
	reqObj = map[string]string{"four": "4"}
	req, err = http.NewRequest(http.MethodPost, "/config/four", test.JsonReaderForObject(reqObj))
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	expected["four"] = "4"
	getCfgAndTest(s, expected)

	// Test deleting a value that exists succeeds
	req, err = http.NewRequest(http.MethodDelete, "/config/four", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	delete(expected, "four")
	getCfgAndTest(s, expected)

	// Test deleting a value that doesn't exist succeeds
	req, err = http.NewRequest(http.MethodDelete, "/config/four", nil)
	s.Require().NoError(err)

	respWr = test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	getCfgAndTest(s, expected)

	updateNumPumpsConfigKey(s, 8)
	updateNumPumpsConfigKey(s, 6)
	updateNumPumpsConfigKey(s, 10)

	expected = wire.Config{db.NumPumpsConfigKey: "0"}
	setConfig(s, wire.Config{})
	getCfgAndTest(s, expected)
	testPumpsAndFluids(s, 0)
}

func updateNumPumpsConfigKey(s *testSuite, numPumps int) {
	cfg := map[string]string{db.NumPumpsConfigKey: strconv.FormatInt(int64(numPumps), 10)}

	req, err := http.NewRequest(http.MethodPost, "/config", test.JsonReaderForObject(cfg))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	getCfgAndTest(s, cfg)
	testPumpsAndFluids(s, numPumps)
}

func setConfig(s *testSuite, cfg wire.Config) {
	req, err := http.NewRequest(http.MethodPost, "/config", test.JsonReaderForObject(cfg))
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())
}

func getCfgAndTest(s *testSuite, expected wire.Config) {
	req, err := http.NewRequest(http.MethodGet, "/config", nil)
	s.Require().NoError(err)

	respWr := test.NewResponseWriter()
	s.Api.Handle(respWr, req)
	s.Require().Equal(http.StatusOK, respWr.StatusCode())

	var cfg wire.Config
	respBody := respWr.Body()
	err = json.Unmarshal(respBody, &cfg)
	s.Require().NoError(err)
	s.Require().Equal(expected, cfg)
}

// mapCopy uses generic key k and value v to create a copy of map m
func mapCopy[K comparable, V any](m map[K]V) map[K]V {
	c := make(map[K]V)
	for k, v := range m {
		c[k] = v
	}
	return c
}

func testPumpsAndFluids(s *testSuite, expected int) {
	ctx := context.Background()

	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)
	defer tx.Rollback()

	var numPumps int
	err = tx.Select("count(*)").From(db.PumpsTable).LoadOneContext(ctx, &numPumps)
	s.Require().NoError(err)
	s.Require().Equal(expected, numPumps)

	var numFluids int
	err = tx.Select("count(*)").From(db.FluidsTable).LoadOneContext(ctx, &numFluids)
	s.Require().NoError(err)
	s.Require().Equal(expected, numFluids)
}
