package main

import (
	"crypto/md5"
	"encoding/hex"
)

type User struct {
	id int
	name string
	hash string
	admin bool
}

type Message struct {
	User     string `json:"user"` //Maker
	Dest     string `json:"dest"` //Objetive
	Action   string `json:"action"` //Extra actions, default text
	Text	 string `json:"text"` //Message in it
}

const GARBAGE string = "1234abcd"

var users = make(map[string]User)

func existUser(user string) (*User, bool) {

	for _, b := range users {
		if b.name == user {
			return &b, true
		}
	}
	return nil, false
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
