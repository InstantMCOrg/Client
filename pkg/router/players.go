package router

import (
	"github.com/gorilla/mux"
	"github.com/instantminecraft/client/pkg/mcserver"
	"net/http"
)

func op(w http.ResponseWriter, r *http.Request) {
	targetName := mux.Vars(r)["name"]

	mcserver.OpPlayer(targetName)

	w.WriteHeader(http.StatusOK)
}
