package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strings"
	"errors"
	// "gopkg.in/validator.v2"
)

func listCarts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(carts)
}


func processRemoval(removeFromCartRequest RemoveFromCartRequest, carts map[int]map[int]CartItem) (error, map[int]map[int]CartItem) {
	err := errors.New("this item/user doesnt exist in the cart, cannot reduce")
	userId := removeFromCartRequest.UserID
	cartItems, userOk := carts[userId]
	// fmt.Println(carts)
	// fmt.Println(cartItems)
	// fmt.Println(userOk)
	if userOk {
		cartItem, itemOk := cartItems[removeFromCartRequest.CartItem.ItemID]
		if itemOk {
			// update the existing cart item
			cartItem.Amount -= removeFromCartRequest.CartItem.Amount
			if (cartItem.Amount <= 0) {
				delete(carts[userId], cartItem.ItemID)
				// w.WriteHeader(http.StatusOK)
				// return
			} else {
				cartItems[removeFromCartRequest.CartItem.ItemID] = cartItem
			}
			// fmt.Println(cartItem)
			// fmt.Println(itemOk)
			err = nil
		}
	}
	return err, carts
}

func removeFromCart(w http.ResponseWriter, r *http.Request) {
	// var itemID int
	// err := json.NewDecoder(r.Body).Decode(&itemID)
	var removeFromCartRequest RemoveFromCartRequest
	err := json.NewDecoder(r.Body).Decode(&removeFromCartRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err = addInventory(
		removeFromCartRequest.CartItem.ItemID,
		removeFromCartRequest.CartItem.Amount,
	)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err, carts = processRemoval(removeFromCartRequest, carts)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
	}

	// Remove from cart logic (not implemented in this example)
	w.WriteHeader(http.StatusOK)
}

func processAddition(addToCartRequest AddToCartRequest, carts map[int]map[int]CartItem) map[int]map[int]CartItem {
	userId := addToCartRequest.UserID
	cartItems, userOk := carts[userId]
	var item CartItem
	if !userOk {
		// create new item entry in cart
		carts[userId] = map[int]CartItem{}
		cartItems = carts[userId]
	}
	item, itemOk := cartItems[addToCartRequest.CartItem.ItemID]
	if itemOk {
		// update the existing cart item
		// item.CartItem.Amount += addToCartRequest.CartItem.Amount
		cartItem := CartItem{
			ItemID: item.ItemID,
			Amount: addToCartRequest.CartItem.Amount + item.Amount,
		}
		// cartItem.Amount += addToCartRequest.CartItem.Amount
		cartItems[addToCartRequest.CartItem.ItemID] = cartItem
	} else {
		// create new entry
		cartItems[addToCartRequest.CartItem.ItemID] = addToCartRequest.CartItem
	}
	carts[userId] = cartItems
	return carts
}

func addToCart(w http.ResponseWriter, r *http.Request) {
	var addToCartRequest AddToCartRequest
	err := json.NewDecoder(r.Body).Decode(&addToCartRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	err = reduceInventory(
		addToCartRequest.CartItem.ItemID,
		addToCartRequest.CartItem.Amount,
	)

	if err != nil {
		http.Error(w, "Bad Request, we do not have this much inventory", http.StatusUnauthorized)
		return
	}
	carts = processAddition(addToCartRequest, carts)

	// Add to cart logic (not implemented in this example)
	w.WriteHeader(http.StatusCreated)
}