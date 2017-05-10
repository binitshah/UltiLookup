/*
 * Hello.go
 *
 * Author: Binit Shah
 * Description: A simple web service that gather people information from mongodb and displays an angular website with the information.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Address   string `json:"address"`
	State     string `json:"state"`
	City      string `json:"city"`
	Zip       string `json:"zip"`
}

func main() {
	session, err := mgo.Dial("10.71.74.148")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	router := mux.NewRouter()
	router.HandleFunc("/", getAllPeople(session)).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS()(router)))
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("People").C("People")

	index := mgo.Index{
		Key:        []string{"_id"},
		Unique:     false,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func getAllPeople(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB("People").C("People")

		var people []Person
		err := c.Find(bson.M{}).All(&people)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed get all people: ", err)
			return
		}

		respBody, err := json.MarshalIndent(people, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
