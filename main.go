package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"send_email_poc/mailer"
)

func main() {
	err := mailer.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/sendEmails", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			do500(w, "failed to read request body: "+err.Error())
			return
		}
		defer r.Body.Close()

		var emails mailer.Emails
		err = json.Unmarshal(body, &emails)
		if err != nil {
			do400(w, "failed to parse request body: "+err.Error())
			return
		}

		// create our mailer (using the correct interface)
		sender, err := mailer.NewMailer()
		if err != nil {
			log.Print(err)
			do500(w, "error getting mailer: "+err.Error())
			return
		}

		for _, email := range emails.Emails {
			email := email

			err := sender.SendEmail(r.Context(), &email)
			if err != nil {
				log.Printf("Error sending email %#v: %v", email, err)
			}
		}

		// TODO: match output from spec/real BOP
		w.Write([]byte("OK"))
	}))

	log.Printf("Starting /sendEmail server on :8000")
	http.ListenAndServe(":8000", nil)
}

func do500(w http.ResponseWriter, msg string) {
	doError(w, msg, 500)
}

func do400(w http.ResponseWriter, msg string) {
	doError(w, msg, 400)
}

func doError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(msg))
}
