package handler

import (
	"achyuBackend/auth"
	"achyuBackend/db"
	"achyuBackend/yolp"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type requestPost struct {
	Token     string
	Latitude  float64
	Longitude float64
	Content   string
	ImageURL  string
	Where     string
}

func PostPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body requestPost
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyToken := body.Token
	bodyLatitude := body.Latitude
	bodyLongitude := body.Longitude
	bodyContent := body.Content
	bodyImageURL := body.ImageURL
	bodyWhere := body.Where
	check, email, err := auth.CheckJWT(bodyToken)
	if err != nil {
		fmt.Println("Error @ CheckJWT", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		ip = strings.TrimSpace(ips[0])
	}
	displayName, err := db.GetDisplayName(email)
	if err != nil {
		fmt.Println("Error @ GetDisplayName", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postId := time.Now().UnixNano() / int64(time.Millisecond)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, postId)
	hash := md5.Sum([]byte(email))
	strHash := hex.EncodeToString(hash[:])
	strId := base64.StdEncoding.EncodeToString(buf.Bytes()) + strHash
	addressData, err := yolp.GetAddressData(bodyLatitude, bodyLongitude)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error @ yolp.GetAddressData", err)
		return
	}

	err = db.InsertPost(strId, email, ip, displayName, bodyContent, bodyImageURL, bodyWhere, bodyLatitude, bodyLongitude, addressData)
	if err != nil {
		fmt.Println("Error @ db.InsertPost", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	return
}

type requestDelete struct {
	Token  string
	PostId string
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body requestDelete
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyToken := body.Token
	bodyPostId := body.PostId
	check, email, err := auth.CheckJWT(bodyToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = db.MovePost(email, bodyPostId)
	if err != nil {
		fmt.Println("Error @ db.MovePost", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type requestGet struct {
	Token     string
	Latitude  float64
	Longitude float64
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var body requestGet
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyToken := body.Token
	bodyLatitude := body.Latitude
	bodyLongitude := body.Longitude
	check, _, err := auth.CheckJWT(bodyToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !check {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	postData, err := db.GetPost(bodyLatitude, bodyLongitude)
	if err != nil {
		fmt.Println("Error @ db.GetPost", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("Error @ json.Marshal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}
