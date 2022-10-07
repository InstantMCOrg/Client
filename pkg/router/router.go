package router

import "github.com/gorilla/mux"

func Register() *mux.Router {
	r := mux.NewRouter()
	r.Use(authMiddleware)
	r.HandleFunc("/", rootRoute).Methods("GET")
	r.HandleFunc("/server/start", start).Methods("GET")
	r.HandleFunc("/server/stop", stop).Methods("GET")
	r.HandleFunc("/server/player/op/{name}", op).Methods("GET") // give player operator rights
	return r
}
