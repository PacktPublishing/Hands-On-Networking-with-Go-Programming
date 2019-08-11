package jsonerror

import "time"

type HTTPError struct {
	RequestID string    `json:"rid"`
	Time      time.Time `json:"t"`
	Issuer    string    `json:"iss"`
	ClientID  string    `json:"clientId"`
	UserID    string    `json:"userId"`
	Status    int       `json:"status"`
	Code      string    `json:"code"`
	Message   string    `json:"msg"`
}
