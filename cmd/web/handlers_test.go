package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	rr := httptest.NewRecorder()

	// init a new dummy request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		// Fatal anywhere in this test completely stops any further execution
		t.Fatal(err)
	}

	// call the ping handler function
	ping(rr, r)

	// the Result() method on the http.ResponseRecorder gets the
	// http.Response generated by the ping handler.
	rs := rr.Result()

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
