package middleware

import (
	"context"
	"net/http"
)

func GetID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		ctx := context.WithValue(r.Context(), "UserID", r.FormValue("id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
