package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	b, err := ioutil.ReadFile("example_public.pem")
	if err != nil {
		fmt.Println("failed to read key:", err)
		os.Exit(1)
	}
	pk, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		fmt.Println("failed to parse key:", err)
		os.Exit(1)
	}

	v := JWTValidator{
		key:  pk,
		next: WhoAmI{},
	}

	err = http.ListenAndServe(":8001", v)
	if err != nil {
		fmt.Println(err)
	}
}

type claimsKey string

var claimsContext claimsKey = "claimsContext"

type JWTValidator struct {
	key  *rsa.PublicKey
	next http.Handler
}

func (v JWTValidator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the token from the request.
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "missing Authorization header", http.StatusUnauthorized)
		return
	}

	s := strings.TrimPrefix(auth, "Bearer ")

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return v.key, nil
	})
	if err != nil {
		http.Error(w, "auth error", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	// Add the user's details to the context.
	r = r.WithContext(context.WithValue(r.Context(), claimsContext, claims))
	v.next.ServeHTTP(w, r)
}

type WhoAmI struct {
}

func (wai WhoAmI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cv := r.Context().Value(claimsContext)
	if cv == nil {
		http.Error(w, "can't find user", http.StatusUnauthorized)
		return
	}
	claims, ok := cv.(jwt.MapClaims)
	if !ok {
		http.Error(w, "claims aren't of that type", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("%+v\n", claims["sub"])))
}
