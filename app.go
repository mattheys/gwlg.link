package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var version = "undefined"

// Response is just a very basic example.
type Response struct {
	Status string `json:"status,omitempty"`
}

// GetVendor returns always the same response.
func GetVendor(w http.ResponseWriter, _ *http.Request) {
	b := Response{Status: "idle"}
	json.NewEncoder(w).Encode(b)
}

// GetRoot returns always the same response.
func GetRoot(w http.ResponseWriter, r *http.Request) {

	b := Response{Status: "YesIAmNew"}
	json.NewEncoder(w).Encode(b)
}

func main() {
	fmt.Println(version)
	router := mux.NewRouter()
	router.HandleFunc("/", GetRoot).Methods("GET")
	router.HandleFunc("/vendor", GetVendor).Methods("Get")
	log.Fatal(http.ListenAndServe(":10987", router))
}
