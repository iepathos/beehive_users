// users_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"

	simplejson "github.com/bitly/go-simplejson"
	r "gopkg.in/gorethink/gorethink.v3"

	"github.com/iepathos/beehive/rego"
)

func TestCreateUser(t *testing.T) {
	// lookup user in rethinkdb and make sure it now exists
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	rego.CreateDatabase("test")
	rego.CreateTable(TableName)

	url := "/create"
	jsonStr := []byte(`{"username":"Saitama"}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("CreateUser handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body is what we expect.
	reqJSON, err := simplejson.NewFromReader(rr.Body)
	if err != nil {
		t.Errorf("Error while reading request JSON: %s", err)
	}
	username := reqJSON.Get("username").MustString()
	if username != "Saitama" {
		t.Errorf("Expected request JSON response to have username Saitama")
	}
	wins := reqJSON.Get("wins").MustInt()
	if wins != 0 {
		t.Errorf("Expected request JSON response to have wins 0")
	}

	db := r.DB("test")
	cursor, err := db.Table(TableName).Count().Run(session)
	if err != nil {
		log.Fatalln(err.Error())
	}
	var count int
	cursor.One(&count)
	cursor.Close()
	if count != 1 {
		t.Errorf("Expected RethinkDB users table to have count of 1")
	}
	rego.DropDatabase("test")
}

func TestGetUser(t *testing.T) {
	rego.CreateDatabase("test")
	rego.CreateTable(TableName)

	jsonStr := []byte(`{"username":"Saitama"}`)
	var user User
	err := json.Unmarshal(jsonStr, &user)
	if err != nil {
		t.Fatal(err)
	}
	InsertUser(user)

	url := "/user/Saitama"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUser)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("CreateUser handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	reqJSON, err := simplejson.NewFromReader(rr.Body)
	if err != nil {
		t.Errorf("Error while reading request JSON: %s", err)
	}
	username := reqJSON.Get("username").MustString()
	if username != "Saitama" {
		t.Errorf("Expected request JSON response to have username Saitama")
	}
	wins := reqJSON.Get("wins").MustInt()
	if wins != 0 {
		t.Errorf("Expected request JSON response to have wins 0")
	}

	rego.DropDatabase("test")
}
