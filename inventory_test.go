package main

import (
	"fmt"
	"errors"
	"testing"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func TestAddInventory(t *testing.T) {
	// Initialize test data
	items = map[int]Item{
		1: {ID: 1, Stock: 10},
		2: {ID: 2, Stock: 5},
	}

	// Test cases
	testCases := []struct {
		itemID       int
		amount       int
		expectedErr  error
		expectedItem Item
	}{
		// Valid case: Item exists
		{1, 5, nil, Item{ID: 1, Stock: 15}},

		// Invalid case: Item doesn't exist
		{3, 5, errors.New("this item doesnt exist"), Item{ID: 0, Stock: 0}},
	}

	for _, tc := range testCases {
		err := addInventory(tc.itemID, tc.amount)

		// Check if the error matches the expected error
		if (err == nil && tc.expectedErr != nil) || (err != nil && tc.expectedErr == nil) || (err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
			t.Errorf("addInventory(%d, %d) error = %v, expected %v", tc.itemID, tc.amount, err, tc.expectedErr)
		}

		// Check if the item stock matches the expected stock
		item := items[tc.itemID]
		if item.Stock != tc.expectedItem.Stock {
			t.Errorf("addInventory(%d, %d) stock = %d, expected %d", tc.itemID, tc.amount, item.Stock, tc.expectedItem.Stock)
		}
	}
}

func TestReduceInventory(t *testing.T) {
	// Initialize test data
	items = map[int]Item{
		1: {ID: 1, Stock: 10},
		2: {ID: 2, Stock: 5},
		3: {ID: 3, Stock: 0},
	}

	// Test cases
	testCases := []struct {
		itemID       int
		amount       int
		expectedErr  error
		expectedItem Item
	}{
		// Valid case: Sufficient stock
		{1, 5, nil, Item{ID: 1, Stock: 5}},

		// Valid case: Exact stock, item should be deleted
		{2, 5, nil, Item{ID: 2, Stock: 0}},

		// Invalid case: Insufficient stock
		{1, 10, errors.New("this much inventory is not available"), Item{ID: 1, Stock: 5}},

		// Invalid case: Item not found
		{4, 1, errors.New("this much inventory is not available"), Item{ID: 0, Stock: 0}},
	}

	for _, tc := range testCases {
		err := reduceInventory(tc.itemID, tc.amount)

		// Check if the error matches the expected error
		if (err == nil && tc.expectedErr != nil) || (err != nil && tc.expectedErr == nil) || (err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error()) {
			t.Errorf("reduceInventory(%d, %d) error = %v, expected %v", tc.itemID, tc.amount, err, tc.expectedErr)
		}

		// Check if the item stock matches the expected stock
		item := items[tc.itemID]
		if item.Stock != tc.expectedItem.Stock {
			t.Errorf("reduceInventory(%d, %d) stock = %d, expected %d", tc.itemID, tc.amount, item.Stock, tc.expectedItem.Stock)
		}
	}
}


func TestAddItem(t *testing.T) {
	// Initialize the test server
	ts := httptest.NewServer(http.HandlerFunc(addItem))
	defer ts.Close()

	// Test cases
	testCases := []struct {
		name          string
		payload       interface{} // Request payload (item)
		expectedCode  int         // Expected HTTP response code
		expectedItems map[int]Item // Expected state of items after adding
	}{
		// Valid case
		{
			name: "ValidItem",
			payload: Item{
				ID:    4,
				Name: "asd",
				Price: 15,
				Stock: 15,
			},
			expectedCode: http.StatusCreated,
			expectedItems: map[int]Item{
				1: {ID: 1, Stock: 5},
				0: {ID: 0, Stock: 0},
				4: {ID: 4, Name: "asd", Price: 15, Stock: 15},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert payload to JSON
			payloadBytes, _ := json.Marshal(tc.payload)
			client := &http.Client{}

			// Send POST request to the test server
			req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(payloadBytes))
			if err != nil {
				fmt.Println(err)
				return
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Authorization", "admin:admin")
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send POST request: %v", err)
			}
			defer resp.Body.Close()

			// Check HTTP response code
			if resp.StatusCode != tc.expectedCode {
				t.Errorf("Expected status code %d, but got %d", tc.expectedCode, resp.StatusCode)
			}
			// fmt.Println(items)
			// Check the state of items
			for id, expectedItem := range tc.expectedItems {
				actualItem := items[id]
				// fmt.Println(actualItem)
				// fmt.Println(expectedItem)
				if actualItem != expectedItem {
					t.Errorf("Expected item %+v, but got %+v", expectedItem, actualItem)
				}
			}
		})
	}
}
