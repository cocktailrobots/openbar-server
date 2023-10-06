package openbarapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db/openbardb"
	"github.com/cocktailrobots/openbar-server/pkg/hardware"
	"github.com/cocktailrobots/openbar-server/pkg/util"
	"github.com/cocktailrobots/openbar-server/pkg/util/test"
	"github.com/gocraft/dbr/v2"
	"net/http"
	"time"
)

func (s *testSuite) TestGetPumpIndices() {
	tests := []struct {
		name         string
		reqFluidVols []wire.FluidVolume
		fluids       []openbardb.Fluid
		idxVolTuples []idxVolTuple
		expectErr    bool
		runTimes     map[int]time.Duration
	}{
		{
			name: "no duplicates",
			reqFluidVols: []wire.FluidVolume{
				{Fluid: "gin", VolumeMl: 50},
				{Fluid: "campari", VolumeMl: 30},
				{Fluid: "sweet_vermouth", VolumeMl: 40},
			},
			fluids: []openbardb.Fluid{
				{Idx: 0, Fluid: util.Ptr("gin")},
				{Idx: 1, Fluid: util.Ptr("vodka")},
				{Idx: 2, Fluid: util.Ptr("tequila")},
				{Idx: 3, Fluid: util.Ptr("campari")},
				{Idx: 4, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 5, Fluid: util.Ptr("dry_vermouth")},
				{Idx: 6, Fluid: util.Ptr("triple_sec")},
				{Idx: 7, Fluid: util.Ptr("lime_juice")},
			},
			idxVolTuples: []idxVolTuple{
				{Idx: 0, VolMl: 50},
				{Idx: 3, VolMl: 30},
				{Idx: 4, VolMl: 40},
			},
		},
		{
			name: "duplicates with no runtimes",
			reqFluidVols: []wire.FluidVolume{
				{Fluid: "gin", VolumeMl: 50},
				{Fluid: "campari", VolumeMl: 30},
				{Fluid: "sweet_vermouth", VolumeMl: 40},
			},
			fluids: []openbardb.Fluid{
				{Idx: 0, Fluid: util.Ptr("gin")},
				{Idx: 1, Fluid: util.Ptr("campari")},
				{Idx: 2, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 3, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 4, Fluid: util.Ptr("campari")},
				{Idx: 5, Fluid: util.Ptr("gin")},
				{Idx: 6, Fluid: util.Ptr("gin")},
				{Idx: 7, Fluid: util.Ptr("campari")},
			},
			idxVolTuples: []idxVolTuple{
				{Idx: 0, VolMl: 50},
				{Idx: 1, VolMl: 30},
				{Idx: 2, VolMl: 40},
			},
		},
		{
			name: "duplicates with runtimes in reverse pump order",
			reqFluidVols: []wire.FluidVolume{
				{Fluid: "gin", VolumeMl: 50},
				{Fluid: "campari", VolumeMl: 30},
				{Fluid: "sweet_vermouth", VolumeMl: 40},
			},
			fluids: []openbardb.Fluid{
				{Idx: 0, Fluid: util.Ptr("gin")},
				{Idx: 1, Fluid: util.Ptr("campari")},
				{Idx: 2, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 3, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 4, Fluid: util.Ptr("campari")},
				{Idx: 5, Fluid: util.Ptr("gin")},
				{Idx: 6, Fluid: util.Ptr("gin")},
				{Idx: 7, Fluid: util.Ptr("campari")},
			},
			idxVolTuples: []idxVolTuple{
				{Idx: 6, VolMl: 50},
				{Idx: 7, VolMl: 30},
				{Idx: 3, VolMl: 40},
			},
			runTimes: func() map[int]time.Duration {
				m := make(map[int]time.Duration)
				for i := 0; i < 8; i++ {
					m[i] = time.Duration(7-i) * time.Second
				}
				return m
			}(),
		},
		{
			name: "fluid not found",
			reqFluidVols: []wire.FluidVolume{
				{Fluid: "gin", VolumeMl: 50},
				{Fluid: "campari", VolumeMl: 30},
				{Fluid: "sweet_vermouth", VolumeMl: 40},
			},
			fluids: []openbardb.Fluid{
				{Idx: 0, Fluid: util.Ptr("gin")},
				{Idx: 1, Fluid: util.Ptr("campari")},
				{Idx: 2, Fluid: util.Ptr("gin")},
				{Idx: 3, Fluid: util.Ptr("campari")},
				{Idx: 4, Fluid: util.Ptr("gin")},
				{Idx: 5, Fluid: util.Ptr("campari")},
				{Idx: 6, Fluid: util.Ptr("gin")},
				{Idx: 7, Fluid: util.Ptr("campari")},
			},
			expectErr: true,
		},
		{
			name: "unordered fluids",
			reqFluidVols: []wire.FluidVolume{
				{Fluid: "gin", VolumeMl: 50},
				{Fluid: "campari", VolumeMl: 30},
				{Fluid: "sweet_vermouth", VolumeMl: 40},
			},
			fluids: []openbardb.Fluid{
				{Idx: 1, Fluid: util.Ptr("campari")},
				{Idx: 7, Fluid: util.Ptr("campari")},
				{Idx: 2, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 3, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 5, Fluid: util.Ptr("gin")},
				{Idx: 0, Fluid: util.Ptr("gin")},
				{Idx: 4, Fluid: util.Ptr("campari")},
				{Idx: 6, Fluid: util.Ptr("gin")},
			},
			idxVolTuples: []idxVolTuple{
				{Idx: 6, VolMl: 50},
				{Idx: 7, VolMl: 30},
				{Idx: 3, VolMl: 40},
			},
			runTimes: func() map[int]time.Duration {
				m := make(map[int]time.Duration)
				for i := 0; i < 8; i++ {
					m[i] = time.Duration(7-i) * time.Second
				}
				return m
			}(),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := wire.MakeRequest{FluidVolumes: tt.reqFluidVols}

			if tt.runTimes != nil {
				testHw := s.Api.hw.(*hardware.TestHardware)
				for idx, dur := range tt.runTimes {
					testHw.SetRuntime(idx, dur)
				}
			}

			idxVolTuples, err := s.Api.getPumpIndices(req, tt.fluids)

			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.idxVolTuples, idxVolTuples)
			}
		})
	}
}

func pumpsOfSpeed(speed int, numPumps int) []openbardb.Pump {
	pumps := make([]openbardb.Pump, numPumps)
	for i := range pumps {
		pumps[i] = openbardb.Pump{Idx: i, MlPerSec: float64(speed)}
	}
	return pumps
}

func (s *testSuite) TestGetPumpTimes() {
	tests := []struct {
		name      string
		idxVols   []idxVolTuple
		pumps     []openbardb.Pump
		expected  []time.Duration
		expectErr bool
	}{
		{
			name:    "no fluids no times",
			idxVols: []idxVolTuple{},
			pumps:   pumpsOfSpeed(1, 8),
			expected: []time.Duration{
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
			},
		},
		{
			name: "negroni slow pumps",
			idxVols: []idxVolTuple{
				{Idx: 0, VolMl: 50},
				{Idx: 3, VolMl: 30},
				{Idx: 4, VolMl: 40},
			},
			pumps: pumpsOfSpeed(1, 8),
			expected: []time.Duration{
				50 * time.Second,
				0,
				0,
				30 * time.Second,
				40 * time.Second,
				0,
				0,
				0,
			},
		},
		{
			name: "negroni fast pumps",
			idxVols: []idxVolTuple{
				{Idx: 0, VolMl: 50},
				{Idx: 3, VolMl: 30},
				{Idx: 4, VolMl: 40},
			},
			pumps: pumpsOfSpeed(100, 8),
			expected: []time.Duration{
				500 * time.Millisecond,
				0,
				0,
				300 * time.Millisecond,
				400 * time.Millisecond,
				0,
				0,
				0,
			},
		},
		{
			name: "negroni out of order pumps",
			pumps: []openbardb.Pump{
				{Idx: 1, MlPerSec: 100},
				{Idx: 0, MlPerSec: 100},
				{Idx: 3, MlPerSec: 100},
				{Idx: 2, MlPerSec: 100},
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			times, err := s.Api.getPumpTimes(tt.idxVols, tt.pumps)

			if tt.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tt.expected, times)
			}
		})
	}
}

func (s *testSuite) TestMakeHandler() {
	tests := []struct {
		name      string
		fluidVols []wire.FluidVolume
		expected  []time.Duration
		dbFluids  []openbardb.Fluid
		pumps     []openbardb.Pump
	}{
		{
			name: "negroni",
			fluidVols: []wire.FluidVolume{
				{Fluid: "gin", VolumeMl: 50},
				{Fluid: "campari", VolumeMl: 30},
				{Fluid: "sweet_vermouth", VolumeMl: 40},
			},
			dbFluids: []openbardb.Fluid{
				{Idx: 0, Fluid: util.Ptr("gin")},
				{Idx: 1, Fluid: util.Ptr("vodka")},
				{Idx: 2, Fluid: util.Ptr("tequila")},
				{Idx: 3, Fluid: util.Ptr("campari")},
				{Idx: 4, Fluid: util.Ptr("sweet_vermouth")},
				{Idx: 5, Fluid: util.Ptr("dry_vermouth")},
				{Idx: 6, Fluid: util.Ptr("triple_sec")},
				{Idx: 7, Fluid: util.Ptr("lime_juice")},
			},
			pumps: pumpsOfSpeed(100, 8),
			expected: []time.Duration{
				500 * time.Millisecond,
				0,
				0,
				300 * time.Millisecond,
				400 * time.Millisecond,
				0,
				0,
				0,
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ctx := context.Background()
			s.setupPumpsAndFluids(ctx, tt.dbFluids, tt.pumps)

			reqJson, err := json.Marshal(wire.MakeRequest{FluidVolumes: tt.fluidVols})
			s.Require().NoError(err)

			req, err := http.NewRequest(http.MethodPost, "/make", bytes.NewBuffer(reqJson))
			s.Require().NoError(err)

			respWr := test.NewResponseWriter()
			s.Api.Handle(respWr, req)

			s.Require().Equal(http.StatusOK, respWr.StatusCode())

			thw, ok := s.Api.hw.(*hardware.TestHardware)
			s.Require().True(ok)

			for i := 0; i < thw.NumPumps(); i++ {
				s.isClose(tt.expected[i], thw.TimeRun(i))
			}
		})
	}
}

func (s *testSuite) setupPumpsAndFluids(ctx context.Context, fluids []openbardb.Fluid, pumps []openbardb.Pump) {
	err := s.Transaction(ctx, func(tx *dbr.Tx) error {
		err := openbardb.UpdateFluids(ctx, tx, fluids)
		s.Require().NoError(err)

		err = openbardb.UpdatePumps(ctx, tx, pumps)
		s.Require().NoError(err)

		return tx.Commit()
	})
	s.Require().NoError(err)
}

func (s *testSuite) isClose(a, b time.Duration) {
	diff := a - b
	s.Require().True((diff >= 0 && diff < 10*time.Millisecond) || (diff < 0 && diff > -10*time.Millisecond), "expected %s to be close to %s, but is %s different", a.String(), b.String(), diff.String())
}
