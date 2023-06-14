package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Ticket struct {
	FirstName    string
	LastName     string
	Email        string
	NumOfTickets int
}

var (
	tickets          []Ticket
	remainingTickets = 50
	mutex            sync.Mutex
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/book", bookHandler)
	http.HandleFunc("/thankyou", thankyouHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	if err := tmpl.Execute(w, remainingTickets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	firstName := r.FormValue("firstname")
	lastName := r.FormValue("lastname")
	email := r.FormValue("email")
	numOfTicketsStr := r.FormValue("numoftickets")

	numOfTickets, err := strconv.Atoi(numOfTicketsStr)
	if err != nil {
		http.Error(w, "Invalid number of tickets", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if numOfTickets > remainingTickets {
		http.Error(w, "Insufficient tickets", http.StatusBadRequest)
		return
	}

	ticket := Ticket{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		NumOfTickets: numOfTickets,
	}
	tickets = append(tickets, ticket)
	remainingTickets -= numOfTickets

	http.Redirect(w, r, "/thankyou", http.StatusSeeOther)
}

func thankyouHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("thankyou.html"))
	if err := tmpl.Execute(w, remainingTickets); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
