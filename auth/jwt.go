package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenJwtGeneral(email string) (string, error) {
	expirationTime := time.Now().AddDate(0, 3, 0)
	retToken, err := genJwt(expirationTime, email)
	if err != nil {
		return "", err
	}
	return retToken, nil
}

func GenJwtVerify(email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	retToken, err := genJwt(expirationTime, email)
	if err != nil {
		return "", err
	}
	return retToken, nil
}

func genJwt(expirationTime time.Time, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"expTime": expirationTime.Format(time.RFC3339Nano),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := jwtSecret()
	retToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return retToken, nil
}

func CheckJWT(tokenString string) (bool, string, error) {
	secret := jwtSecret()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { return secret, nil })
	if err != nil {
		fmt.Println(err)
		return false, "", nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, "", nil
	}
	exp, ok := claims["expTime"].(string)
	if !ok {
		return false, "", errors.New("can not get exp")
	}
	expTime, err := time.Parse(time.RFC3339Nano, exp)
	if err != nil {
		return false, "", err
	}
	now := time.Now()
	if now.After(expTime) {
		return false, "", errors.New("token expired")
	}
	email, ok := claims["email"].(string)
	if !ok {
		return false, "", errors.New("can not get email")
	}
	return true, email, nil
}
