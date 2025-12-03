package order

import (
	"dz4/configs"
	"dz4/internal/models"
	"dz4/internal/user"
	"dz4/pkg/middleware"
	"dz4/pkg/req"
	"dz4/pkg/res"
	"net/http"
	"strconv"
	"time"

	"gorm.io/datatypes"
)

type OrderHandlerDeps struct {
	OrderRepository *OrderRepository
	UserRepository  *user.UserRepository
	Config          *configs.Config
}

type OrderHandler struct {
	OrderRepository *OrderRepository
	UserRepository  *user.UserRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderRepository: deps.OrderRepository,
		UserRepository: deps.UserRepository,
	}
	router.Handle("POST /order", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("GET /order/{id}", middleware.IsAuthed(handler.GetById(), deps.Config))
	router.Handle("GET /my-orders", middleware.IsAuthed(handler.GetByUserId(), deps.Config))
}

func (handler *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId, ok := r.Context().Value(middleware.ContextSessionKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := handler.UserRepository.FindBySessionId(sessionId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		body, err := req.HandleBody[CreateOrderRequest](&w, r)
		if err != nil {
			return
		}

		parsedDate, err := time.Parse("2006-01-02", body.Date)
		if err != nil {
			http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}

		order := &models.Order{
			UserId: user.ID,
			Date:   datatypes.Date(parsedDate),
		}

		createdOrder, err := handler.OrderRepository.Create(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(body.ProductIDs) > 0 {
			err = handler.OrderRepository.AddProductsToOrder(createdOrder.ID, body.ProductIDs)
			if err != nil {
				_ = handler.OrderRepository.Delete(createdOrder.ID)
				http.Error(w, "Failed to add products to order", http.StatusInternalServerError)
				return
			}
		}

		fullOrder, err := handler.OrderRepository.GetById(createdOrder.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := toOrderResponse(fullOrder)
		res.Json(w, response, http.StatusCreated)
	}
}

func (handler *OrderHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId, ok := r.Context().Value(middleware.ContextSessionKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := handler.UserRepository.FindBySessionId(sessionId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32) // десятичная система, uint32
		if err != nil {
			http.Error(w, "Invalid order ID", http.StatusBadRequest)
			return
		}

		order, err := handler.OrderRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		if order.UserId != user.ID {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		response := toOrderResponse(order)
		res.Json(w, response, http.StatusOK)
	}
}

func (handler *OrderHandler) GetByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionId, ok := r.Context().Value(middleware.ContextSessionKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := handler.UserRepository.FindBySessionId(sessionId)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		orders, err := handler.OrderRepository.GetByUserId(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := make([]OrderResponse, len(orders))
		for i, order := range orders {
			response[i] = toOrderResponse(&order)
		}

		res.Json(w, response, http.StatusOK)
	}
}

func toProductResponse(product models.Product) ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Images:      product.Images,
	}
}


func toOrderResponse(order *models.Order) OrderResponse {
	products := make([]ProductResponse, len(order.Products))
	for i, product := range order.Products {
		products[i] = toProductResponse(product)
	}

	return OrderResponse{
		ID:        order.ID,
		UserID:    order.UserId,
		Date:      time.Time(order.Date).Format("2006-01-02"),
		Products:  products,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		UpdatedAt: order.UpdatedAt.Format(time.RFC3339),
	}
}