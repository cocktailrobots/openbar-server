package openbardb

import (
	"context"
	"github.com/cocktailrobots/openbar-server/pkg/util"
)

func (s *testSuite) TestFluids() {
	s.Run("ListFluids", s.testListFluids)
	s.Run("UpdateFluids", s.testUpdateFluids)
}

func (s *testSuite) testListFluids() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)
	err = InitFluids(ctx, tx, 8)
	s.Require().NoError(err)

	fluids, err := ListFluids(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(fluids, 8)

	for i := 0; i < len(fluids); i++ {
		s.Require().Equal(i, fluids[i].Idx)
		s.Require().Nil(fluids[i].Fluid)
	}
}

func (s *testSuite) testUpdateFluids() {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx)
	s.Require().NoError(err)
	err = InitFluids(ctx, tx, 8)
	s.Require().NoError(err)

	fluids := []Fluid{
		{
			Fluid: util.Ptr("zero"),
			Idx:   0,
		},
		{
			Fluid: util.Ptr("two"),
			Idx:   2,
		},
		{
			Fluid: util.Ptr("four"),
			Idx:   4,
		},
		{
			Fluid: util.Ptr("six"),
			Idx:   6,
		},
	}

	err = UpdateFluids(ctx, tx, fluids)
	s.Require().NoError(err)

	listed, err := ListFluids(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(listed, 8)

	for i := range listed {
		if i%2 == 0 {
			s.Require().Equal(i, listed[i].Idx)
			s.Require().NotNil(listed[i].Fluid)
		} else {
			s.Require().Equal(i, listed[i].Idx)
			s.Require().Nil(listed[i].Fluid)
		}
	}

	fluids = []Fluid{
		{
			Fluid: util.Ptr("zero"),
			Idx:   0,
		},
		{
			Fluid: util.Ptr("one"),
			Idx:   1,
		},
		{
			Fluid: util.Ptr("two"),
			Idx:   2,
		},
		{
			Fluid: util.Ptr("three"),
			Idx:   3,
		},
	}

	err = UpdateFluids(ctx, tx, fluids)
	s.Require().NoError(err)

	listed, err = ListFluids(ctx, tx)
	s.Require().NoError(err)
	s.Require().Len(listed, 8)

	for i := range listed {
		if i < 4 {
			s.Require().Equal(i, listed[i].Idx)
			s.Require().NotNil(listed[i].Fluid)
		} else {
			s.Require().Equal(i, listed[i].Idx)
			s.Require().Nil(listed[i].Fluid)
		}
	}
}
