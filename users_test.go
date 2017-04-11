// handlers_test.go
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
)

func TestCreateUser(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
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
	// expected := `{"username":"Saitama","wins":0}`
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
	// if rr.Body.String() != expected {
	//  t.Errorf("CreateUser handler returned unexpected body: got %v want %v",
	//      rr.Body.String(), expected)
	// }
}
