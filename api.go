package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct{
	listenAddr string
}


func NewAPIServer(listenAddr string) *APIServer{
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (server *APIServer) Run(){
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(handleRoot))

	log.Printf("JSON API server running on: http://localhost%v", server.listenAddr)

	http.ListenAndServe(server.listenAddr, router)
}

func handleRoot(w http.ResponseWriter, r *http.Request) error{
	type SayHello struct{	Text string	`json:"text"`	}

	return writeJSON(w, http.StatusOK, SayHello{Text: "Hello from root"})
}

type APIError struct{
	Error string `json:"error"`
}

func makeHTTPHandleFunc(handler apiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err := handler(w, r); err != nil{
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		} 
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) error{
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}