package apiserver

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Result interface{} `json:"result"`
}

func constructListResponse[T any](result []T) *Response {
	response := &Response{}
	if len(result) == 0 {
		response.Result = make([]string, 0)
	} else {
		response.Result = result
	}
	return response
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	json_response, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_response)
}
