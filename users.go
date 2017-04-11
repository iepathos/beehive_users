// simple users service
package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"

	r "gopkg.in/gorethink/gorethink.v3"
)

// database name
var DbName = os.Getenv("DBNAME")

// table name for users service
var TableName = "users"

// user model
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Wins     int    `json:"wins"`
}

// db calls
func InsertUser(user User) {
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
}

// user views

// create a new user
func CreateUser(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := req.Body.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(422)
		log.Println(err.Error())
	}

	InsertUser(user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// get
func GetUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	username := vars["username"]
	username = "Saitama"
	log.Println(username)
	var user User
	// connect to db
	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost:28015",
	})
	if err != nil {
		log.Fatalln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	db := r.DB(DbName)

	res, err := db.Table(TableName).Filter(r.Row.Field("username").Eq(username)).Run(session)
	if err != nil {
		log.Fatalln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = res.One(&user)
	if err != nil {
		log.Fatalln(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// list

// login
// logout

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/create", CreateUser)
	router.HandleFunc("/user/{username}", GetUser)
	server := &http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}
	log.Fatal(server.ListenAndServe())
}
