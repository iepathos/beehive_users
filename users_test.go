// handlers_test.go
package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	r "gopkg.in/gorethink/gorethink.v3"
)

func createDatabase(databaseName string) {
	// connect to rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Creating database", databaseName)
	_, err = r.DBCreate(databaseName).Run(session)
	if err != nil {
		log.Println(err.Error())
	}
}

func createTable(tableName string) {
	// connect to rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	db := r.DB(os.Getenv("DBNAME"))

	log.Println("Creating table", tableName)
	if _, err := db.TableCreate(tableName).RunWrite(session); err != nil {
		log.Println(err)
	}
}

func dropDatabase(databaseName string) {
	// connect to rethinkdb
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Dropping database", databaseName)
	_, err = r.DBDrop(databaseName).Run(session)
	if err != nil {
		log.Println(err.Error())
	}
}

func TestCreateUser(t *testing.T) {
	// lookup user in rethinkdb and make sure it now exists
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	createDatabase("test")
	createTable(TableName)

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
	id := reqJSON.Get("id").MustInt()
	if id != 0 {
		t.Errorf("Expected request JSON response to have id 0 got %v", id)
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
	dropDatabase("test")
}

func TestGetUser(t *testing.T) {

}
