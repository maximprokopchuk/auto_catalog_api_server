package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	address_service_api "github.com/maximprokopchuk/address_service/pkg/api"
	storehouse_service_api "github.com/maximprokopchuk/storehouse_service/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APIServer) handleGetCountry() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := &address_service_api.GetAddressByIdRequest{
			Id: int32(id),
		}
		city, err := s.GrpcClient.AddressClient.GetAddressById(ctx, params)
		if err != nil {
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if city.Result.Type != "country" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		renderJSON(w, city)
	}
}

func (s *APIServer) handleListCountries() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		params := &address_service_api.ListAddressesByParentIdAndTypeRequest{
			Type: "country",
		}
		countries, err := s.GrpcClient.AddressClient.ListAddressesByParentAndType(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := constructListResponse(countries.Result)
		renderJSON(w, response)
	}
}

type CreateCountryRequestBody struct {
	Name string `json:"name"`
}

func (s *APIServer) handleCreateCountry() http.HandlerFunc {
	ctx := context.Background()
	var body CreateCountryRequestBody
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
		params := &address_service_api.CreateAddressRequest{
			Type: "country",
			Name: body.Name,
		}
		country, err := s.GrpcClient.AddressClient.CreateAddress(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		renderJSON(w, country)
	}
}

func (s *APIServer) handlerDeleteCountry() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var cityIds []int32
		get_cities_params := &address_service_api.ListAddressesByParentIdAndTypeRequest{
			ParentId: int32(id),
			Type:     "city",
		}
		cities, err := s.GrpcClient.AddressClient.ListAddressesByParentAndType(ctx, get_cities_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, city := range cities.Result {
			cityIds = append(cityIds, city.Id)
		}
		delete_storehouses_params := &storehouse_service_api.DeleteStorehousesByCityIdRequest{
			CityIds: cityIds,
		}
		_, err = s.GrpcClient.StoreHouseClient.DeleteStorehousesByCityId(ctx, delete_storehouses_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		delete_country_params := &address_service_api.DeleteAddressRequest{
			Id:   int32(id),
			Type: "country",
		}
		_, err = s.GrpcClient.AddressClient.DeleteAddress(ctx, delete_country_params)
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

func (s *APIServer) handleGetCity() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := &address_service_api.GetAddressByIdRequest{
			Id: int32(id),
		}
		city, err := s.GrpcClient.AddressClient.GetAddressById(ctx, params)
		if err != nil {
			if e, ok := status.FromError(err); ok && e.Code() == codes.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if city.Result.Type != "city" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		renderJSON(w, city)
	}
}

func (s *APIServer) handleListCities() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		country_id_str := r.URL.Query().Get("country_id")
		if country_id_str == "" {
			http.Error(w, "Missing 'country_id' param", http.StatusBadRequest)
			return
		}
		country_id, err := strconv.Atoi(country_id_str)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params := &address_service_api.ListAddressesByParentIdAndTypeRequest{
			Type:     "city",
			ParentId: int32(country_id),
		}
		countries, err := s.GrpcClient.AddressClient.ListAddressesByParentAndType(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		response := constructListResponse(countries.Result)
		renderJSON(w, response)
	}
}

type CreateCityRequestBody struct {
	CountryId int32  `json:"country_id"`
	Name      string `json:"name"`
}

func (s *APIServer) handleCreateCity() http.HandlerFunc {
	ctx := context.Background()
	var body CreateCityRequestBody
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
		if body.CountryId == 0 {
			http.Error(w, "Missing 'country_id' param", http.StatusBadRequest)
			return
		}
		params := &address_service_api.CreateAddressRequest{
			Type:     "city",
			Name:     body.Name,
			ParentId: body.CountryId,
		}
		city, err := s.GrpcClient.AddressClient.CreateAddress(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		renderJSON(w, city)
	}
}

func (s *APIServer) handleDeleteCity() http.HandlerFunc {
	ctx := context.Background()
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		delete_storehouses_params := &storehouse_service_api.DeleteStorehousesByCityIdRequest{
			CityIds: []int32{int32(id)},
		}
		_, err = s.GrpcClient.StoreHouseClient.DeleteStorehousesByCityId(ctx, delete_storehouses_params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		delete_city_params := &address_service_api.DeleteAddressRequest{
			Id:   int32(id),
			Type: "city",
		}

		_, err = s.GrpcClient.AddressClient.DeleteAddress(ctx, delete_city_params)
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
