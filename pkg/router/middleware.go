package router

import (
	"github.com/instantminecraft/client/pkg/auth"
	"net/http"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if auth.HasAuthKey() {
			// Check the Auth header
			if r.Header.Get("auth") != auth.GetAuthKey() {
				// Authentication failed
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
