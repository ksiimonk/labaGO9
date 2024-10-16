package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	baseURL = "http://localhost:8000/users"
)

// User структура для представления пользователя
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Session токен для авторизации
var sessionToken string

// Выполнение GET запроса для получения всех пользователей
func getAllUsers() {
	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error getting users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var users []User
		json.NewDecoder(resp.Body).Decode(&users)
		fmt.Println("Users:")
		for _, user := range users {
			fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d\n", user.ID, user.Name, user.Email, user.Age)
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Error: %s - %s", resp.Status, body)
	}
}

// Выполнение GET запроса для получения конкретного пользователя
func getUser(id int) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", baseURL, id), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error getting user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user User
		json.NewDecoder(resp.Body).Decode(&user)
		fmt.Printf("User: %+v\n", user)
	} else {
		log.Printf("Error: %s", resp.Status)
	}
}

// Выполнение POST запроса для создания нового пользователя
func createUser(user User) {
	userData, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(userData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var newUser User
		json.NewDecoder(resp.Body).Decode(&newUser)
		fmt.Printf("User created: %+v\n", newUser)
	} else {
		log.Printf("Error: %s", resp.Status)
	}
}

// Выполнение PUT запроса для обновления пользователя
func updateUser(id int, user User) {
	user.ID = id
	userData, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%d", baseURL, id), bytes.NewBuffer(userData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error updating user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var updatedUser User
		json.NewDecoder(resp.Body).Decode(&updatedUser)
		fmt.Printf("User updated: %+v\n", updatedUser)
	} else {
		log.Printf("Error: %s", resp.Status)
	}
}

// Выполнение DELETE запроса для удаления пользователя
func deleteUser(id int) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d", baseURL, id), nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+sessionToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error deleting user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("User deleted")
	} else {
		log.Printf("Error: %s", resp.Status)
	}
}

func authorize(username, password string) {
	authData := map[string]string{"username": username, "password": password}
	jsonData, _ := json.Marshal(authData)

	req, err := http.NewRequest("POST", "http://localhost:8000/login", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error authorizing: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var token map[string]string
		json.NewDecoder(resp.Body).Decode(&token)
		sessionToken = token["token"]
		fmt.Println("Successfully authorized")
	} else {
		log.Printf("Authorization failed: %s", resp.Status)
	}
}

// Главное меню
func main() {
	// Пример авторизации
	authorize("user", "password")

	for {
		fmt.Println("\n1. Get All Users")
		fmt.Println("2. Get User")
		fmt.Println("3. Create User")
		fmt.Println("4. Update User")
		fmt.Println("5. Delete User")
		fmt.Println("6. Exit")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getAllUsers() // Запрос на получение всех пользователей
		case 2:
			var id int
			fmt.Print("Enter user ID: ")
			fmt.Scan(&id)
			getUser(id)
		case 3:
			var user User
			fmt.Print("Enter name: ")
			fmt.Scan(&user.Name)
			fmt.Print("Enter email: ")
			fmt.Scan(&user.Email)
			fmt.Print("Enter age: ")
			fmt.Scan(&user.Age)
			createUser(user)
		case 4:
			var id int
			var user User
			fmt.Print("Enter user ID: ")
			fmt.Scan(&id)
			fmt.Print("Enter new name: ")
			fmt.Scan(&user.Name)
			fmt.Print("Enter new email: ")
			fmt.Scan(&user.Email)
			fmt.Print("Enter new age: ")
			fmt.Scan(&user.Age)
			updateUser(id, user)
		case 5:
			var id int
			fmt.Print("Enter user ID: ")
			fmt.Scan(&id)
			deleteUser(id)
		case 6:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
