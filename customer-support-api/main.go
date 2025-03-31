package main

import "fmt"

type UserQuery struct {
	Query          string `json:"query"`
	ConversationId int    `json:"conversation_id"`
}

func receiveUserQuery(query UserQuery) {

}

func main() {
	fmt.Println("hello")
}
