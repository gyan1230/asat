package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gyan1230/asat/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//User ...
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Fullname string             `json:"fullname,omitempty" bson:"fullname,omitempty"`
}

//Claims ...
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//Resp :
type Resp struct {
	Data interface{} `json:"data"`
}

//Session ...
type Session struct {
	Un           string
	LastActivity time.Time
}

//Allusers :
func Allusers(r *http.Request) ([]User, error) {
	var users []User
	collection := config.Client.Database("userDb").Collection("user")
	// ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(r.Context(), bson.M{})
	if err != nil {
		log.Println("error in find....", err)
		return nil, err
	}
	defer cursor.Close(r.Context())
	for cursor.Next(r.Context()) {
		var person User
		cursor.Decode(&person)
		users = append(users, person)
	}
	if err := cursor.Err(); err != nil {
		log.Println("error in cursor....", err)
		return nil, err
	}
	log.Println("Sucessfully return all users from db:::")
	return users, nil
}

//GetUser :
func GetUser(ctx context.Context, em string) (*User, error) {
	person := &User{}
	collection := config.Client.Database("userDb").Collection("user")
	err := collection.FindOne(ctx, User{Email: em}).Decode(person)
	if err != nil {
		log.Println("No user exists in User Database :", err)
		return nil, err
	}
	log.Println("User exists :", person.Email)
	return person, nil
}

/*
//FindUser :
func FindUser(w http.Request, r *http.Request) bool {
	u := User{}
	json.NewDecoder(r.Body).Decode(&u)
	err := config.Users.Find(bson.M{"email": u.Email}).One(&u)
	if err != nil {
		return true
	}
	return false

}



//VerifyPassword :
func VerifyPassword(rcv string, r *http.Request) (bool, error) {
	u := User{}
	json.NewDecoder(r.Body).Decode(&u)

	//TODO fetch from DB and check
	err := config.Users.Find(bson.M{"password": u.Password}).One(&u)
	if err != nil {
		log.Println("Error find in password in DB")
		return false, errors.New("500. Internal Server Error." + err.Error())

	}
	// does the entered password match the stored password?
	err = bcrypt.CompareHashAndPassword([]byte(rcv), []byte(u.Password))
	if err != nil {
		log.Println("Email and/or password do not match")
		return false, errors.New("Email and/or password do not match" + err.Error())
	}
	log.Println("Password verify done.")
	return true, nil

}

//GetPassword :
func GetPassword(r *http.Request) (bool, error) {
	u := User{}
	json.NewDecoder(r.Body).Decode(&u)

	//TODO fetch from DB and check
	err := config.Users.Find(bson.M{"password": u.Password}).One(&u)
	if err != nil {
		log.Println("Error find in password in DB")
		return false, errors.New("500. Internal Server Error." + err.Error())

	}
	return true, nil

}

//DeleteUser :
func DeleteUser(r *http.Request) error {
	u := User{}
	json.NewDecoder(r.Body).Decode(&u)

	err := config.Users.Remove(bson.M{"email": u.Email})
	if err != nil {
		return errors.New("500. Internal Server Error." + err.Error())
	}

	return nil
}
*/
