package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rcliao/sql-unit-test/web"
)

/*
Server.go runs the web server for the sql-unit-test to provide easier access to
test sql without installing CLI or its dependencies (e.g. MySQL)
*/

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", web.Hello()).Methods("GET")

	http.ListenAndServe(":8000", r)
}
