package apiserver

import (
	"context"
	"net/http"
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
