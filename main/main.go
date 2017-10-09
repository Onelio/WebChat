package main

import (
	"fmt"
	"log"
	"net/http"
)

import (
	"github.com/jcuga/golongpoll"
)

//Server Settings
var serverSettings = golongpoll.Options{
	LoggingEnabled: false,
}

func main() {
	//Set configuration
	manager, err := golongpoll.StartLongpoll(serverSettings)
	if err != nil {
		log.Fatalf("Failed to create manager: %q", err)
		return
	}
	//Set Handlers for webpage
	HandleHttpRequests()
	//Set Handlers for api
	HandleApiRequests(manager)

	//Start Listening
	fmt.Println("Serving webpage at http://127.0.0.1:8081/")
	http.ListenAndServe("127.0.0.1:8081", nil) //127.0.0.1
}
