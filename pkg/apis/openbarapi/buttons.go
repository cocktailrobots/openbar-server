package openbarapi

import (
	"context"
	"encoding/json"
	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/hardware"
	"net/http"
	"time"
)

func (api *OpenBarAPI) ButtonsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions}, w, r)
	case http.MethodPost:
		api.setButtonState(ctx, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

func (api *OpenBarAPI) setButtonState(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req wire.ButtonState
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		api.Respond(w, r, nil, apis.ErrBadRequest)
		return
	}

	duration := 250 * time.Millisecond
	if req.DurationMs != 0 {
		duration = time.Duration(req.DurationMs) * time.Millisecond
	}

	runTimes := make([]time.Duration, api.hw.NumPumps())
	for i := 0; i < api.hw.NumPumps(); i++ {
		runTimes[i] = 0
		for _, j := range req.DepressedButtons {
			if i == j {
				runTimes[i] = duration
				break
			}
		}
	}

	direction := hardware.Forward
	if !req.Forward {
		direction = hardware.Backward
	}

	if req.Async {
		err = api.ashwr.RunPumps(direction, runTimes)
	} else {
		err = api.hw.RunForTimes(runTimes)
	}

	api.Respond(w, r, nil, err)
}
