package tracking

import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/net/context"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/marcusolsson/goddd/cargo"
)

// MakeHandler ...
func MakeHandler(ctx context.Context, ts Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	trackCargoHandler := kithttp.NewServer(
		ctx,
		makeTrackCargoEndpoint(ts),
		decodeTrackCargoRequest,
		encodeResponse,
		opts...,
	)

	r.Handle("/tracking/v1/cargos/{id}", trackCargoHandler).Methods("GET")
	r.Handle("/tracking/v1/docs", http.StripPrefix("/tracking/v1/docs", http.FileServer(http.Dir("tracking/docs"))))

	return r
}

func decodeTrackCargoRequest(r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errors.New("bad route")
	}
	return trackCargoRequest{ID: id}, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(w, e.error())
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(w http.ResponseWriter, err error) {
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case cargo.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
