package handler

import (
	"achyuBackend/auth"
	"achyuBackend/db"
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"net/http"
)

func OAuthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.OAuthLoginURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func OAuthCallBack(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if !auth.CheckStateString(state) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	code := r.FormValue("code")
	token, err := auth.CodeExchange(code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	email, ok := userInfo["email"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jwtToken, err := auth.GenJwtGeneral(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	exists, err := db.CheckExistMail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	displayName, ok := userInfo["given_name"].(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		err = db.UserInsert(email, "", displayName, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	cookie := &http.Cookie{
		Name:     "token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   false,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "http://localhost:3000/home", http.StatusFound)
}

/*
map[
	email:hirogoshawk3249@gmail.com
	family_name:t
	given_name:h
	id:
	name:h t
	picture:https://lh3.googleusercontent.com/a/
	verified_email:true
]
*/
