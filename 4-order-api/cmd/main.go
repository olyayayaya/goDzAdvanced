package main

import (
	"dz4/configs"
	"dz4/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf) //инициализируем базу данных
	router := http.NewServeMux()

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
 
	fmt.Println("server is lixtening on port 8081")
	server.ListenAndServe()
}
 