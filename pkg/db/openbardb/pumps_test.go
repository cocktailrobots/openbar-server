package openbardb

import "context"

func (s *testSuite) TestPumps() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)

	n, err := CountPumpRows(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(0, n)

	err = SetConfig(ctx, tx, map[string]string{NumPumpsConfigKey: "8"})
	s.Require().NoError(err)

	n, err = CountPumpRows(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal(8, n)

	pumps := []Pump{
		{Idx: 0, MlPerSec: 100},
		{Idx: 2, MlPerSec: 102},
		{Idx: 4, MlPerSec: 104},
		{Idx: 6, MlPerSec: 106},
	}
	err = UpdatePumps(ctx, tx, pumps)
	s.Require().NoError(err)

	pumps, err = ListPumps(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal([]Pump{
		{Idx: 0, MlPerSec: 100},
		{Idx: 1, MlPerSec: 0},
		{Idx: 2, MlPerSec: 102},
		{Idx: 3, MlPerSec: 0},
		{Idx: 4, MlPerSec: 104},
		{Idx: 5, MlPerSec: 0},
		{Idx: 6, MlPerSec: 106},
		{Idx: 7, MlPerSec: 0},
	}, pumps)

	err = SetConfig(ctx, tx, map[string]string{NumPumpsConfigKey: "6"})
	s.Require().NoError(err)

	pumps, err = ListPumps(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal([]Pump{
		{Idx: 0, MlPerSec: 100},
		{Idx: 1, MlPerSec: 0},
		{Idx: 2, MlPerSec: 102},
		{Idx: 3, MlPerSec: 0},
		{Idx: 4, MlPerSec: 104},
		{Idx: 5, MlPerSec: 0},
	}, pumps)

	err = SetConfig(ctx, tx, map[string]string{NumPumpsConfigKey: "10"})
	s.Require().NoError(err)

	pumps, err = ListPumps(ctx, tx)
	s.Require().NoError(err)
	s.Require().Equal([]Pump{
		{Idx: 0, MlPerSec: 100},
		{Idx: 1, MlPerSec: 0},
		{Idx: 2, MlPerSec: 102},
		{Idx: 3, MlPerSec: 0},
		{Idx: 4, MlPerSec: 104},
		{Idx: 5, MlPerSec: 0},
		{Idx: 6, MlPerSec: 0},
		{Idx: 7, MlPerSec: 0},
		{Idx: 8, MlPerSec: 0},
		{Idx: 9, MlPerSec: 0},
	}, pumps)
}
