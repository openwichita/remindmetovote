package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	c, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", c.DB)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/twilio/incoming", twilioIncomingHandler)

	log.Print("Binding to ", c.Bind)
	log.Fatal(http.ListenAndServe(c.Bind, nil))
}

func twilioIncomingHandler(rw http.ResponseWriter, r *http.Request) {
	from := r.FormValue("From")
	body := strings.ToLower(r.FormValue("Body"))

	log.Print("From: ", from)
	log.Print("Body: ", body)

	switch body {
	case "sign up", "sign-up", "signup", "subscribe":
		if _, err := db.Exec("insert into subscriptions (number) values ($1)", from); err != nil {
			respondError(rw, err)
			return
		}

		rw.Write([]byte("You've been subscribed!"))
	case "stop", "remove", "unsubscribe":
		if _, err := db.Exec("delete from subscriptions where number = $1", from); err != nil {
			respondError(rw, err)
			return
		}

		rw.Write([]byte("You've been unsubscribed!"))
	default:
		rw.WriteHeader(http.StatusBadRequest)
	}
}

func respondError(rw http.ResponseWriter, err error) {
	log.Print(err)
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}
