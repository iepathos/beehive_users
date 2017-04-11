// simple users service
package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	r "gopkg.in/gorethink/gorethink.v3"
)

// database name
var DbName = "beehive"

// table name for users service
var TableName = "users"

// user model
type User struct {
	Username string `json:"username"`
	Wins     int    `json:"wins"`
}

// user views

// create a new user
func CreateUser(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := req.Body.Close(); err != nil {
		panic(err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(422)
		log.Println(err.Error())
	}

	// connect to db
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	db := r.DB(DbName)

	// insert user db
	err = db.Table(TableName).Insert(map[string]interface{}{
		"username": user.Username,
		"wins":     0,
	}).Exec(session)
	if err != nil {
		log.Println(err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		panic(err)
	}
}

// list
// get
// login
// logout

func main() {
	http.HandleFunc("/create", CreateUser)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
