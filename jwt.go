package main

import (
	"fmt"
)

type jwt struct {
	Header    header
	Payload   string
	Signature string
}

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
}

func authenticate_jwt(token string) bool {
	fmt.Println(token)

	return true
}

// think will need a jwt middleware to auth the route and use this file
// to build the jwt
// how does data look like on client and server?
// what encoding do i need to handle and think about
