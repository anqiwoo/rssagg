package main

import (
	"github.com/anqiwoo/rssagg/internal/auth"
	"github.com/anqiwoo/rssagg/internal/database"
	"net/http"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Error getting API key")
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Error getting user")
			return
		}

		handler(w, r, user)
	}
}
