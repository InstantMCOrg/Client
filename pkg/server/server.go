package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const port = 25575

func Handle(r *mux.Router) {
	log.Printf("Starting HTTP Server on port %d...\n", port)
	for true { // Handle forever
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
		if err != nil {
			fmt.Println("An error occured while serving http endpoint:", err)
		}
	}
}
