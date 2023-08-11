package main

import (
	// "errors"
	// "fmt"
	"testing"
)

func TestProcessRemoval(t *testing.T) {
	carts := make(map[int]map[int]CartItem)

	// Create a test user and cart item
	userID := 2
	itemID := 1
	initialAmount := 5
	removeAmount := 3

	// Initialize the cart with a test user and item
	cartItems := make(map[int]CartItem)
	cartItems[itemID] = CartItem{ItemID: itemID, Amount: initialAmount}
	carts[userID] = cartItems
	// fmt.Println(carts)

	// Create a removeFromCartRequest for testing
	removeFromCartRequest := RemoveFromCartRequest{
		UserID: userID,
		CartItem: CartItem{
			ItemID: itemID,
			Amount: removeAmount,
		},
	}

	// Call the function to be tested
	err, carts := processRemoval(removeFromCartRequest, carts)
	// fmt.Println(carts)

	// Check for errors
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Check if cart item was properly updated or deleted
	updatedCartItem, cartItemExists := carts[userID][itemID]
	if !cartItemExists {
		t.Errorf("Expected cart item to exist, but it was deleted")
	}
	if updatedCartItem.Amount != (initialAmount - removeAmount) {
		t.Errorf("Expected cart item amount to be %d, but got %d", initialAmount-removeAmount, updatedCartItem.Amount)
	}
}

func TestProcessAddition(t *testing.T) {
	carts := make(map[int]map[int]CartItem)

	// Create a test user and cart item
	userID := 123
	itemID := 456
	addAmount := 3
	initialAmount := 3
	cartItems := make(map[int]CartItem)
	cartItems[itemID] = CartItem{ItemID: itemID, Amount: initialAmount}
	carts[userID] = cartItems

	// Create an addToCartRequest for testing
	addToCartRequest := AddToCartRequest{
		UserID: userID,
		CartItem: CartItem{
			ItemID: itemID,
			Amount: addAmount,
		},
	}

	// Call the function to be tested
	updatedCarts := processAddition(addToCartRequest, carts)

	// Check if the cart was updated as expected
	cartItems, userOk := updatedCarts[userID]
	if !userOk {
		t.Errorf("Expected user cart to exist, but it wasn't created")
	}

	updatedCartItem, cartItemExists := cartItems[itemID]
	if !cartItemExists {
		t.Errorf("Expected cart item to exist, but it wasn't created")
	}

	
	expectedAmount := 0
	if cartItemExists {
		expectedAmount = initialAmount + addAmount
	}

	if updatedCartItem.Amount != expectedAmount {
		t.Errorf("Expected cart item amount to be %d, but got %d", expectedAmount, updatedCartItem.Amount)
	}
}