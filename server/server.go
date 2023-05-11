package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/kaashmonee/gcounter/utilities"
)

type server struct {
	*log.Logger
}

type Server interface {
	Serve()
}

var views int

func (s *server) displayPage(w http.ResponseWriter, r *http.Request) {
	s.Println("got / request")
	io.WriteString(w, strconv.Itoa(views))
}

func NewServer() Server {
	return &server{
		Logger: utilities.NewLogger(nil),
	}
}

func (s *server) Serve() {
	http.HandleFunc("/", s.displayPage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("err:", err)
	}
	return
}
