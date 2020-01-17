package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) authMiddleware(v http.HandlerFunc) http.HandlerFunc {
	type Data struct {
		User User `json:"user"`
	}
	client := &http.Client{}
	authURL := fmt.Sprintf("%s/api/auth/user/", s.config.AppServer)

	return func(w http.ResponseWriter, r *http.Request) {
		authReq, err := http.NewRequest(http.MethodGet, authURL, nil)
		if err != nil {
			http.Error(w, "Auth error", http.StatusInternalServerError)
			return
		}
		authReq.Header.Set("Host", r.Host)
		authReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
		for _, cookie := range r.Cookies() {
			authReq.AddCookie(cookie)
		}
		resp, err := client.Do(authReq)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var data Data
			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), contextKeyUser, &data.User)
			r = r.WithContext(ctx)
		}
		v(w, r)
	}
}

func (s *Server) loginRequired(v http.HandlerFunc) http.HandlerFunc {
	return s.authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(contextKeyUser).(*User)
		if !ok || user.IsGuest {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		v(w, r)
	})
}
