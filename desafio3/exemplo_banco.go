package main

import (
	"log"
	"encoding/json"
	"net/http"
	"os"
	"time"
	"database/sql"
	"strings"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

//PORT port to be used
const PORT = "8080"

var db, _ = sql.Open("sqlite3", "file:database.sqlite?cache=shared")

type Script struct {
	Actor 	string	`json:"actor"`
	Quote 	string	`json:"quote"`
}

func main() {
	r := mux.NewRouter()
	http.Handle("/", r)
	r.Handle("/v1/quote", quote()).Methods("GET", "OPTIONS")
	r.Handle("/v1/quote/{actor}", quoteByActor()).Methods("GET", "OPTIONS")
	logger := log.New(os.Stderr, "logger: ", log.Lshortfile)
	srv := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":" + PORT,
		Handler:      context.ClearHandler(http.DefaultServeMux),
		ErrorLog:     logger,
	}
	errr := srv.ListenAndServe()
	if errr != nil {
		log.Fatal(errr)
	}
}

func quote() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT actor, detail as quote FROM scripts WHERE actor IS NOT NULL AND detail IS NOT NULL LIMIT 1")
		defer rows.Close()
	
		if err != nil {
			log.Fatal(err)
		}
		
		script := new(Script)
		rows.Next()
		rows.Scan(&script.Actor, &script.Quote)

		response, _ := json.Marshal(script)

		w.Write(response)
	})
}

func quoteByActor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		actor := strings.Replace(params["actor"], "+", " ", 1)

		rows, err := db.Query("SELECT actor, detail as quote FROM scripts WHERE actor IS NOT NULL AND detail IS NOT NULL AND actor like $1 LIMIT 1", actor)
		defer rows.Close()
	
		if err != nil {
			log.Fatal(err)
		}
		script := new(Script)
		rows.Next()
		rows.Scan(&script.Actor, &script.Quote)
		response, _ := json.Marshal(script)

		w.Write(response)
	})
}