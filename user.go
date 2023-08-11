package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"
    // "errors"
    validator "gopkg.in/validator.v2"
)

// this method should send successful message
// needs to support user registry in specific format
// 1. incremental id for user
// 2. unique user
// 3. pass role in the request body as well, 
//    and indicate possible roles when sending response
func register(w http.ResponseWriter, r *http.Request) {
    var newUser User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    fmt.Println("error1")
    if errs := validator.Validate(newUser); errs != nil {
        fmt.Println(errs)
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
        // values not valid, return reason/s for error
    }
    fmt.Println(err)
    fmt.Println("error2")

    users = append(users, newUser)
    response, err := json.Marshal(users)
    fmt.Println(users)
    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
    w.Write(response)
}

func login(w http.ResponseWriter, r *http.Request) {
    var credentials struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&credentials)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    var user User
    for _, u := range users {
        if u.Username == credentials.Username && u.Password == credentials.Password {
            user = u
            break
        }
    }

    if user.Username == "" {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token := fmt.Sprintf("%s:%s", user.Username, user.Role)
    w.Header().Set("Authorization", token)
    w.WriteHeader(http.StatusOK)
}


func removeUser(users []User, removeID int) ([]User, int) {
	indexToRemove := -1
    for i, user := range users {
        if user.ID == removeID {
            indexToRemove = i
            break
        }
    }

    // If the item was found, remove it from the slice
    if indexToRemove >= 0 {
        users = append(users[:indexToRemove], users[indexToRemove+1:]...)
    }
    return users, indexToRemove
}

func authenticate(requiredRole Role, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(token, ":")
        if len(parts) != 2 {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        username := parts[0]
        role := Role(parts[1])
        fmt.Println(role)

        var user User
        for _, u := range users {
            if u.Username == username {
                user = u
                break
            }
        }

        if user.Username == "" || user.Role != requiredRole {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    }
}

func suspendUser(w http.ResponseWriter, r *http.Request) {
    var userID int
    err := json.NewDecoder(r.Body).Decode(&userID)
    if err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Suspend user logic (not implemented in this example)
    users, _ = removeUser(users, userID)
    // if (indexToRemove < 0) {
    //  http.Error(w, "user with this id does not exist", http.StatusBadRequest)
    //  return
    // }
    w.WriteHeader(http.StatusOK)
}