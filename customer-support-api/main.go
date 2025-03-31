package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserQuery struct {
	Query          string `json:"query"`
	ConversationId int    `param:"conversationId" json:"conversation_id"`
}

// Knowledge Base
var answersDB = map[string]string{
	"request_help":    "I will help you as much as I can!",
	"password_reset":  "To reset your password, navigate to your profile...",
	"order_status":    "You can check the status of your order in dashboard",
	"refund_request":  "Refunds are discussed directly with the seller",
	"account_issue":   "What is the exact issue you're having?", // Probably should be a better intent, but ehh
	"billing_inquiry": "Issues with billing should be discussed with the bank",
	"technical_issue": "We're working to fix this issue!",
}

// The record of a conversation
var supportMessages = map[int][]string{}

func receiveUserQuery(c echo.Context) error {
	userQuery := UserQuery{}

	err := (&echo.DefaultBinder{}).BindBody(c, &userQuery)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	if len(userQuery.Query) == 0 {
		return c.String(http.StatusBadRequest, "No user query specified")
	}

	conversationIdStr := c.Param("conversationId")
	conversationId, err := strconv.Atoi(conversationIdStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Couldn't convert conversation id")
	}

	userQuery.ConversationId = conversationId

	// If the conversation does not exist, then create it in the database
	_, exists := supportMessages[userQuery.ConversationId]
	if !exists {
		supportMessages[userQuery.ConversationId] = make([]string, 0)
	}

	supportMessages[userQuery.ConversationId] = append(supportMessages[userQuery.ConversationId], userQuery.Query)
	// In a proper implementation there needs a check whether the conversation has been delegated to a human or not
	// I see no reason to implement it yet
	err = sendUserQueryToNLP(userQuery)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Sent for processing")
}

type UserIntent struct {
	Confidence     float64 `json:"confidence"`
	Intent         string  `json:"intent"`
	ConversationId int     `json:"conversation_id"`
}

func receiveUserConversationIntent(c echo.Context) error {
	// Assuming this endpoint is not accessible by normal users, as there should be some sort of validation in real world
	// So this data can mostly be trusted (Basically my excuse to do, frankly, lackluster validation)
	var userIntent UserIntent
	err := c.Bind(&userIntent)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	_, exists := supportMessages[userIntent.ConversationId]
	if !exists {
		return c.String(http.StatusNotFound, "The conversation was not found")
	}

	// If the nlp confidence is low, then delegate to a human
	if userIntent.Confidence < 0.5 {
		// Mark the conversation here as delegated to a human agent
		// Could just be a flag or an enum. Not implemented in the mock db, as no real need
		supportMessages[userIntent.ConversationId] = append(supportMessages[userIntent.ConversationId], "A human agent will be in contact with you soon")
		return c.String(http.StatusOK, "Received user intent")
	}

	// Retrieve answer from knowledge base
	intentReply := answersDB[userIntent.Intent]
	supportMessages[userIntent.ConversationId] = append(supportMessages[userIntent.ConversationId], intentReply)

	return c.String(http.StatusOK, "Received used intent")
}

func getSupportMessages(c echo.Context) error {
	conversationIdStr := c.Param("conversationId")
	conversationId, err := strconv.Atoi(conversationIdStr)
	if err != nil {
		return c.String(http.StatusBadRequest, "Couldn't convert conversation id")
	}

	conversation, exists := supportMessages[conversationId]
	if !exists {
		return c.String(http.StatusNotFound, "Conversation not found")
	}

	return c.JSON(http.StatusOK, conversation)
}

func main() {
	e := echo.New()
	e.POST("/support/:conversationId/query", receiveUserQuery)
	e.POST("/nlp/intent", receiveUserConversationIntent)

	// For realtime communication this could be a websocket connection or any other polling form of connection
	// However this is simpler and demonstrates the idea well enough
	e.GET("/support/:conversationId", getSupportMessages)
	e.Logger.Fatal(e.Start(":1323"))
}

func sendUserQueryToNLP(query UserQuery) error {
	url := "http://localhost:8000/intent"

	jsonStr, err := json.Marshal(&query)
	if err != nil {
		return errors.New("failed to convert to json")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("nlp server error")
	}

	return nil
}
