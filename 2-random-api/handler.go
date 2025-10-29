package main

import (
	"math/rand"
	"net/http"
	"strconv"
)

type Handler struct{}

func NewHandler(router *http.ServeMux) {
	handler := &Handler{}
	router.HandleFunc("/", handler.Hello())
}

func (handler *Handler) Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		randNum := rand.Intn(6) + 1
		w.Write([]byte(strconv.Itoa(randNum)))
	}
}
