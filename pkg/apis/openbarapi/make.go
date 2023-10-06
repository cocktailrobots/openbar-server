package openbarapi

import (
	"encoding/json"
	"fmt"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/db/openbardb"
	"github.com/gocraft/dbr/v2"
	"net/http"
	"time"
)

func (api *OpenBarAPI) MakeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Method == http.MethodOptions {
		api.OptionsResponse([]string{http.MethodOptions, http.MethodPost}, w, r)
		return
	} else if r.Method != http.MethodPost {
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
		return
	}

	var req wire.MakeRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	var pumps []openbardb.Pump
	var fluids []openbardb.Fluid
	err = api.Transaction(ctx, func(tx *dbr.Tx) error {
		var err error
		pumps, err = openbardb.ListPumps(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to list pumps: %w", err)
		}

		fluids, err = openbardb.ListFluids(ctx, tx)
		if err != nil {
			return fmt.Errorf("failed to list fluids: %w", err)
		}

		return nil
	})

	if err != nil {
		api.Respond(w, r, nil, err)
		return
	}

	if len(pumps) != len(fluids) {
		api.Respond(w, r, nil, fmt.Errorf("pumps and fluids do not match"))
		return
	}

	if len(pumps) != api.hw.NumPumps() {
		api.Respond(w, r, nil, fmt.Errorf("pumps and hardware do not match"))
		return
	}

	pumpIndices, err := api.getPumpIndices(req, fluids)
	if err != nil {
		api.Respond(w, r, nil, err)
		return
	}

	timesForPumps, err := api.getPumpTimes(pumpIndices, pumps)
	if err != nil {
		api.Respond(w, r, nil, err)
		return
	}

	err = api.hw.RunForTimes(timesForPumps)
	api.Respond(w, r, nil, err)
}

type idxVolTuple struct {
	Idx   int
	VolMl uint
}

func (api *OpenBarAPI) getPumpIndices(req wire.MakeRequest, fluids []openbardb.Fluid) ([]idxVolTuple, error) {
	indicesPerFluid := make([][]int, len(req.FluidVolumes))
	for i, fv := range req.FluidVolumes {
		for _, fluid := range fluids {
			if *fluid.Fluid == fv.Fluid {
				indicesPerFluid[i] = append(indicesPerFluid[i], fluid.Idx)
			}
		}
	}

	// choose pumps. If more than one index choose the one that's been run the least
	indexPerFluid := make([]idxVolTuple, len(req.FluidVolumes))
	for i, fluidIndices := range indicesPerFluid {
		if len(fluidIndices) == 0 {
			return nil, fmt.Errorf("fluid %s not found", req.FluidVolumes[i].Fluid)
		}

		runtime := time.Duration(0x7fffffffffffffff)
		for _, idx := range fluidIndices {
			pumpRuntime := api.hw.TimeRun(idx)
			if pumpRuntime < runtime {
				runtime = pumpRuntime
				indexPerFluid[i] = idxVolTuple{
					Idx:   idx,
					VolMl: req.FluidVolumes[i].VolumeMl,
				}
			}
		}
	}

	return indexPerFluid, nil
}

func (api *OpenBarAPI) getPumpTimes(pumpIndicesAndVols []idxVolTuple, pumps []openbardb.Pump) ([]time.Duration, error) {
	for i, p := range pumps {
		if p.Idx != i {
			return nil, fmt.Errorf("pump indices are not sequential")
		}
	}

	timesForPumps := make([]time.Duration, len(pumps))
	for _, idxVol := range pumpIndicesAndVols {
		pump := pumps[idxVol.Idx]
		seconds := float64(idxVol.VolMl) / pump.MlPerSec
		timesForPumps[idxVol.Idx] = time.Duration(seconds * float64(time.Second))
	}

	return timesForPumps, nil
}
