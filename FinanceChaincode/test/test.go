package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name       string   `json:"name"`
	Uid        string   `json:"uid"`
	CompactIDs []string `json:"compactIDs"`
}

func main() {
	var user1 User
	var user2 User
	user1.Name = "heihei"
	user1.Uid = "ajldfladf"
	// user := User{Name: "heihei", Uid: "ajfdlajfdl"}
	userBytes, _ := json.Marshal(user1)
	err := json.Unmarshal(userBytes, &user2)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(string(userBytes))
	fmt.Println(user1)
	fmt.Println(user2)
}
