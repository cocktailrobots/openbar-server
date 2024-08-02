package openbarapi

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os/exec"
	"time"

	"github.com/cocktailrobots/openbar-server/pkg/apis"
)

func (api *OpenBarAPI) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodPost}, w, r)
	case http.MethodPost:
		api.shutdown(ctx, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

func (api *OpenBarAPI) shutdown(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	api.Logger().Info("Shutting down")

	go func() {
		time.Sleep(10 * time.Second)
		cmd := exec.Command("poweroff")
		if err := cmd.Run(); err != nil {
			api.Logger().Info("failed to shutdown", zap.Error(err))
		}
	}()

	api.Respond(w, r, nil, nil)
}
