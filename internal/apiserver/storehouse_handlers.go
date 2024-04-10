package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	auto_reference_catalog_service "github.com/maximprokopchuk/auto_reference_catalog_service/pkg/api"
	storehouse_service_api "github.com/maximprokopchuk/storehouse_service/pkg/api"
)

func (s *APIServer) handleGetStorehouseItems() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		city_id_arg := r.URL.Query().Get("city_id")
		if city_id_arg == "" {
			http.Error(w, "Missing 'city_id' param", http.StatusBadRequest)
			return
		}
		city_id, err := strconv.Atoi(city_id_arg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		componentsIdsArg := r.URL.Query().Get("components_ids")
		if componentsIdsArg == "" {
			http.Error(w, "Missing 'components_ids' components_ids", http.StatusBadRequest)
			return
		}
		componentsIdsStr := strings.Split(componentsIdsArg, ",")
		componentsIds := make([]int32, 0)
		for _, idStr := range componentsIdsStr {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			componentsIds = append(componentsIds, int32(id))
		}
		storehouses_params := &storehouse_service_api.GetStorehousesListByCityIdRequest{
			CityId: int32(city_id),
		}
		storehouses, err := s.GrpcClient.StoreHouseClient.GetStorehousesListByCityId(ctx, storehouses_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(storehouses.Result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var storehouse = storehouses.Result[0]

		storehouse_items_params := &storehouse_service_api.GetStorehouseItemsByStorehouseIdAndComponentsIdsRequest{
			StorehouseId: storehouse.Id,
			ComponentIds: componentsIds,
		}
		storehouse_items, err := s.GrpcClient.StoreHouseClient.GetStorehouseItemsByStorehouseIdAndComponentsIds(ctx, storehouse_items_params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(storehouse_items.Result) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		renderJSON(w, storehouse_items)
	}
}

type CreateStorehouseItem struct {
	CityId      int32 `json:"city_id"`
	ComponentId int32 `json:"component_id"`
	Count       int32 `json:"count"`
}

func (s *APIServer) handleCreateStorehouseItem() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateStorehouseItem
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.ComponentId == 0 {
			http.Error(w, "Missing 'component_id' param", http.StatusBadRequest)
			return
		}

		if body.CityId == 0 {
			http.Error(w, "Missing 'city_id' param", http.StatusBadRequest)
			return
		}

		if body.Count <= 0 {
			http.Error(w, "'count' param should be greater than 0", http.StatusBadRequest)
			return
		}
		storehouses_params := &storehouse_service_api.GetStorehousesListByCityIdRequest{
			CityId: body.CityId,
		}
		storehouses, err := s.GrpcClient.StoreHouseClient.GetStorehousesListByCityId(ctx, storehouses_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var storehouse *storehouse_service_api.Storehouse

		if len(storehouses.Result) == 0 {
			storehouseParams := &storehouse_service_api.CreateStorehouseRequest{
				Name:   "default_storehouse",
				CityId: body.CityId,
			}
			storehouse_resp, err := s.GrpcClient.StoreHouseClient.CreateStorehouse(ctx, storehouseParams)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			storehouse = storehouse_resp.Result
		} else {
			storehouse = storehouses.Result[0]
		}

		item_params := &storehouse_service_api.CreateStorehouseItemForStorehoseRequest{
			StorehouseId: storehouse.Id,
			ComponentId:  body.ComponentId,
			Count:        body.Count,
		}
		item, err := s.GrpcClient.StoreHouseClient.CreateStorehouseItemForStorehouse(ctx, item_params)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		renderJSON(w, item)
	}
}

func (s *APIServer) getAllChildComponentsIds(ctx context.Context, componentId int32) ([]int32, error) {
	params := &auto_reference_catalog_service.GetChildComponentsByComponentRequest{
		ParentId: componentId,
	}
	components, err := s.GrpcClient.AutoReferenceCatalogClient.GetChildComponentsByComponent(ctx, params)
	if err != nil {
		return nil, err
	}
	var component_ids = make([]int32, 0)
	for _, component := range components.Result {
		component_ids = append(component_ids, component.Id)
		childs, err := s.getAllChildComponentsIds(ctx, component.Id)
		if err != nil {
			return nil, err
		}
		component_ids = append(component_ids, childs...)
	}
	return component_ids, nil
}

func (s *APIServer) handleDeleteStorehouseItem() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		get_params := &storehouse_service_api.GetStorehouseItemByIdRequest{
			Id: int32(id),
		}
		item, err := s.GrpcClient.StoreHouseClient.GetStorehouseItemById(ctx, get_params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// we need to remove all storehouse items related to child components of current item's compoennt
		component_ids, err := s.getAllChildComponentsIds(ctx, item.Result.ComponentId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		component_ids = append(component_ids, item.Result.ComponentId)
		params := &storehouse_service_api.DeleteStorehouseItemsByComponentIdsRequest{
			StorehouseId: item.Result.StorehouseId,
			ComponentIds: component_ids,
		}
		_, err = s.GrpcClient.StoreHouseClient.DeleteStorehouseItemsByComponentIds(ctx, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type UpdateStorehouseItemBody struct {
	Count int32 `json:"count"`
}

func (s *APIServer) handleUpdateStorehouseItem() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		var body CreateStorehouseItem
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if body.Count <= 0 {
			http.Error(w, "'count' param should be greater than 0", http.StatusBadRequest)
			return
		}
		params := &storehouse_service_api.UpdateStorehouseItemRequest{
			Id:    int32(id),
			Count: body.Count,
		}
		rec, err := s.GrpcClient.StoreHouseClient.UpdateStorehouseItem(ctx, params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderJSON(w, rec)
	}
}
