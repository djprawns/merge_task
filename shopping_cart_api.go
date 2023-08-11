package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strings"
	// "errors"
	"bytes"
	"os"
	"os/signal"
	"syscall"
	"time"
	_ "gopkg.in/validator.v2"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
	Role     Role   `json:"role" validate:"nonnil"`
}

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" validate:"nonzero"`
	Price float64 `json:"price" validate:"min=0"`
	Stock int     `json:"stock" validate:"min=0"`
}

type CartItem struct {
	ItemID int `json:"item_id" validate:"nonnil"`
	Amount int `json:"amount" validate:"min=1"`
}

type AddToCartRequest struct {
	CartItem   CartItem `json:"cart_item" validate:"nonnil"`
	UserID     int      `json:"user_id" validate:"nonzero"`
}

type RemoveFromCartRequest struct {
	CartItem   CartItem `json:"cart_item" validate:"nonnil"`
	UserID     int      `json:"user_id" validate:"nonzero"`
}


var users = []User{
	{ID: 1, Username: "admin", Password: "admin", Role: AdminRole},
	{ID: 2, Username: "user", Password: "user", Role: UserRole},
}

var items = map[int]Item{
	1 : {ID: 1, Name: "Item 1", Price: 10.0, Stock: 5},
	2 : {ID: 2, Name: "Item 2", Price: 15.0, Stock: 10},
}

// userid - cart_item list
// var carts = map[int][]CartItem{}
var carts = map[int]map[int]CartItem{}

func main() {
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/admin/add-item", authenticate(AdminRole, addItem))
	http.HandleFunc("/admin/suspend-user", authenticate(AdminRole, suspendUser))
	// this lists all the carts
	http.HandleFunc("/admin/list-carts", authenticate(AdminRole, listCarts))
	http.HandleFunc("/user/list-items", authenticate(UserRole, listItems))
	http.HandleFunc("/user/add-to-cart", authenticate(UserRole, addToCart))
	http.HandleFunc("/user/remove-from-cart", authenticate(UserRole, removeFromCart))

	port := 8080
	serverDone := make(chan struct{})

	go func() {
		fmt.Printf("Server started on port %d\n", port)
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		close(serverDone)
	}()

	// Wait for the server to start
	<-serverDone

	// Run a method after the server starts
	afterServerStart()

	// Wait for a termination signal without blocking
	waitForTerminationSignal()
}

func afterServerStart() {
	e2e_test()
	// fmt.Println("Method executed after the server starts")
}

func waitForTerminationSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals
	fmt.Println("Received termination signal. Shutting down...")
	// Add cleanup or other termination logic here if needed
	time.Sleep(time.Second) // Give some time for cleanup if necessary
	os.Exit(0)
}


func e2e_test() {
	// E2E tests done here
	// 1. register user,
	// 2. login user,
	// 3. then use the user to list inventory/items, 
	// 4. then add to the users cart, 
	// 5. and display the carts contents
	// 6. then remove from the users cart, 
	// 7. and display the carts contents
	// 8. let the admin add an item
	// 9. list the inventory/items available
	// 10. finally suspend the user


	baseURL := "http://localhost:8080"

	// 1. Register user
	registerUser := User{
		Username: "newuser",
		Password: "newpassword",
		Role:     UserRole,
	}
	registerUserJSON, _ := json.Marshal(registerUser)
	resp, _ := http.Post(baseURL+"/register", "application/json", bytes.NewBuffer(registerUserJSON))
	fmt.Println("Register User Response:", resp.Status)

	// 2. Login user
	loginCredentials := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "newuser",
		Password: "newpassword",
	}
	loginCredentialsJSON, _ := json.Marshal(loginCredentials)
	resp, _ = http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(loginCredentialsJSON))
	fmt.Println("Login User Response:", resp.Status)
	token := resp.Header.Get("Authorization")

	// 2.1. Login Admin
	loginCredentialsAdmin := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "admin",
		Password: "admin",
	}
	loginCredentialsAdminJSON, _ := json.Marshal(loginCredentialsAdmin)
	resp, _ = http.Post(baseURL+"/login", "application/json", bytes.NewBuffer(loginCredentialsJSON))
	fmt.Println("Login Admin Response:", resp.Status)
	loginCredentialsAdminToken := resp.Header.Get("Authorization")
	fmt.Println("Login Admin Token:", resp.loginCredentialsAdminToken)

	// 3. List inventory/items
	resp, _ = http.Get(baseURL + "/user/list-items")
	fmt.Println("List Items Response:", resp.Status)

	// 4. Add to the user's cart
	addToCartRequest := AddToCartRequest{
		CartItem: CartItem{
			ItemID: 1, // Replace with actual item ID
			Amount: 2,
		},
		UserID: 2, // Replace with the actual user ID
	}
	addToCartJSON, _ := json.Marshal(addToCartRequest)
	req, _ := http.NewRequest("POST", baseURL+"/user/add-to-cart", bytes.NewBuffer(addToCartJSON))
	req.Header.Set("Authorization", token)
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("Add to Cart Response:", resp.Status)

	// 5. Display cart's contents
	resp, _ = http.Get(baseURL + "/admin/list-carts")
	fmt.Println("List Carts Response:", resp.Status)

	// 6. Remove from the user's cart
	removeFromCartRequest := RemoveFromCartRequest{
		CartItem: CartItem{
			ItemID: 1, // Replace with actual item ID
			Amount: 1,
		},
		UserID: 2, // Replace with the actual user ID
	}
	removeFromCartJSON, _ := json.Marshal(removeFromCartRequest)
	req, _ = http.NewRequest("POST", baseURL+"/user/remove-from-cart", bytes.NewBuffer(removeFromCartJSON))
	req.Header.Set("Authorization", token)
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("Remove from Cart Response:", resp.Status)

	// 7. Display cart's contents
	resp, _ = http.Get(baseURL + "/admin/list-carts")
	fmt.Println("List Carts Response:", resp.Status)

	// 8. Admin adds an item
	newItem := Item{
		ID:    3, // Replace with a unique item ID
		Name:  "New Item",
		Price: 20.0,
		Stock: 10,
	}
	newItemJSON, _ := json.Marshal(newItem)
	req, _ = http.NewRequest("POST", baseURL+"/admin/add-item", bytes.NewBuffer(newItemJSON))
	req.Header.Set("Authorization", token)
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("Add Item Response:", resp.Status)

	// 9. List inventory/items available
	resp, _ = http.Get(baseURL + "/user/list-items")
	fmt.Println("List Items Response:", resp.Status)

	// 10. Suspend the user
	userIDToRemove := 2 // Replace with the actual user ID
	suspendJSON, _ := json.Marshal(userIDToRemove)
	req, _ = http.NewRequest("POST", baseURL+"/admin/suspend-user", bytes.NewBuffer(suspendJSON))
	req.Header.Set("Authorization", token)
	resp, _ = http.DefaultClient.Do(req)
	fmt.Println("Suspend User Response:", resp.Status)
}


