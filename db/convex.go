package db

import (
	"fmt"
)

var CONVEX_API_KEY = "dev:patient-dachshund-614|eyJ2MiI6ImI4NDg2YThiMjkxZTQxMTNhM2U5Nzc4ZWVhYjhhNjA3In0="
var HTTP_ACTIONS_URL = "https://patient-dachshund-614.convex.cloud"
var ACTIONS_URL = HTTP_ACTIONS_URL + "/api/action"

type Db struct{}

type ConvexClient struct{}

func New() *Db {
	return &Db{}
}

func (db *Db) NewClient() *ConvexClient {
	return &ConvexClient{}
}

func (c *ConvexClient) Log(id string, message string) {
	fmt.Printf("[%s] %s\n", id, message)
}
