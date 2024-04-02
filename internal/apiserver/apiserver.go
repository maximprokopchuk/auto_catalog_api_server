package apiserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	config *Config
	router *mux.Router
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Handle("/hello", s.handleHello())
}

type Address struct {
	Id       int32
	Name     string
	Type     string
	ParentId int32
}

func (s *APIServer) handleHello() http.HandlerFunc {
	address := &Address{
		Id:       1,
		Name:     "Moscow",
		Type:     "City",
		ParentId: 1,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := json.Marshal(address)
		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.WriteString(w, string(res)); err != nil {
			log.Fatal(err)
		}
	}
}
