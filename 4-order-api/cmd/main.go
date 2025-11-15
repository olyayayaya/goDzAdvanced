package main

import (
	"dz4/configs"
	"dz4/internal/auth"
	"dz4/internal/product"
	"dz4/internal/user"
	"dz4/pkg/db"
	"dz4/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf) //инициализируем базу данных
	router := http.NewServeMux()

	// Repositories
	productRepository := product.NewProductRepository(db)
	userRepository := user.NewUserRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
		AuthService: authService,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.Logging(router),
	}

	fmt.Println("server is lixtening on port 8081")
	server.ListenAndServe()
}
