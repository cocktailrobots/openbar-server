package apis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/cocktailrobots/openbar-server/pkg/util/dbutils"
	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var ErrMethodNotAllowed = errors.New("method not allowed")
var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")
var ErrAlreadyExists = errors.New("already exists")

type API struct {
	logger  *zap.Logger
	txp     dbutils.TxProvider
	handler http.Handler
}

func NewAPI(logger *zap.Logger, txp dbutils.TxProvider, handler http.Handler) *API {
	return &API{
		logger:  logger,
		txp:     txp,
		handler: handler,
	}
}

func (api *API) Logger() *zap.Logger {
	return api.logger
}

func (api *API) Handle(w http.ResponseWriter, r *http.Request) {
	api.handler.ServeHTTP(w, r)
}

func (api *API) Close() error {
	return nil
}

func (api *API) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	api.Respond(w, r, nil, nil)
}

func (api *API) OptionsResponse(allowedMethods []string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
}

func (api *API) Respond(w http.ResponseWriter, r *http.Request, respObj any, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, ErrBadRequest) {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if errors.Is(err, dbr.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if errors.Is(err, ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, ErrMethodNotAllowed) {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			w.WriteHeader(http.StatusConflict)
			return
		}

		api.logger.Info("Error processing "+r.URL.Path, zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if respObj == nil {
			respObj = struct{}{}
		}

		jsonData, err := json.Marshal(respObj)
		if err != nil {
			api.logger.Info("Error marshaling response", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonData)
		if err != nil {
			api.logger.Info("Error writing response", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Error(err))
			return
		}
	}
}

func (api *API) Transaction(ctx context.Context, fn func(tx *dbr.Tx) error) error {
	return api.txp.Transaction(ctx, fn)
}

func GetPathTokens(r *http.Request) []string {
	tokens := strings.Split(r.URL.Path, "/")

	for i := 0; i < len(tokens); i++ {
		if tokens[i] == "" {
			tokens = append(tokens[:i], tokens[i+1:]...)
			i--
		}
	}

	return tokens
}
