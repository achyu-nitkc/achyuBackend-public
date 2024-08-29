package handler

import (
	"achyuBackend/auth"
	"achyuBackend/db"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type requestLoginBody struct {
	Email    string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body requestLoginBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyEmail := body.Email
	bodyPassword := body.Password
	hash := sha256.New()
	hash.Write([]byte(bodyPassword))
	hashed := hash.Sum(nil)
	passwordHash := hex.EncodeToString(hashed)
	check, err := db.CheckMailHash(bodyEmail, passwordHash)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := auth.GenJwtGeneral(bodyEmail)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", token)
	return
}

type requestCreateBody struct {
	Email       string
	Password    string
	DisplayName string
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body requestCreateBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyEmail := body.Email
	bodyPassword := body.Password
	bodyDisplayName := body.DisplayName
	hash := sha256.New()
	hash.Write([]byte(bodyPassword))
	hashed := hash.Sum(nil)
	passwordHash := hex.EncodeToString(hashed)
	checkUser, err := db.CheckExistMail(bodyEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if checkUser {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	checkVerify, err := db.CheckExistVerify(bodyEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if checkVerify {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	b := make([]byte, 8)
	for i := range b {
		b[i] = '0' + byte(rand.Intn(10))
	}
	verifyCode := string(b)
	err = db.VerifyInsert(bodyEmail, passwordHash, bodyDisplayName, verifyCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = auth.SendVerifyCode(bodyEmail, verifyCode)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token, err := auth.GenJwtVerify(bodyEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", token)
}

type requestVerifyBody struct {
	Token string
	Code  string
}

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body requestVerifyBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyToken := body.Token
	bodyCode := body.Code
	check, email, err := auth.CheckJWT(bodyToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	check, err = db.VerifyCodeCmp(email, bodyCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := auth.GenJwtGeneral(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = db.MoveData(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", token)
}
