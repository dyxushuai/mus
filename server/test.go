package main

import (
//	"github.com/JohnSmithX/mus/server/api"
	"fmt"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

type Response2 struct {
	Page   int      //`json:"page,omitempty"`
	Fruits []string `json:"fruits_a"`
}
func main() {
	a, b, c, d, e := uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4(), uuid.NewV4()
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(e)

	w := &Response2{
		Page: 3,
		Fruits: []string{"apple"},
	}
	r, _ := json.Marshal(w)
	fmt.Println(string(r))
}

