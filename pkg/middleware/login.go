package middleware

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		t, err := template.ParseFiles("assets/templates/login.html")
		if err != nil {
			log.Fatal("login ", r.Method, ": ", err)
		}
		t.Execute(w, nil)
	} else {
		min := 18
		max := 36
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"][0:])
		fmt.Println("password:", r.Form["password"])
		l := len(r.Form["password"])
		if !(l >= min && l <= max) {
			fmt.Println("invalid password")
		}
	
		redirectURL := fmt.Sprintf("/login?username=%s", r.Form["username"])
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}