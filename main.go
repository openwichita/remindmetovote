package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

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
	to := r.FormValue("To")

	var message string

	if err := logMessage(from, to, body); err != nil {
		respondError(rw, err)
		return
	}

	switch body {
	case "sign up", "sign-up", "signup", "subscribe":
		if _, err := db.Exec("insert into subscriptions (number) values ($1)", from); err != nil {
			respondError(rw, err)
			return
		}
		message = "You've been subscribed! Text STOP to unsubscribe."
	case "stop", "remove", "unsubscribe":
		if _, err := db.Exec("delete from subscriptions where number = $1", from); err != nil {
			respondError(rw, err)
			return
		}
		message = "You've been unsubscribed! Text SIGNUP to resubscribe."
	default:
		message = "Text SIGNUP to subscribe for voting reminders."
	}

	if err := logMessage(to, from, message); err != nil {
		respondError(rw, err)
		return
	}
	rw.Write([]byte(message))
}

func respondError(rw http.ResponseWriter, err error) {
	log.Print(err)
	rw.WriteHeader(http.StatusBadRequest)
}

func logMessage(from string, to string, message string) error {
	_, err := db.Exec("insert into message_log (from_number, to_number, body, created_at) values ($1, $2, $3, $4)", from, to, message, time.Now())
	return err
}
