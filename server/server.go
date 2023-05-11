package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type server struct {
}

type Server interface {
	Serve()
}

var views int

func displayPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got / request")
	io.WriteString(w, strconv.Itoa(views))
}

func NewServer() Server {
	return &server{}
}

func (s *server) Serve() {
	http.HandleFunc("/", displayPage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err:", err)
	}
	return
}
