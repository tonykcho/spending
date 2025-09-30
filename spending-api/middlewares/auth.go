package middlewares

import (
	"encoding/base64"
	"net/http"
	"spending/utils"
	"strings"
)

// Endpoint to skip authentication
var ExemptedPaths = map[string]bool{
	"/metrics": true,
	"/health":  true,
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ExemptedPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Basic") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(header, "Basic "))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 {
			http.Error(w, "bad credentials format", http.StatusBadRequest)
			return
		}

		username, password := parts[0], parts[1]
		if username != utils.GetBasicAuthUser() || password != utils.GetBasicAuthPassword() {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
