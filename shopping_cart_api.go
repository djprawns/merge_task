package main

import (
	// "encoding/json"
	"fmt"
	"net/http"
	// "strings"
	// "errors"
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
	http.HandleFunc("/admin/list-carts", authenticate(AdminRole, listCarts))
	http.HandleFunc("/user/list-items", authenticate(UserRole, listItems))
	http.HandleFunc("/user/add-to-cart", authenticate(UserRole, addToCart))
	http.HandleFunc("/user/remove-from-cart", authenticate(UserRole, removeFromCart))

	port := 8080
	fmt.Printf("Server started on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}







