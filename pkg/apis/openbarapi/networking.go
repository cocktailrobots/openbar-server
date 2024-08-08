package openbarapi

import (
	"context"
	"go.uber.org/zap"
	"log"
	"net/http"

	"github.com/cocktailrobots/openbar-server/pkg/apis"
	"github.com/cocktailrobots/openbar-server/pkg/apis/wire"
	"github.com/cocktailrobots/openbar-server/pkg/networking"
)

func (api *OpenBarAPI) NetworkingHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodOptions:
		api.OptionsResponse([]string{http.MethodOptions, http.MethodGet}, w, r)
	case http.MethodGet:
		api.getNetworking(ctx, w, r)
	default:
		api.Respond(w, r, nil, apis.ErrMethodNotAllowed)
	}
}

func (api *OpenBarAPI) getNetworking(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	const wifiInterface = "wlan0"
	const wiredInterface = "eth0"

	// this is all best effort. Log, but ignore errors.
	ssid, err := networking.SSIDForHotspot()
	if err != nil {
		api.Logger().Info("Failed to get wifi SSID", zap.String("interface", wifiInterface), zap.Error(err))
	}

	wifiAddr, err := networking.AddrForInterface(wifiInterface)
	if err != nil {
		log.Println("Failed to get wifi address", zap.String("interface", wifiInterface), zap.Error(err))
	}

	wiredAddr, err := networking.AddrForInterface(wiredInterface)
	if err != nil {
		log.Println("Failed to get wired address", zap.String("interface", wiredInterface), zap.Error(err))
	}

	api.Respond(w, r, wire.NetworkingResponse{
		WifiSSID:  ssid,
		WifiAddr:  wifiAddr,
		WiredAddr: wiredAddr,
	}, nil)
}
