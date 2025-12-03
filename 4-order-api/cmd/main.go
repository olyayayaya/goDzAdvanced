package main

import (
	"dz4/configs"
	"dz4/internal/auth"
	"dz4/internal/order"
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
	orderRepository := order.NewOrderRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)

	// Handlers
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
	})
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderRepository: orderRepository,
		UserRepository: userRepository,
		Config: conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("server is lixtening on port 8081")
	server.ListenAndServe()
}
