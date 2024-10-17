package main

import (
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
)

type application struct {
	Session *scs.SessionManager
}

func main() {
	// set up app config
	app := application{}
	// get a session manager
	app.Session = getSession()
	// print out a message
	log.Println("Starting server on: http://localhost:8080")
	// start the server
	err := http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
