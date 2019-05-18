package main

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"

	"github.com/dgrijalva/jwt-go"
)

const privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA1IE11W+JDl1+JhB5g2V3Ux/284zkyAN08s6OLMCHVIp5/7KA
WDTe23voRiKaPKDWNlh9gUZqkgKBb4kxhI5ctAgy6p22hPz6gsz82/0jMvJKojC4
Duv0vAc8r/xS/Uwq6xkBWs8KTqeQhkqInsthzj92FzHukQc3UUkrTq8Iu3nVTBzg
oBXZOvBxFW2swt/9tE4zhtyhAgDlnrpjVP6Va0MBJu+6bA1R18AKYSK9uNpMkof9
ltk44NstKgO5YpV9BdyNBlnU6CfpsoeWdveCJmSEPzv7YCHnQZdijg0DKCITZQxr
VVKUh561x0b5242kbWJOCUyLBOj6nkYXyF+qrQIDAQABAoIBAQCnN22XIAcnSKZl
aX1UydkVjgeTKoE0apPyJFt4F5/mBHlvnZSk1CWxbFUgK0ZXAvDNHuDTgweFEXes
vrY6apPEDteSCrx+9Vpi5s7qhMzX4BSef9u10jJoawF0MgdTzkXPbYPFYznnHq/5
HFlZKw0xcHqKUf46HQWIbx0m81DZw9vEtofBsNNvNBGp7GtbvBhMQoFdMHUsbD2N
Yf6unJb4KXixwtVXma6Wy2BQbJj1ifQ+xmZFak/130w1+64OY9+QMUa8zNKK2OHO
iVK4BrPYk3Vt49XsMRTLB6h3MDXvcvQ66x7jY8RunP1jr9bD8eK3ZX7uO6uyCa4Q
iNafQ10BAoGBAO9HWhYR2GAxVeIekkf1UwC6dc7MEioeIhxqP9rekEurVZRh42ko
zz5yD07CUEEz/iWZ+pqk/XYNTBv9si3PLqadN1b+o7XOgm67rVyp69rjJA2eKrlu
xOyMfXoeEi2CCLbHi7MTb6EdLk4MAjHMcFI6VOBt4bpvQiPQUBFpXEDhAoGBAONa
4IiwxIf6uUg3MRF248VYoFcmObLh1n75K2UGkqMzFUNA6lPaDreRhhXmJGrbHz0O
wjMZ9zc799HVnZr5YpQbidMGc+S0LrmqAI1zINIvpZ9RgYwMEzy/5qGLpYFt3Gx7
Mal7SdaVDdPbVLG4tzSqm3/0svQtxoeQ6utbGgdNAoGBAI4ZSZ6hqmY15lMK5MRn
JIviL+RHvOHWU1ucnZ9VXUwSzBf6qhrhaXIkOoMDUrXmMqAR+YmtQfjBnNliqFYc
HBBGfX7kakSmBz/LpQDKyI6NJfQQYj8NUVVJeZr0EMeF2bbyejw25qw/sCgZaZQ5
XNr4WT+PAea9/AFYzLQKZgcBAoGAY5g6xgZRgZPWuIjc6N6g9qFVU/f9zJvb37F9
TfssH2vQQ67bN7JNQiLwjwVLLLginheqAMK+JicR74zZRrs6cNEDdjrcZ/J6iYCs
T0qAtTKEJh+JVXUwtCsId/n5nZInvinVXn4QoXyYGxd4qYXWU67tAYeLISYwUtCr
6D/3Tf0CgYAm9LOsTzyAbrNExdM5Aj64R23DWaiH3q91z3pBB9wo99gcnCcdkm7n
2XWOLf3/103MaMy+TckAWIRdr3CLLuc4LevjsdM2T66wXDU9HoqVPjepiPYfzAq0
6D5hL7RE6SCte1/u4QNKfVL8smZg06joL2PVSN0IXG4VjzoJEJjFCg==
-----END RSA PRIVATE KEY-----`

func main() {
	pk, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		fmt.Println("failed to parse key:", err)
		os.Exit(1)
	}

	h := Handler{
		key: pk,
		// Using plain text passwords is a bad idea. Don't do it. :)
		users: map[string]string{
			"user1@example.com": "BadPassword1",
			"user2@example.com": "BadPassword2",
		},
	}

	http.ListenAndServe(":8000", h)
}

type Handler struct {
	key   *rsa.PrivateKey
	users map[string]string
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if ap, userExists := h.users[username]; !userExists || ap != password {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		Id:        uuid.Must(uuid.NewV4()).String(),
		Audience:  "any",
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().UTC().Add(time.Hour).Unix(),
		Issuer:    "my_issuer",
		Subject:   username,
	})

	s, err := t.SignedString(h.key)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(s))
}
