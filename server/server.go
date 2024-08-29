package server

import (
	"achyuBackend/handler"
	"fmt"
	"net/http"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Start() {
	http.Handle("/login", enableCORS(http.HandlerFunc(handler.LoginHandler)))
	http.Handle("/create", enableCORS(http.HandlerFunc(handler.CreateHandler)))
	http.Handle("/verify", enableCORS(http.HandlerFunc(handler.VerifyHandler)))
	http.Handle("/post", enableCORS(http.HandlerFunc(handler.PostPostHandler)))
	http.Handle("/delete", enableCORS(http.HandlerFunc(handler.DeletePostHandler)))
	http.Handle("/get", enableCORS(http.HandlerFunc(handler.GetPostHandler)))

	http.HandleFunc("/Oauth", handler.OAuthGoogleLogin)
	http.HandleFunc("/OauthCallback", handler.OAuthCallBack)

	fmt.Println("Start Server @ http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
