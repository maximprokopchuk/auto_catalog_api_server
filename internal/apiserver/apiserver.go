package apiserver

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/maximprokopchuk/auto_catalog_api_server/internal/grpcclient"
)

type APIServer struct {
	config     *Config
	router     *mux.Router
	GrpcClient *grpcclient.GRPCClient
}

func New(config *Config, grpc_client *grpcclient.GRPCClient) *APIServer {
	return &APIServer{
		config:     config,
		router:     mux.NewRouter(),
		GrpcClient: grpc_client,
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()
	return http.ListenAndServe(s.config.BindAddr, s.router)
}
