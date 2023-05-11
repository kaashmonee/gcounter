package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/kaashmonee/gcounter/server/workers"
	"github.com/kaashmonee/gcounter/utilities"
)

type server struct {
	*log.Logger
	nodes []workers.Worker
}

type Server interface {
	Serve()
}

var views int

// displayPage - chooses a node to display the page
func (s *server) displayPage(w http.ResponseWriter, r *http.Request) {
	s.Println("got / request")
	io.WriteString(w, strconv.Itoa(views))
}

func NewServer(numWorkers int) Server {
	workersSlice := make([]workers.Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		w := workers.NewWorker(numWorkers, i)
		workersSlice[i] = w
	}
	return &server{
		Logger: utilities.NewLogger(nil),
		nodes:  workersSlice,
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
