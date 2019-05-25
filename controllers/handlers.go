package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gyan1230/asat/config"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

//ShowAll :
func ShowAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	usr, err := Allusers(r)
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usr)
	return
}

//Register ...
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("content-type", "application/json")
	var person User
	_ = json.NewDecoder(r.Body).Decode(&person)
	p, err := GetUser(r.Context(), person.Email)
	if err != nil {
		bs, _ := bcrypt.GenerateFromPassword([]byte(person.Password), bcrypt.MinCost)
		person.Password = string(bs)
		collection := config.Client.Database("userDb").Collection("user")
		result, _ := collection.InsertOne(r.Context(), person)
		json.NewEncoder(w).Encode(result)
		w.WriteHeader(http.StatusCreated)
		return
	}
	if p != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User already exist"))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("User could not be created"))
	return
}

//Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	if alreadyLoggedIn(w, r) {
		log.Println("Already login return to index...")
		w.Header().Set("content-type", "application/json")

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// is there a Email?
	var person User
	_ = json.NewDecoder(r.Body).Decode(&person)
	u, err := GetUser(r.Context(), person.Email) // return user (if present), nil in u,err OR nil,err (if not present user) in u,err
	if u == nil {
		log.Println("Email not exists")
		http.Error(w, "Email not exists", http.StatusForbidden)
		return
	}
	// does the entered password match the stored password?
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(person.Password))
	if err != nil {
		log.Println("Email and/or password do not match")
		http.Error(w, "Email and/or password do not match", http.StatusForbidden)
		return
	}

	log.Println("sucessfully login")

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: u.Email,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string >> get the complete signed token
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// create session
	c := &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	}
	// c.MaxAge = sessionLength
	fmt.Println("login cookie :", c.Value)
	http.SetCookie(w, c)

	// r := models.Resp{Data: u}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(u)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

//Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(w, r) {
		log.Println("Return to index :::")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, _ := r.Cookie("token")
	// delete the session
	// delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	c, err := r.Cookie("token")
	if err != nil {
		return false
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	//get value if user is valid
	//fgf
	s := c.Value
	c.Value = s

	// refresh session
	c.Expires = expirationTime
	http.SetCookie(w, c)
	return true
}
