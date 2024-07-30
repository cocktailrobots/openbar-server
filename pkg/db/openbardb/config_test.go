package openbardb

import "context"

func (s *testSuite) TestConfig() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	config, err := GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey:   "0",
		DefaultVolConfigKey: "133",
	}, config)

	err = SetConfig(ctx, tx, map[string]string{
		"test": "test",
	})
	s.Require().NoError(err)

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey:   "0",
		"test":              "test",
		DefaultVolConfigKey: "133",
	}, config)

	err = SetConfig(ctx, tx, map[string]string{
		NumPumpsConfigKey: "1",
	})

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey:   "1",
		DefaultVolConfigKey: "133",
	}, config)

	err = DeleteConfigValues(ctx, tx, NumPumpsConfigKey)
	s.Require().NoError(err)

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey:   "0",
		DefaultVolConfigKey: "133",
	}, config)

	err = SetConfig(ctx, tx, map[string]string{
		NumPumpsConfigKey: "1",
		"foo":             "bar",
		"baz":             "qux",
	})

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey:   "1",
		DefaultVolConfigKey: "133",
		"foo":               "bar",
		"baz":               "qux",
	}, config)

	err = DeleteConfigValues(ctx, tx, NumPumpsConfigKey, "foo")
	s.Require().NoError(err)

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey:   "0",
		DefaultVolConfigKey: "133",
		"baz":               "qux",
	}, config)

}
