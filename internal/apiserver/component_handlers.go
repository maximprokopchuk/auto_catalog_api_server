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

func (s *APIServer) handleGetCarModel() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := &api.GetCarModelByIdRequest{
			Id: int32(id),
		}
		car_model, err := s.GrpcClient.AutoReferenceCatalogClient.GetCarModelById(ctx, params)
		if err != nil {
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		renderJSON(w, car_model)
	}
}

func (s *APIServer) handleListCarModels() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		car_models, err := s.GrpcClient.AutoReferenceCatalogClient.ListCarModels(ctx, nil)
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
		params := &api.CreateCarModelRequest{
			Name: body.Name,
		}
		response, err := s.GrpcClient.AutoReferenceCatalogClient.CreateCarModel(ctx, params)
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
		_, err = s.GrpcClient.AutoReferenceCatalogClient.DeleteCarModel(ctx, params)
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

func (s *APIServer) handleGetComponents() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		car_model_id_str := r.URL.Query().Get("car_model_id")
		parent_component_id_str := r.URL.Query().Get("parent_component_id")
		if parent_component_id_str == "" && car_model_id_str == "" {
			http.Error(w, "Missing 'parent_component_id' or 'car_model_id' param", http.StatusBadRequest)
			return
		}
		var components *api.ListComponentResponse
		if car_model_id_str != "" {
			car_model_id, err := strconv.Atoi(car_model_id_str)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			components_params := &api.GetTopLevelComponentsByCarModelRequest{
				CarModelId: int32(car_model_id),
			}
			components, err = s.GrpcClient.AutoReferenceCatalogClient.GetTopLevelComponentsByCarModel(ctx, components_params)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			parent_component_id, err := strconv.Atoi(parent_component_id_str)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			components_params := &api.GetChildComponentsByComponentRequest{
				ParentId: int32(parent_component_id),
			}
			components, err = s.GrpcClient.AutoReferenceCatalogClient.GetChildComponentsByComponent(ctx, components_params)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		renderJSON(w, components)

	}
}

type CreateComponentBody struct {
	Name              string `json:"name"`
	CarModelId        int32  `json:"car_model_id"`
	ParentComponentId int32  `json:"parent_component_id"`
}

func (s *APIServer) handleCreateComponent() http.HandlerFunc {
	ctx := context.Background()
	var body CreateComponentBody
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
		if body.CarModelId == 0 && body.ParentComponentId == 0 {
			http.Error(w, "Missing 'car_model_id' or 'parent_component_id' param", http.StatusBadRequest)
			return
		}
		params := &api.CreateComponentRequest{
			Name: body.Name,
		}
		if body.ParentComponentId == 0 {
			params.CarModelId = body.CarModelId
		} else {
			params.ParentId = body.ParentComponentId
		}
		response, err := s.GrpcClient.AutoReferenceCatalogClient.CreateComponent(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		renderJSON(w, response)
	}
}

type UpdateComponentBody struct {
	Name string `json:"name"`
}

func (s *APIServer) handleUpdateComponent() http.HandlerFunc {
	ctx := context.Background()
	var body CreateComponentBody
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.Name == "" {
			http.Error(w, "Missing 'name' param", http.StatusBadRequest)
			return
		}
		params := &api.UpdateComponentRequest{
			Id:   int32(id),
			Name: body.Name,
		}
		response, err := s.GrpcClient.AutoReferenceCatalogClient.UpdateComponent(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		renderJSON(w, response)
	}
}
