package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/grpcclient"
	"github.com/rs/cors"
)

type APIServer struct {
	config     *Config
	router     *mux.Router
	handler    *http.Handler
	GrpcClient *grpcclient.GRPCClient
}

func New(config *Config, grpc_client *grpcclient.GRPCClient) *APIServer {
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.AllowedOrigin},
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)
	return &APIServer{
		config:     config,
		router:     router,
		handler:    &handler,
		GrpcClient: grpc_client,
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()
	return http.ListenAndServe(s.config.BindAddr, *s.handler)
}
