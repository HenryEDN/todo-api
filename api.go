package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct{
	listenAddr string
	database DatabaseReader
}

type Handler struct{
	database DatabaseReader
}

func NewHandler(database DatabaseReader) *Handler{
	return &Handler{database: database}
}


func NewAPIServer(listenAddr string, database DatabaseReader) *APIServer{
	return &APIServer{
		listenAddr: listenAddr,
		database: database,
	}
}

func (server *APIServer) Run(){
	router := mux.NewRouter()

	handler := NewHandler(server.database)

	router.HandleFunc("/", makeHTTPHandleFunc(handler.handleRoot))
	router.HandleFunc("/users", makeHTTPHandleFunc(handler.handleUsers))

	log.Printf("JSON API server running on: http://localhost%v", server.listenAddr)

	http.ListenAndServe(server.listenAddr, router)
}

func (h Handler) handleRoot(w http.ResponseWriter, r *http.Request) error{
	type SayHello struct{	Text string	`json:"text"`	}

	return writeJSON(w, http.StatusOK, SayHello{Text: "Hello from root"})
}

func newUser(user *CreateUserDTO)(*User, error){
	encryptPW, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil{
		return nil, err
	}
	return &User{
		UserID: uint64(rand.Intn(1000000)),
		Email: user.Email,
		Username: user.Username,
		Password: string(encryptPW),
		Creation_time: time.Now(),
	}, nil
}

func (h Handler) handleUsers(w http.ResponseWriter, r *http.Request)(error){
	switch(r.Method){
	case "GET":
		return h.getUsers(w,r)
	case "POST":
		return h.createUser(w,r)
	default:
		return writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: fmt.Sprintf("method is not allowed, %v", r.Method)})
	}
}

func (h Handler) createUser(w http.ResponseWriter, r *http.Request)(error){
	if r.Method != "POST"{
		return writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: fmt.Sprintf("method is not allowed, %v", r.Method)})
	}

	createUserData := new(CreateUserDTO)

	err := json.NewDecoder(r.Body).Decode(&createUserData)
	if err != nil{
		return writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	newUser, err := newUser(createUserData)
	if err != nil{
		return writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	resp, err := h.database.CreateUser(*newUser)
	if err != nil{
		return writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, resp)
}

func (h Handler) getUsers(w http.ResponseWriter, r *http.Request)(error){
	if r.Method != "GET"{
		return writeJSON(w, http.StatusMethodNotAllowed, APIError{Error: fmt.Sprintf("method is not allowed, %v", r.Method)})
	}

	users, err := h.database.GetUsers()
	if err != nil{
		return writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return writeJSON(w, http.StatusOK, users)
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