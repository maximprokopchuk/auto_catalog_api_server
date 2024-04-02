package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/maximprokopchuk/auto_reference_catalog_service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APIServer) handleListCarModels() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		car_models, err := s.GrpcClient.AutoCatalogClient.ListCarModels(ctx, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := constructListResponse(car_models.Result)
		renderJSON(w, response)
	}
}

type CreateCarModelBody struct {
	Name string `json:"name"`
}

func (s *APIServer) handleCreateCarModel() http.HandlerFunc {
	ctx := context.Background()
	var body CreateCarModelBody
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "Missing 'name' param", http.StatusBadRequest)
			return
		}
		params := api.CreateCarModelRequest{
			Name: body.Name,
		}
		response, err := s.GrpcClient.AutoCatalogClient.CreateCarModel(ctx, &params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		renderJSON(w, response)
	}
}

func (s *APIServer) handleDeleteCarModel() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := &api.DeleteCarModelRequest{
			Id: int32(id),
		}
		_, err = s.GrpcClient.AutoCatalogClient.DeleteCarModel(ctx, params)
		if err != nil {
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
