package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/mux"
	"github.com/rcliao/sql-unit-test/db"
	"github.com/rcliao/sql-unit-test/web"

	_ "github.com/go-sql-driver/mysql"
)

/*
Server.go runs the web server for the sql-unit-test to provide easier access to
test sql without installing CLI or its dependencies (e.g. MySQL)
*/

func main() {
	r := mux.NewRouter()
	sqlDB := getDB()
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	factory := db.NewFactory(sqlDB, session)

	r.HandleFunc("/hello", web.Hello()).Methods("GET")
	r.HandleFunc("/", web.Index(factory)).Methods("GET")
	r.HandleFunc("/health", web.HealthCheck(sqlDB, session)).Methods("GET")
	r.HandleFunc("/{subject}", web.Index(factory)).Methods("GET")
	r.HandleFunc("/api/test", web.RunTest(factory)).Methods("POST")
	r.PathPrefix("/static").Handler(web.Static())

	log.Println("Running web server at port 8000")
	http.ListenAndServe(":8000", r)
}

func getDB() *sql.DB {
	defaultProtocol := "tcp"
	defaultPort := "3306"

	sqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/",
		os.Getenv("MYSQL_USERNAME"),
		os.Getenv("MYSQL_PASSWORD"),
		defaultProtocol,
		os.Getenv("MYSQL_HOST"),
		defaultPort,
	)

	db, err := sql.Open("mysql", sqlDSN)
	if err != nil {
		panic(err)
	}

	return db
}
