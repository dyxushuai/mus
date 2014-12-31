package main

import (
	"fmt"
	"encoding/json"
)

type A struct {

}

type B struct {

}


type serverOptions struct {
	A
	B
	Port     string `json:"port"`
	Password string `json:"password"`
	Method   string `json:"method"`
	current  int64	`json:"current"`
	Limit    int64  `json:"limit"`
	Timeout  int64  `json:"timeout"`
}


func main() {
	var aa = A{}
	var bb = B{}
	a := &serverOptions{
		aa,
		bb,
		Port: "9090",
		Password: "123456",
		Method:   "rc4",
	}

	b, _ := json.Marshal(a)
	fmt.Println(string(b))
}
