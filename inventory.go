package main

import (
	"encoding/json"
	// "fmt"
	"net/http"
	// "strings"
	"errors"
	// "gopkg.in/validator.v2"
)

func listItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(items)
}


func addInventory(itemId int, amount int) error {
	err := errors.New("this item doesnt exist")
	item, ok := items[itemId]
	// If the key exists
	if ok {
	    // Do something
	    item.Stock += amount
	    items[itemId] = item
	    err = nil
	}
    return err
}

func reduceInventory(itemId int, amount int) error {
	err := errors.New("this much inventory is not available")
	item, ok := items[itemId]
	// If the key exists
	if ok {
	    // Do something
	    if (item.Stock - amount >= 0) {
	    	err = nil
    		item.Stock -= amount
    		if item.Stock == 0 {
    			delete(items, itemId)
    			return err
    		}
    		items[itemId] = item
    	}
	}
    return err
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// items = append(items, newItem)
	items[newItem.ID] = newItem
	w.WriteHeader(http.StatusCreated)
}
