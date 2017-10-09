package main

import (
	"net/http"
)

import (
	"./templates"
	"time"
	"github.com/jcuga/golongpoll"
)

//Set Handlers for content requests
func HandleHttpRequests(lpManager *golongpoll.LongpollManager) {

	//!!!!__________GENERAL STYLES__________!!!!
	//Home Page
	http.HandleFunc("/chat/style.css", templates.Style)

	//!!!!__________GENERAL WEBPAGES__________!!!!
	//Home Page
	http.HandleFunc("/", templates.Index)
	//Login Page
	http.HandleFunc("/chat", Login(lpManager))
	//Logout Page
	http.HandleFunc("/chat/logout", Logout(lpManager))
	//Chat Page
	http.HandleFunc("/chat/main", Chat)

	http.NotFoundHandler()

}

//Login
func Login(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	// Creates closure that captures the LongpollManager
	return func(w http.ResponseWriter, r *http.Request) {
		//Get username
		user := r.URL.Query().Get("user")

		//Confirm a valid user
		if user == "" || user == "System" {
			templates.ErrorInvalid(w, r)
			return
		}

		//Username not taken
		_, exist := existUser(user)
		if exist {
			templates.ErrorTaken(w, r)
			return
		}

		//Continue if positive
		hash := GetMD5Hash(user + GARBAGE)
		admin := user == "Onelio" //Give Admin just to Onelio
		//Set Cookie in Client
		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "hash", Value: hash, Expires: expiration}

		http.SetCookie(w, &cookie)
		templates.Login(w, r)
		//Add user to pool
		manager.Publish( "public_actions", Message{Action:"notice", User:"System", Text:"User " + user + " has entered"})
		users[hash] = User{len(users), user, hash, admin}
	}

}

//Logout
func Logout(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	// Creates closure that captures the LongpollManager
	return func(w http.ResponseWriter, r *http.Request) {
		hash, _ := r.Cookie("hash")
		//Confirm hash is valid
		if hash == nil {
			templates.ErrorInvalid(w, r)
			return
		}
		user, ok := users[string(hash.Value)]
		if !ok {
			templates.ErrorInvalid(w, r)
			return
		}
		//Delete local user
		manager.Publish( "public_actions", Message{Action:"notice", User:"System", Text:"User " + user.name + " has left"})
		delete(users, hash.Value)

		//Continue if positive
		templates.Logout(w, r)
	}
}

//Chat
func Chat(w http.ResponseWriter, r *http.Request) {
	// Creates closure that captures the LongpollManager
	hash, _ := r.Cookie("hash")
	//Confirm hash is valid
	if hash == nil {
		templates.ErrorInvalid(w, r)
		return
	}
	value, ok := users[string(hash.Value)]
	if !ok {
		templates.ErrorInvalid(w, r)
		return
	}

	//Continue if positive
	templates.Chat(w, r, value.name)
}

