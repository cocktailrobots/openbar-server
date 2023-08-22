package openbardb

import "context"

func (s *testSuite) TestConfig() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	config, err := GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey: "0",
	}, config)

	err = SetConfig(ctx, tx, map[string]string{
		"test": "test",
	})
	s.Require().NoError(err)

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey: "0",
		"test":            "test",
	}, config)

	err = SetConfig(ctx, tx, map[string]string{
		NumPumpsConfigKey: "1",
	})

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey: "1",
	}, config)

	err = DeleteConfigValues(ctx, tx, NumPumpsConfigKey)
	s.Require().NoError(err)

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey: "0",
	}, config)

	err = SetConfig(ctx, tx, map[string]string{
		NumPumpsConfigKey: "1",
		"foo":             "bar",
		"baz":             "qux",
	})

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey: "1",
		"foo":             "bar",
		"baz":             "qux",
	}, config)

	err = DeleteConfigValues(ctx, tx, NumPumpsConfigKey, "foo")
	s.Require().NoError(err)

	config, err = GetConfig(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(map[string]string{
		NumPumpsConfigKey: "0",
		"baz":             "qux",
	}, config)

}
