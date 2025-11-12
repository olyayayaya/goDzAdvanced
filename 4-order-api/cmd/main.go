package main

import (
	"dz4/configs"
	"dz4/internal/product"
	"dz4/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf) //инициализируем базу данных
	router := http.NewServeMux()

	// Repositories
	productRepository := product.NewProductRepository(db)

	// Handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("server is lixtening on port 8081")
	server.ListenAndServe()
}
