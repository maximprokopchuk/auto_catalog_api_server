package apiserver

func (s *APIServer) configureRouter() {
	s.router.Handle("/countries/", s.handleListCountries()).Methods("GET")
	s.router.Handle("/countries/", s.handleCreateCountry()).Methods("POST")
	s.router.Handle("/countries/{id}/", s.handleGetCountry()).Methods("GET")
	s.router.Handle("/countries/{id}/", s.handlerDeleteCountry()).Methods("DELETE")

	s.router.Handle("/cities/", s.handleListCities()).Methods("GET")
	s.router.Handle("/cities/", s.handleCreateCity()).Methods("POST")
	s.router.Handle("/cities/{id}/", s.handleGetCity()).Methods("GET")
	s.router.Handle("/cities/{id}/", s.handleDeleteCity()).Methods("DELETE")

	s.router.Handle("/car_models/", s.handleListCarModels()).Methods("GET")
	s.router.Handle("/car_models/", s.handleCreateCarModel()).Methods("POST")
	s.router.Handle("/car_models/{id}", s.handleDeleteCarModel()).Methods("DELETE")
}
