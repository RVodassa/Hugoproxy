package server

// import (
// 	"encoding/json"
// 	"net/http"

// 	"golang.org/x/crypto/bcrypt"
// )

// type User struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// var users = make(map[string]string) // email -> hashed password

// func RegisterHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User

// 	// Декодируем данные пользователя из тела запроса
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}

// 	// Проверяем, существует ли пользователь с таким email
// 	if _, exists := users[user.Email]; exists {
// 		http.Error(w, "User already exists", http.StatusConflict)
// 		return
// 	}

// 	// Хешируем пароль
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		http.Error(w, "Error hashing password", http.StatusInternalServerError)
// 		return
// 	}

// 	// Сохраняем пользователя (email и хеш пароля)
// 	users[user.Email] = string(hashedPassword)

// 	// Возвращаем успешный ответ
// 	w.Write([]byte("User registered successfully"))
// }
