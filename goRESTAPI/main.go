package main

import (
	"math/rand"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/gorilla/mux"

	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"html/template"
)

type ViewData struct {
	Answer string
}

var mySigningKey = []byte("secret_key_5c0ea8459c058006f38fd53d1c912d1e_kMsrKf7124389351789c6c544ed71cf91b77b")

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "Elliot Forbes"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("secret_key_5c0ea8459c058006f38fd53d1c912d1e_kMsrKf7124389351789c6c544ed71cf91b77b")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}

type Book struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Stars  int    `json:"stars"`
}

var books = []Book{
	{Id: "1", Title: "War and Peace", Author: "Tolstoi", Stars: 4},
	{Id: "2", Title: "Genius", Author: "Draiser", Stars: 5},
	{Id: "3", Title: "Doctor Jivago", Author: "Pasternak", Stars: 3},
}

func homepage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("templates/index.html")
		data := ViewData{Answer: "Enter password"}
		tmpl.Execute(w, data)
	}
	if r.Method == "POST" {
		if (r.FormValue("login") == "root") && (r.FormValue("password") == "root") {
			a, _ := GenerateJWT()
			tmpl, _ := template.ParseFiles("templates/index.html")
			data := ViewData{Answer: a}
			tmpl.Execute(w, data)
		} else {
			w.Write([]byte("Incorrect!"))
		}
	}
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range books {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing Authorization Header"))
		return
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	_, err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Error verifying JWT token: " + err.Error()))
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		for index, item := range books {
			if item.Id == params["id"] {
				books = append(books[:index], books[index+1:]...)
				break
			}
			json.NewEncoder(w).Encode(books)
		}

	}

}
func createBook(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Missing Authorization Header"))
		return
	}
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	_, err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Error verifying JWT token: " + err.Error()))
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		var book Book
		_ = json.NewDecoder(r.Body).Decode(&book)
		book.Id = strconv.Itoa(rand.Intn(1000000))
		books = append(books, book)
		json.NewEncoder(w).Encode(book)
	}
}
func main() {
	r := mux.NewRouter()
	// a, _ := GenerateJWT()
	// fmt.Println(a)
	//r.HandleFunc("/", homepage).Methods("GET", "POST")
	r.HandleFunc("/", homepage).Methods("GET", "POST")
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	http.ListenAndServe(":8000", r)
}
