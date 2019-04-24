package main

import (
	"encoding/json"
	"flag"
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
	dbLocation := flag.String("db", "my.db", "Full path to the database file")
	flag.Parse()

	apiKey = os.Getenv("APIKEY")
	if apiKey == "" {
		log.Fatal("No API Key setup")
	}

	var err error
	db, err = bolt.Open(*dbLocation, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Urls"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}

// Response is just a very basic example.
type Response struct {
	Version string `json:"version,omitempty"`
}

// GetVersion returns always the same response.
func GetVersion(w http.ResponseWriter, _ *http.Request) {
	b := Response{Version: version}
	json.NewEncoder(w).Encode(b)
}

// GetRoot returns always the same response.
func GetRoot(w http.ResponseWriter, r *http.Request) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Urls"))
		url := b.Get([]byte(r.URL.Path))
		http.Redirect(w, r, string(url), 302)
		return nil
	})
}

// SetRoot returns always the same response.
func SetRoot(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("ApiKey") != apiKey {
		http.Error(w, "Invalid Key", 401)
	} else {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Urls"))
			err := b.Put([]byte(r.URL.Path), []byte(r.Header.Get("Url")))
			return err
		})

		if err != nil {
			log.Fatal(err)
		}

	}

}

func main() {
	fmt.Println("Starting Version " + version)
	router := mux.NewRouter()
	router.HandleFunc("/version", GetVersion).Methods("Get")
	router.HandleFunc("/{url}", GetRoot).Methods("GET")
	router.HandleFunc("/{url}", SetRoot).Methods("POST")
	//log.Fatal(http.ListenAndServe(":10987", router))
	log.Fatal(http.ListenAndServe(":5000", router))
}
