package main

import (
	"github.com/squiidz/bone"
	"net/http"
	"fmt"
)


type Home struct {

}

func (self *Home) ServeHTTP(rw http.ResponseWriter,rq *http.Request) {
	fmt.Println(rw, rq)
}

func main() {
	home := &Home{}
	mux := bone.New()
	mux.Get("/", home)
	http.ListenAndServe(":2444", mux)
}
