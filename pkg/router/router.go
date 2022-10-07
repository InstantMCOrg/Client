package router

import "github.com/gorilla/mux"

func Register() *mux.Router {
	r := mux.NewRouter()
	r.Use(authMiddleware)
	r.HandleFunc("/", rootRoute).Methods("GET")
	r.HandleFunc("/server/start", start).Methods("GET")
	return r
}
