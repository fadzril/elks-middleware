package main

import "time"

// Subscription struct
type Subscription struct {
	Type             string
	App_Id           string
	Component        string
	Hostname         string
	Transaction_Type string
	Subscription_Id  string
	Category         string
	Status           string
	Messages         string
	Tags             []string
	Ts               time.Time `json: "created"`
}

// Subscriptions List struct
type Subscriptions struct {
	Total   int
	Message string
	Result  []Subscription
}
