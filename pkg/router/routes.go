package router

import (
	"encoding/json"
	"github.com/instantminecraft/client/pkg/mcserver"
	"net/http"
)

func rootRoute(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(map[string]interface{}{
		"server": map[string]interface{}{
			"running": mcserver.IsRunning(),
		},
	})
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
