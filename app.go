package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var (
	version = "undefined"
	apiKey  = ""
	db      *bolt.DB
)

func init() {
	apiKey = os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal("No API Key setup")
	}
	db, _ = bolt.Open("/db/my.db", 0600, nil)
}

// Response is just a very basic example.
type Response struct {
	Status string `json:"status,omitempty"`
}

// GetVersion returns always the same response.
func GetVersion(w http.ResponseWriter, _ *http.Request) {
	b := Response{Status: "idle"}
	json.NewEncoder(w).Encode(b)
}

// GetRoot returns always the same response.
func GetRoot(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://www.google.com", 302)
}

// SetRoot returns always the same response.
func SetRoot(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("ApiKey") != apiKey {
		http.Error(w, "Invalid Key", 401)
	} else {
		http.Redirect(w, r, "https://www.google.com", 302)
	}

}

func main() {
	fmt.Println("Starting Version " + version)
	router := mux.NewRouter()
	router.HandleFunc("/", SetRoot).Methods("POST")
	router.HandleFunc("/", GetRoot).Methods("GET")
	router.HandleFunc("/vendor", GetVersion).Methods("Get")
	log.Fatal(http.ListenAndServe(":10987", router))
}
