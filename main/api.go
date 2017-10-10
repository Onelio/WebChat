package main

import (
	"net/http"
	"github.com/jcuga/golongpoll"
	"os"
)

//Set Handlers for content requests
func HandleApiRequests(lpManager *golongpoll.LongpollManager) {

	//!!!!__________GENERAL EXE__________!!!!
	//Send Messages
	http.HandleFunc("/chat/api/sendMessage", SendMessage(lpManager))

	//Request Messages
	http.HandleFunc("/chat/api/requestMessage", GetMessage(lpManager))

}

//Execute Actions
func Execute(action string, fuser User, user string, manager *golongpoll.LongpollManager) (bool){
	message := Message{Action:"notice", User:"System"}
	switch action {
	case "help":
		message.Text = "Command List:  /w <user> <message> To Shout"
		manager.Publish(fuser.name + "_actions", message)
		break
	case "exit":
		if fuser.admin {
			message.Text = "Server shuting down by " + fuser.name
			manager.Publish("public_actions", message)
			os.Exit(0)
		}
		break
	case "shout":
		if fuser.admin {
			return false
		}
		break
	case "promote":
		tuser, exist := existUser(user)
		if exist && fuser.admin {
			tuser.admin = !tuser.admin
			users[tuser.hash] = *tuser
		}
		break
	case "ban":
		tuser, exist := existUser(user)
		if exist && fuser.admin {
			delete(users, tuser.hash)
			users[tuser.hash + "ban"] = *tuser
			message.Text = "User " + user + " has been banned by " + fuser.name
			manager.Publish("public_actions", message)
		}
		break
	case "kick":
		tuser, exist := existUser(user)
		if exist && fuser.admin {
			delete(users, tuser.hash)
			message.Text = "User " + user + " has been kicked by " + fuser.name
			manager.Publish("public_actions", message)
		}
		break
	default:
		return false
	}
	return true
}

//Allow users to publish messages
func SendMessage(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	// Creates closure that captures the LongpollManager
	return func(w http.ResponseWriter, r *http.Request) {
		//Composition
		hash, _ := r.Cookie("hash") //From
		user := r.URL.Query().Get("user") //To
		action := r.URL.Query().Get("action")
		text := r.URL.Query().Get("text")
		message := Message{Action:action, Text:text}
		isPrivate := false

		//Confirm that hash is valid
		fuser, ok := users[string(hash.Value)]
		if hash == nil || !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing required URL param."))
			return
		}
		message.User = fuser.name //Set message maker

		//Confirm that has text
		if len(text) == 0 && action=="say" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing required URL param."))
			return
		}

		if len(user) > 0 { //If user is bigger than 0 it's a wishper
			isPrivate = true
			tuser, exist := existUser(user)
			if exist {
				message.Dest = tuser.name // Set destination
			} else { //If not exist let's just return
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Missing required URL param."))
				return
			}
		}

		if Execute(action, fuser, user, manager) {
			return
		}

		if !isPrivate {
			manager.Publish("public_actions", message)
		} else {
			//From
			message.User = "To " + user
			manager.Publish(fuser.name + "_actions", message)
			//To
			message.User = "From " + fuser.name
			manager.Publish(user + "_actions", message)
		}
	}
}

//Allow users to get updates from messages
func GetMessage(manager *golongpoll.LongpollManager) func(w http.ResponseWriter, r *http.Request) {
	// Creates closure that captures the LongpollManager
	return func(w http.ResponseWriter, r *http.Request) {
		hash, _ := r.Cookie("hash")
		category := r.URL.Query().Get("category")

		//Confirm that hash is valid
		user, ok := users[string(hash.Value)]
		if hash == nil || !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing required URL param."))
			return
		}

		// Only allow supported subscription categories:
		if category != "public_actions" && category != user.name + "_actions" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Subscription channel does not exist."))
			return
		}

		// Client is either requesting the public stream, or a private
		// stream that they're allowed to see.
		// Go ahead and let the subscription happen:
		manager.SubscriptionHandler(w, r)
	}
}
