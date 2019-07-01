package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gyan1230/asat/controllers"
)

func main() {
	log.Println("Starting application :::::::::")
	http.HandleFunc("/", index)
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/admin", controllers.ShowAll)
	http.HandleFunc("/login", controllers.Login)
	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/tweetData", controllers.GetTweetData)
	http.HandleFunc("/data", controllers.DisplayAllPowerData)
	http.HandleFunc("/storeData", controllers.StoreEnergyData)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Server listen error:", err)
	}

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	log.Println("In index:::::::::::::::")
	json.NewEncoder(w).Encode("In Index:::")
}
