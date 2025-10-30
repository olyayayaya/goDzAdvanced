package main

import (
	"dz3/configs"
	"dz3/internal/verify"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()

	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		Config: conf,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
 
	fmt.Println("server is lixtening on port 8081")
	server.ListenAndServe()
}
 