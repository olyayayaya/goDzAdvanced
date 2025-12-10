package main

import (
	"bytes"
	"dz4/internal/auth"
	"dz4/internal/models"
	"dz4/internal/order"
	"dz4/internal/user"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func createTestUser(db *gorm.DB) (uint, int, string) {
	user := &user.User{
		PhoneNumber: "+79991112233",
		SessionId:   "123qwertyuiopasd",
		Code:        1234,
	}
	db.Create(user)
	return user.ID, user.Code, user.SessionId
}

func createTestProduct(db *gorm.DB) uint {
	product := &models.Product{
		Name:        "Test Product",
		Description: "Test product description",
		Images:      []string{"test.jpg"},
	}
	db.Create(product)
	return product.ID
}

func removeData(db *gorm.DB, productID uint) {
	db.Exec("DELETE FROM order_products")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products WHERE id = ?", productID)
	db.Exec("DELETE FROM users WHERE phone_number = ?", "+79991112233")
}

func TestCreateOrder(t *testing.T) {
	// Prepare
	db := initDb()
	userId, userCode, userSessionId := createTestUser(db)
	productId := createTestProduct(db)

	// создаем тестовый сервер
	ts := httptest.NewServer(App())
	defer ts.Close() // закрываем тестовый сервер

	data, _ := json.Marshal(&auth.ValidationCodeRequest{
		SessionId: userSessionId,
		Code:      userCode,
	})

	res, err := http.Post(ts.URL+"/auth/checkValidationCode", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// тут мы уже залогинились и получили токен
	var resData auth.ValidationCodeResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}

	// Отладочный вывод
	t.Logf("Token received: %s", resData.Token)
	t.Logf("User ID: %d", userId)

	if resData.Token == "" {
		t.Fatalf("token empty")
	}

	orderData, _ := json.Marshal(&order.CreateOrderRequest{
		ProductIDs: []uint{productId},
		Date:       "2025-12-10",
	})

	req, err := http.NewRequest("POST", ts.URL+"/order", bytes.NewReader(orderData))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+resData.Token)

	orderRes, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer orderRes.Body.Close()

	if orderRes.StatusCode == 401 {
		errorBody, _ := io.ReadAll(orderRes.Body)
		t.Logf("401 Error body: %s", string(errorBody))
		t.Logf("Response headers: %v", orderRes.Header)
	}

	if orderRes.StatusCode != 201 {
		t.Fatalf("expected %d got %d", 201, orderRes.StatusCode)
	}

	orderBody, err := io.ReadAll(orderRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resOrderData order.OrderResponse
	err = json.Unmarshal(orderBody, &resOrderData)
	if err != nil {
		t.Fatal(err)
	}
	if resOrderData.UserID == userId {
		t.Fatalf("user error")
	}

	removeData(db, productId)
}
